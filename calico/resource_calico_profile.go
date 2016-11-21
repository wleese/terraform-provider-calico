package calico

import (
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/projectcalico/libcalico-go/lib/api"
	"github.com/projectcalico/libcalico-go/lib/errors"
	caliconet "github.com/projectcalico/libcalico-go/lib/net"
)

func resourceCalicoProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceCalicoProfileCreate,
		Read:   resourceCalicoProfileRead,
		Update: resourceCalicoProfileUpdate,
		Delete: resourceCalicoProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"interface": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expected_ips": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"profiles": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dToProfileMetadata(d *schema.ResourceData) api.ProfileMetadata {
	metadata := api.ProfileMetadata{
		Name: d.Get("name").(string),
		Node: d.Get("node").(string),
	}

	if v, ok := d.GetOk("labels"); ok {
		labelMap := v.(map[string]interface{})
		labels := make(map[string]string, len(labelMap))

		for k, v := range labelMap {
			labels[k] = v.(string)
		}
		metadata.Labels = labels
	}

	return metadata
}

func dToProfileSpec(d *schema.ResourceData) (api.ProfileSpec, error) {
	spec := api.ProfileSpec{}
	spec.InterfaceName = d.Get("interface").(string)

	if v, ok := d.GetOk("expected_ips.#"); ok {
		ips := make([]caliconet.IP, v.(int))

		for i := range ips {
			ip := d.Get("expected_ips." + strconv.Itoa(i)).(string)
			validIP := net.ParseIP(ip)
			if validIP == nil {
				return spec, fmt.Errorf("expected_ips: %v is not IP", ip)
			}
			ips[i] = caliconet.IP{validIP}
		}

		if len(ips) != 0 {
			spec.ExpectedIPs = ips
		}
	}

	if v, ok := d.GetOk("profiles.#"); ok {
		profiles := make([]string, v.(int))

		for i := range profiles {
			profiles[i] = d.Get("profiles." + strconv.Itoa(i)).(string)
		}

		if len(profiles) != 0 {
			spec.Profiles = profiles
		}
	}

	return spec, nil
}

func resourceCalicoProfileCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	calicoClient := config.Client

	metadata := dToProfileMetadata(d)
	spec, err := dToProfileSpec(d)
	if err != nil {
		return err
	}

	profiles := calicoClient.Profiles()
	if _, err = profiles.Create(&api.Profile{
		Metadata: metadata,
		Spec:     spec,
	}); err != nil {
		return err
	}

	d.SetId(metadata.Name)
	return resourceCalicoProfileRead(d, meta)
}

func resourceCalicoProfileRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	calicoClient := config.Client

	profiles := calicoClient.profiles()
	profile, err := profiles.Get(api.ProfileMetadata{
		Name: d.Get("name").(string),
		Node: d.Get("node").(string),
	})

	// Handle endpoint does not exist
	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		}
	}

	d.SetId(profile.Metadata.Name)
	d.Set("name", profile.Metadata.Name)
	d.Set("node", profile.Metadata.Node)
	d.Set("labels", profile.Metadata.Labels)

	d.Set("profiles", profile.Spec.Profiles)
	d.Set("expected_ips", profile.Spec.ExpectedIPs)
	d.Set("interface", profile.Spec.InterfaceName)

	return nil
}

func resourceCalicoProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	calicoClient := config.Client

	profiles := calicoClient.profiles()

	// Handle non-existant resource
	metadata := dToprofileMetadata(d)
	if _, err := profiles.Get(metadata); err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		}
	}

	// Simply recreate the complete resource
	spec, err := dToprofileSpec(d)
	if err != nil {
		return err
	}

	if _, err = profiles.Apply(&api.Profile{
		Metadata: metadata,
		Spec:     spec,
	}); err != nil {
		return err
	}

	return nil
}

func resourceCalicoProfileDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	calicoClient := config.Client

	profiles := calicoClient.profiles()
	err := profiles.Delete(api.ProfileMetadata{
		Name: d.Get("name").(string),
		Node: d.Get("node").(string),
	})

	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); !ok {
			return fmt.Errorf("ERROR: %v", err)
		}
	}

	return nil
}
