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

func nonEmptyEntityRule(entityRule *api.EntityRule) bool {
	state := false

	if len(entityRule.Tag) > 0 {
		state = true
	}
	if entityRule.Net != nil {
		state = true
	}
	if len(entityRule.Selector) > 0 {
		state = true
	}
	if len(entityRule.Ports) > 0 {
		state = true
	}
	if len(entityRule.NotTag) > 0 {
		state = true
	}
	if entityRule.NotNet != nil {
		state = true
	}
	if len(entityRule.NotSelector) > 0 {
		state = true
	}
	if len(entityRule.NotPorts) > 0 {
		state = true
	}

	return state
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

	d.Set("name", profile.Metadata.Name)
	d.Set("labels", profile.Metadata.Labels)

	setSchemaFieldsForProfileSpec(profile, d)
	return nil
}

func setResourceRulesFromProfile(profileRules []api.Rule) []map[string]interface{} {
	resourceRules := make([]map[string]interface{}, len(profileRules))

	for i, rule := range profileRules {
		resourceRule := make(map[string]interface{})
		if len(rule.Action) > 0 {
			resourceRule["action"] = rule.Action
		}
		if rule.Protocol != nil {
			resourceRule["protocol"] = rule.Protocol.String()
		}
		if rule.ICMP != nil {
			resourceIcmpArray := make([]map[string]int, 1)
			resourceIcmpMap := make(map[string]int)

			resourceIcmpMap["code"] = *rule.ICMP.Code
			resourceIcmpMap["type"] = *rule.ICMP.Type

			resourceIcmpArray[0] = resourceIcmpMap
			resourceRule["icmp"] = resourceIcmpArray
		}
		if rule.NotICMP != nil {
			resourceIcmpArray := make([]map[string]int, 1)
			resourceIcmpMap := make(map[string]int)

			resourceIcmpMap["code"] = *rule.ICMP.Code
			resourceIcmpMap["type"] = *rule.ICMP.Type

			resourceIcmpArray[0] = resourceIcmpMap
			resourceRule["notICMP"] = resourceIcmpArray
		}
		if nonEmptyEntityRule(&rule.Source) {
			resourceSourceArray := make([]map[string]interface{}, 1)

			resourceSourceArray[0] = createResourceEntityRule(rule.Source)
			resourceRule["source"] = resourceSourceArray
		}
		if nonEmptyEntityRule(&rule.Destination) {
			resourceSourceArray := make([]map[string]interface{}, 1)

			resourceSourceArray[0] = createResourceEntityRule(rule.Source)
			resourceRule["destination"] = resourceSourceArray
		}
		resourceRules[i] = resourceRule
	}

	return resourceRules

}
func createResourceEntityRule(entityRule api.EntityRule) map[string]interface{} {
	resourceSourceMap := make(map[string]interface{})

	if entityRule.Net != nil {
		resourceSourceMap["net"] = entityRule.Net.String()
	}
	if len(entityRule.Selector) > 0 {
		resourceSourceMap["selector"] = entityRule.Selector
	}
	if len(entityRule.Ports) > 0 {
		portsArray := make([]string, len(entityRule.Ports))
		for i, v := range entityRule.Ports {
			val := v.String()
			portsArray[i] = val
		}
		resourceSourceMap["ports"] = portsArray
	}
	if entityRule.NotNet != nil {
		resourceSourceMap["notNet"] = entityRule.NotNet.String()
	}
	if len(entityRule.NotSelector) > 0 {
		resourceSourceMap["notSelector"] = entityRule.NotSelector
	}
	if len(entityRule.NotPorts) > 0 {
		notPortsArray := make([]string, len(entityRule.NotPorts))
		for i, v := range entityRule.NotPorts {
			val := v.String()
			notPortsArray[i] = val
		}
		resourceSourceMap["notPorts"] = notPortsArray
	}
	return resourceSourceMap
}

func setSchemaFieldsForProfileSpec(profile *api.Profile, d *schema.ResourceData) {

	// Structure for this is pretty convulted. e.g.
	//   spec.0.ingress.0.rule.0.source.0.ports.0
	specArray := make([]interface{}, 2)
	if profile.Spec.IngressRules != nil && len(profile.Spec.IngressRules) > 0 {
		resourceRules := setResourceRulesFromProfile(profile.Spec.IngressRules)
		ruleMap := make(map[string]interface{})
		ruleMap["rule"] = resourceRules
		ruleMapArray := make([]interface{}, 1)
		ruleMapArray[0] = ruleMap
		ingressMap := make(map[string]interface{})
		ingressMap["ingress"] = ruleMapArray
		specArray[0] = ingressMap
	}
	if profile.Spec.EgressRules != nil && len(profile.Spec.EgressRules) > 0 {
		resourceRules := setResourceRulesFromProfile(profile.Spec.EgressRules)
		ruleMap := make(map[string]interface{})
		ruleMap["rule"] = resourceRules
		ruleMapArray := make([]interface{}, 1)
		ruleMapArray[0] = ruleMap
		egressMap := make(map[string]interface{})
		egressMap["egress"] = ruleMapArray
		specArray[1] = egressMap
	}
	d.Set("spec", specArray)
}

func resourceCalicoProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(config)
	calicoClient := config.Client

	profiles := calicoClient.Profiles()

	// Handle non-existent resource
	metadata := dToProfileMetadata(d)
	if _, err := profiles.Get(metadata); err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		}
	}

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
