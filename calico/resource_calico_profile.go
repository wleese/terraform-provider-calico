package calico

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/projectcalico/libcalico-go/lib/api"
	"github.com/projectcalico/libcalico-go/lib/errors"
	caliconet "github.com/projectcalico/libcalico-go/lib/net"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
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
			"spec": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingress": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"protocol": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"notProtocol": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"icmp": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"code": &schema.Schema{
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"notICMP": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"code": &schema.Schema{
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"source": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"net": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notNet": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"selector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notSelector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"ports": &schema.Schema{
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"notPorts": &schema.Schema{
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"destination": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"net": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notNet": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"selector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notSelector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"ports": &schema.Schema{
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"notPorts": &schema.Schema{
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"egress": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"protocol": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"notProtocol": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"icmp": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"code": &schema.Schema{
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"notICMP": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"code": &schema.Schema{
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"source": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"net": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notNet": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"selector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notSelector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"ports": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notPorts": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"destination": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"net": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notNet": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"selector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notSelector": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"ports": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"notPorts": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dToProfileMetadata(d *schema.ResourceData) api.ProfileMetadata {
	metadata := api.ProfileMetadata{
		Name: d.Get("name").(string),
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

func mapToRule(mapStruct map[string]interface{}) (api.Rule, error) {
	rule := api.Rule{}

	if val, ok := mapStruct["action"]; ok {
		rule.Action = val.(string)
	}
	if val, ok := mapStruct["protocol"]; ok {
		if len(val.(string)) > 0 {
			protocol := numorstring.ProtocolFromString(val.(string))
			rule.Protocol = &protocol
		}
	}
	if val, ok := mapStruct["notProtocol"]; ok {
		if len(val.(string)) > 0 {
			notProtocol := numorstring.ProtocolFromString(val.(string))
			rule.NotProtocol = &notProtocol
		}
	}
	if val, ok := mapStruct["icmp"]; ok {
		icmpList := val.([]interface{})
		if len(icmpList) > 0 {
			for _, v := range icmpList {
				icmpMap := v.(map[string]interface{})
				icmpType := icmpMap["type"].(int)
				icmpCode := icmpMap["code"].(int)
				icmp := api.ICMPFields{
					Type: &icmpType,
					Code: &icmpCode,
				}
				rule.ICMP = &icmp
			}
		}
	}
	if val, ok := mapStruct["notICMP"]; ok {
		icmpList := val.([]interface{})
		if len(icmpList) > 0 {
			for _, v := range icmpList {
				icmpMap := v.(map[string]interface{})
				icmpType := icmpMap["type"].(int)
				icmpCode := icmpMap["code"].(int)
				icmp := api.ICMPFields{
					Type: &icmpType,
					Code: &icmpCode,
				}
				rule.NotICMP = &icmp
			}
		}
	}
	if val, ok := mapStruct["source"]; ok {
		sourceList := val.([]interface{})

		// Only 1 source is allowed however
		if len(sourceList) > 0 {
			entityRule := api.EntityRule{}
			resourceRuleMap := sourceList[0].(map[string]interface{})

			if v, ok := resourceRuleMap["net"]; ok {
				_, n, err := caliconet.ParseCIDR(v.(string))
				if err != nil {
					return rule, err
				}
				entityRule.Net = n
			}
			if v, ok := resourceRuleMap["selector"]; ok {
				entityRule.Selector = v.(string)
			}
			if v, ok := resourceRuleMap["notSelector"]; ok {
				entityRule.NotSelector = v.(string)
			}
			if v, ok := resourceRuleMap["ports"]; ok {
				if resourcePortList, ok := v.([]interface{}); ok {
					portList, err := toPortList(resourcePortList)
					if err != nil {
						return rule, err
					}
					entityRule.Ports = portList
				}
			}
			if v, ok := resourceRuleMap["notPorts"]; ok {
				if resourcePortList, ok := v.([]interface{}); ok {
					portList, err := toPortList(resourcePortList)
					if err != nil {
						return rule, err
					}
					entityRule.NotPorts = portList
				}
			}

			rule.Source = entityRule
		}
	}

	return rule, nil
}

func toPortList(resourcePortList []interface{}) ([]numorstring.Port, error) {
	portList := make([]numorstring.Port, len(resourcePortList))

	for i, v := range resourcePortList {
		p, err := numorstring.PortFromString(v.(string))
		if err != nil {
			return portList, err
		}
		portList[i] = p
	}
	return portList, nil
}

func dToProfileSpec(d *schema.ResourceData) (api.ProfileSpec, error) {
	spec := api.ProfileSpec{}

	if v, ok := d.GetOk("spec.0.ingress.0.rule.#"); ok {
		ingressRules := make([]api.Rule, v.(int))

		for i := range ingressRules {
			mapStruct := d.Get("spec.0.ingress.0.rule." + strconv.Itoa(i)).(map[string]interface{})

			rule, err := mapToRule(mapStruct)
			if err != nil {
				return spec, err
			}

			ingressRules[i] = rule
		}
		spec.IngressRules = ingressRules
	}
	if v, ok := d.GetOk("spec.0.egress.0.rule.#"); ok {
		egressRules := make([]api.Rule, v.(int))

		for i := range egressRules {
			mapStruct := d.Get("spec.0.egress.0.rule." + strconv.Itoa(i)).(map[string]interface{})

			rule, err := mapToRule(mapStruct)
			if err != nil {
				return spec, err
			}

			egressRules[i] = rule
		}
		spec.EgressRules = egressRules
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

	profiles := calicoClient.Profiles()
	profile, err := profiles.Get(api.ProfileMetadata{
		Name: d.Get("name").(string),
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
	d.Set("labels", profile.Metadata.Labels)

	return nil
}

func resourceCalicoProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	calicoClient := config.Client

	profiles := calicoClient.Profiles()

	// Handle non-existant resource
	metadata := dToProfileMetadata(d)
	if _, err := profiles.Get(metadata); err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		}
	}

	// Simply recreate the complete resource
	spec, err := dToProfileSpec(d)
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

	profiles := calicoClient.Profiles()
	err := profiles.Delete(api.ProfileMetadata{
		Name: d.Get("name").(string),
	})

	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); !ok {
			return fmt.Errorf("ERROR: %v", err)
		}
	}

	return nil
}
