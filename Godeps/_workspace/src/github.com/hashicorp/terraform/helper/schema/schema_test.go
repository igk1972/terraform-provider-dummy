package schema

import (
	"reflect"
	"testing"

	"github.com/whitepages/terraform-provider-dummy/Godeps/_workspace/src/github.com/hashicorp/terraform/config"
	"github.com/whitepages/terraform-provider-dummy/Godeps/_workspace/src/github.com/hashicorp/terraform/terraform"
)

func TestSchemaMap_Diff(t *testing.T) {
	cases := []struct {
		Schema          map[string]*Schema
		State           *terraform.InstanceState
		Config          map[string]interface{}
		ConfigVariables map[string]string
		Diff            *terraform.InstanceDiff
		Err             bool
	}{
		/*
		 * String decode
		 */

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"availability_zone": "foo",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Config: map[string]interface{}{},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						NewComputed: true,
						RequiresNew: true,
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.InstanceState{
				ID: "foo",
			},

			Config: map[string]interface{}{},

			Diff: nil,

			Err: false,
		},

		// Computed, but set in config
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Config: map[string]interface{}{
				"availability_zone": "bar",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old: "foo",
						New: "bar",
					},
				},
			},

			Err: false,
		},

		// Default
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Default:  "foo",
				},
			},

			State: nil,

			Config: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old: "",
						New: "foo",
					},
				},
			},

			Err: false,
		},

		// DefaultFunc, value
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					DefaultFunc: func() (interface{}, error) {
						return "foo", nil
					},
				},
			},

			State: nil,

			Config: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old: "",
						New: "foo",
					},
				},
			},

			Err: false,
		},

		// DefaultFunc, configuration set
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					DefaultFunc: func() (interface{}, error) {
						return "foo", nil
					},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"availability_zone": "bar",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},
				},
			},

			Err: false,
		},

		// String with StateFunc
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					StateFunc: func(a interface{}) string {
						return a.(string) + "!"
					},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"availability_zone": "foo",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:      "",
						New:      "foo!",
						NewExtra: "foo",
					},
				},
			},

			Err: false,
		},

		// Variable (just checking)
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"availability_zone": "${var.foo}",
			},

			ConfigVariables: map[string]string{
				"var.foo": "bar",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},
				},
			},

			Err: false,
		},

		// Variable computed
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"availability_zone": "${var.foo}",
			},

			ConfigVariables: map[string]string{
				"var.foo": config.UnknownVariableValue,
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old: "",
						New: "${var.foo}",
					},
				},
			},

			Err: false,
		},

		/*
		 * Int decode
		 */

		{
			Schema: map[string]*Schema{
				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"port": 27,
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"port": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "27",
						RequiresNew: true,
					},
				},
			},

			Err: false,
		},

		/*
		 * Bool decode
		 */

		{
			Schema: map[string]*Schema{
				"port": &Schema{
					Type:     TypeBool,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"port": false,
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"port": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "0",
						RequiresNew: true,
					},
				},
			},

			Err: false,
		},

		/*
		 * Bool
		 */
		{
			Schema: map[string]*Schema{
				"delete": &Schema{
					Type:     TypeBool,
					Optional: true,
					Default:  false,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"delete": "false",
				},
			},

			Config: nil,

			Diff: nil,

			Err: false,
		},

		/*
		 * List decode
		 */

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ports": []interface{}{1, 2, 5},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "0",
						New: "3",
					},
					"ports.0": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "",
						New: "2",
					},
					"ports.2": &terraform.ResourceAttrDiff{
						Old: "",
						New: "5",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ports": []interface{}{1, "${var.foo}"},
			},

			ConfigVariables: map[string]string{
				"var.foo": "2" + config.InterpSplitDelim + "5",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "0",
						New: "3",
					},
					"ports.0": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "",
						New: "2",
					},
					"ports.2": &terraform.ResourceAttrDiff{
						Old: "",
						New: "5",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ports": []interface{}{1, "${var.foo}"},
			},

			ConfigVariables: map[string]string{
				"var.foo": config.UnknownVariableValue +
					config.InterpSplitDelim + "5",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "0",
						New:         "",
						NewComputed: true,
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Config: map[string]interface{}{
				"ports": []interface{}{1, 2, 5},
			},

			Diff: nil,

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "2",
					"ports.0": "1",
					"ports.1": "2",
				},
			},

			Config: map[string]interface{}{
				"ports": []interface{}{1, 2, 5},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "2",
						New: "3",
					},
					"ports.2": &terraform.ResourceAttrDiff{
						Old: "",
						New: "5",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					ForceNew: true,
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ports": []interface{}{1, 2, 5},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "0",
						New:         "3",
						RequiresNew: true,
					},
					"ports.0": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "1",
						RequiresNew: true,
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "2",
						RequiresNew: true,
					},
					"ports.2": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "5",
						RequiresNew: true,
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Config: map[string]interface{}{},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "",
						NewComputed: true,
					},
				},
			},

			Err: false,
		},

		/*
		 * Set
		 */

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ports": []interface{}{5, 2, 1},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "0",
						New: "3",
					},
					"ports.0": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "",
						New: "2",
					},
					"ports.2": &terraform.ResourceAttrDiff{
						Old: "",
						New: "5",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Computed: true,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "0",
				},
			},

			Config: nil,

			Diff: nil,

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: nil,

			Config: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "",
						NewComputed: true,
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ports": []interface{}{"${var.foo}", 1},
			},

			ConfigVariables: map[string]string{
				"var.foo": "2" + config.InterpSplitDelim + "5",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "0",
						New: "3",
					},
					"ports.0": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "",
						New: "2",
					},
					"ports.2": &terraform.ResourceAttrDiff{
						Old: "",
						New: "5",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ports": []interface{}{1, "${var.foo}"},
			},

			ConfigVariables: map[string]string{
				"var.foo": config.UnknownVariableValue +
					config.InterpSplitDelim + "5",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "0",
						New:         "",
						NewComputed: true,
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "2",
					"ports.0": "2",
					"ports.1": "1",
				},
			},

			Config: map[string]interface{}{
				"ports": []interface{}{5, 2, 1},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "2",
						New: "3",
					},
					"ports.0": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "1",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "2",
						New: "2",
					},
					"ports.2": &terraform.ResourceAttrDiff{
						Old: "",
						New: "5",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "2",
					"ports.0": "2",
					"ports.1": "1",
				},
			},

			Config: map[string]interface{}{},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "2",
						New: "0",
					},
					"ports.0": &terraform.ResourceAttrDiff{
						Old:        "1",
						NewRemoved: true,
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old:        "2",
						NewRemoved: true,
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set:      func(v interface{}) int { return v.(int) },
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "bar",
					"ports.#":           "1",
					"ports.0":           "80",
				},
			},

			Config: map[string]interface{}{},

			Diff: nil,

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeSet,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"ports": &Schema{
								Type:     TypeList,
								Optional: true,
								Elem:     &Schema{Type: TypeInt},
							},
						},
					},
					Set: func(v interface{}) int {
						m := v.(map[string]interface{})
						ps := m["ports"].([]interface{})
						result := 0
						for _, p := range ps {
							result += p.(int)
						}
						return result
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ingress.#":         "2",
					"ingress.0.ports.#": "1",
					"ingress.0.ports.0": "80",
					"ingress.1.ports.#": "1",
					"ingress.1.ports.0": "443",
				},
			},

			Config: map[string]interface{}{
				"ingress": []interface{}{
					map[string]interface{}{
						"ports": []interface{}{443},
					},
					map[string]interface{}{
						"ports": []interface{}{80},
					},
				},
			},

			Diff: nil,

			Err: false,
		},

		/*
		 * List of structure decode
		 */

		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"ingress": []interface{}{
					map[string]interface{}{
						"from": 8080,
					},
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "0",
						New: "1",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "8080",
					},
				},
			},

			Err: false,
		},

		/*
		 * ComputedWhen
		 */

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:         TypeString,
					Computed:     true,
					ComputedWhen: []string{"port"},
				},

				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
					"port":              "80",
				},
			},

			Config: map[string]interface{}{
				"port": 80,
			},

			Diff: nil,

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:         TypeString,
					Computed:     true,
					ComputedWhen: []string{"port"},
				},

				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"port": "80",
				},
			},

			Config: map[string]interface{}{
				"port": 80,
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						NewComputed: true,
					},
				},
			},

			Err: false,
		},

		/* TODO
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:         TypeString,
					Computed:     true,
					ComputedWhen: []string{"port"},
				},

				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
					"port":              "80",
				},
			},

			Config: map[string]interface{}{
				"port": 8080,
			},

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "foo",
						NewComputed: true,
					},
					"port": &terraform.ResourceAttrDiff{
						Old: "80",
						New: "8080",
					},
				},
			},

			Err: false,
		},
		*/

		/*
		 * Maps
		 */

		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type: TypeMap,
				},
			},

			State: nil,

			Config: map[string]interface{}{
				"config_vars": []interface{}{
					map[string]interface{}{
						"bar": "baz",
					},
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.bar": &terraform.ResourceAttrDiff{
						Old: "",
						New: "baz",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type: TypeMap,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.foo": "bar",
				},
			},

			Config: map[string]interface{}{
				"config_vars": []interface{}{
					map[string]interface{}{
						"bar": "baz",
					},
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.foo": &terraform.ResourceAttrDiff{
						Old:        "bar",
						NewRemoved: true,
					},
					"config_vars.bar": &terraform.ResourceAttrDiff{
						Old: "",
						New: "baz",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"vars": &Schema{
					Type:     TypeMap,
					Optional: true,
					Computed: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"vars.foo": "bar",
				},
			},

			Config: map[string]interface{}{
				"vars": []interface{}{
					map[string]interface{}{
						"bar": "baz",
					},
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"vars.foo": &terraform.ResourceAttrDiff{
						Old:        "bar",
						New:        "",
						NewRemoved: true,
					},
					"vars.bar": &terraform.ResourceAttrDiff{
						Old: "",
						New: "baz",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"vars": &Schema{
					Type:     TypeMap,
					Computed: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"vars.foo": "bar",
				},
			},

			Config: nil,

			Diff: nil,

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type: TypeList,
					Elem: &Schema{Type: TypeMap},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "1",
					"config_vars.0.foo": "bar",
				},
			},

			Config: map[string]interface{}{
				"config_vars": []interface{}{
					map[string]interface{}{
						"bar": "baz",
					},
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.0.foo": &terraform.ResourceAttrDiff{
						Old:        "bar",
						NewRemoved: true,
					},
					"config_vars.0.bar": &terraform.ResourceAttrDiff{
						Old: "",
						New: "baz",
					},
				},
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type: TypeList,
					Elem: &Schema{Type: TypeMap},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "1",
					"config_vars.0.foo": "bar",
					"config_vars.0.bar": "baz",
				},
			},

			Config: map[string]interface{}{},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "0",
					},
					"config_vars.0.foo": &terraform.ResourceAttrDiff{
						Old:        "bar",
						NewRemoved: true,
					},
					"config_vars.0.bar": &terraform.ResourceAttrDiff{
						Old:        "baz",
						NewRemoved: true,
					},
				},
			},

			Err: false,
		},

		/*
		 * ForceNews
		 */

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					ForceNew: true,
				},

				"address": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "bar",
					"address":           "foo",
				},
			},

			Config: map[string]interface{}{
				"availability_zone": "foo",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "bar",
						New:         "foo",
						RequiresNew: true,
					},

					"address": &terraform.ResourceAttrDiff{
						Old:         "foo",
						New:         "",
						NewComputed: true,
					},
				},
			},

			Err: false,
		},

		// Set
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					ForceNew: true,
				},

				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set:      func(v interface{}) int { return v.(int) },
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "bar",
					"ports.#":           "1",
					"ports.0":           "80",
				},
			},

			Config: map[string]interface{}{
				"availability_zone": "foo",
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "bar",
						New:         "foo",
						RequiresNew: true,
					},

					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "1",
						New:         "",
						NewComputed: true,
					},
				},
			},

			Err: false,
		},
	}

	for i, tc := range cases {
		c, err := config.NewRawConfig(tc.Config)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if len(tc.ConfigVariables) > 0 {
			if err := c.Interpolate(tc.ConfigVariables); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		d, err := schemaMap(tc.Schema).Diff(
			tc.State, terraform.NewResourceConfig(c))
		if (err != nil) != tc.Err {
			t.Fatalf("#%d err: %s", i, err)
		}

		if !reflect.DeepEqual(tc.Diff, d) {
			t.Fatalf("#%d: bad:\n\n%#v", i, d)
		}
	}
}

func TestSchemaMap_Input(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		Config map[string]interface{}
		Input  map[string]string
		Result map[string]interface{}
		Err    bool
	}{
		/*
		 * String decode
		 */

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
				},
			},

			Input: map[string]string{
				"availability_zone": "foo",
			},

			Result: map[string]interface{}{
				"availability_zone": "foo",
			},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
				},
			},

			Config: map[string]interface{}{
				"availability_zone": "bar",
			},

			Input: map[string]string{
				"availability_zone": "foo",
			},

			Result: map[string]interface{}{},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Default:  "foo",
					Optional: true,
				},
			},

			Input: map[string]string{
				"availability_zone": "bar",
			},

			Result: map[string]interface{}{},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type: TypeString,
					DefaultFunc: func() (interface{}, error) {
						return "foo", nil
					},
					Optional: true,
				},
			},

			Input: map[string]string{
				"availability_zone": "bar",
			},

			Result: map[string]interface{}{},

			Err: false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type: TypeString,
					DefaultFunc: func() (interface{}, error) {
						return nil, nil
					},
					Optional: true,
				},
			},

			Input: map[string]string{
				"availability_zone": "bar",
			},

			Result: map[string]interface{}{
				"availability_zone": "bar",
			},

			Err: false,
		},
	}

	for i, tc := range cases {
		if tc.Config == nil {
			tc.Config = make(map[string]interface{})
		}

		c, err := config.NewRawConfig(tc.Config)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		input := new(terraform.MockUIInput)
		input.InputReturnMap = tc.Input

		rc := terraform.NewResourceConfig(c)
		rc.Config = make(map[string]interface{})

		actual, err := schemaMap(tc.Schema).Input(input, rc)
		if (err != nil) != tc.Err {
			t.Fatalf("#%d err: %s", i, err)
		}

		if !reflect.DeepEqual(tc.Result, actual.Config) {
			t.Fatalf("#%d: bad:\n\n%#v", i, actual.Config)
		}
	}
}

func TestSchemaMap_InternalValidate(t *testing.T) {
	cases := []struct {
		In  map[string]*Schema
		Err bool
	}{
		{
			nil,
			false,
		},

		// No optional and no required
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeInt,
					Optional: true,
					Required: true,
				},
			},
			true,
		},

		// No optional and no required
		{
			map[string]*Schema{
				"foo": &Schema{
					Type: TypeInt,
				},
			},
			true,
		},

		// Missing Type
		{
			map[string]*Schema{
				"foo": &Schema{
					Required: true,
				},
			},
			true,
		},

		// Required but computed
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeInt,
					Required: true,
					Computed: true,
				},
			},
			true,
		},

		// Looks good
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeString,
					Required: true,
				},
			},
			false,
		},

		// Computed but has default
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeInt,
					Optional: true,
					Computed: true,
					Default:  "foo",
				},
			},
			true,
		},

		// Required but has default
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeInt,
					Optional: true,
					Required: true,
					Default:  "foo",
				},
			},
			true,
		},

		// List element not set
		{
			map[string]*Schema{
				"foo": &Schema{
					Type: TypeList,
				},
			},
			true,
		},

		// List default
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:    TypeList,
					Elem:    &Schema{Type: TypeInt},
					Default: "foo",
				},
			},
			true,
		},

		// List element computed
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeList,
					Optional: true,
					Elem: &Schema{
						Type:     TypeInt,
						Computed: true,
					},
				},
			},
			true,
		},

		// List element with Set set
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeList,
					Elem:     &Schema{Type: TypeInt},
					Set:      func(interface{}) int { return 0 },
					Optional: true,
				},
			},
			true,
		},

		// Set element with no Set set
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeSet,
					Elem:     &Schema{Type: TypeInt},
					Optional: true,
				},
			},
			true,
		},

		// Required but computed
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:         TypeInt,
					Required:     true,
					ComputedWhen: []string{"foo"},
				},
			},
			true,
		},

		// Sub-resource invalid
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeList,
					Optional: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"foo": new(Schema),
						},
					},
				},
			},
			true,
		},

		// Sub-resource valid
		{
			map[string]*Schema{
				"foo": &Schema{
					Type:     TypeList,
					Optional: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"foo": &Schema{
								Type:     TypeInt,
								Optional: true,
							},
						},
					},
				},
			},
			false,
		},
	}

	for i, tc := range cases {
		err := schemaMap(tc.In).InternalValidate()
		if (err != nil) != tc.Err {
			t.Fatalf("%d: bad: %s\n\n%#v", i, err, tc.In)
		}
	}

}

func TestSchemaMap_Validate(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		Config map[string]interface{}
		Vars   map[string]string
		Warn   bool
		Err    bool
	}{
		// Good
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			Config: map[string]interface{}{
				"availability_zone": "foo",
			},
		},

		// Good, because the var is not set and that error will come elsewhere
		{
			Schema: map[string]*Schema{
				"size": &Schema{
					Type:     TypeInt,
					Required: true,
				},
			},

			Config: map[string]interface{}{
				"size": "${var.foo}",
			},

			Vars: map[string]string{
				"var.foo": config.UnknownVariableValue,
			},
		},

		// Required field not set
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Required: true,
				},
			},

			Config: map[string]interface{}{},

			Err: true,
		},

		// Invalid type
		{
			Schema: map[string]*Schema{
				"port": &Schema{
					Type:     TypeInt,
					Required: true,
				},
			},

			Config: map[string]interface{}{
				"port": "I am invalid",
			},

			Err: true,
		},

		{
			Schema: map[string]*Schema{
				"user_data": &Schema{
					Type:     TypeString,
					Optional: true,
				},
			},

			Config: map[string]interface{}{
				"user_data": []interface{}{
					map[string]interface{}{
						"foo": "bar",
					},
				},
			},

			Err: true,
		},

		// Bad type, interpolated
		{
			Schema: map[string]*Schema{
				"size": &Schema{
					Type:     TypeInt,
					Required: true,
				},
			},

			Config: map[string]interface{}{
				"size": "${var.foo}",
			},

			Vars: map[string]string{
				"var.foo": "nope",
			},

			Err: true,
		},

		// Required but has DefaultFunc
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Required: true,
					DefaultFunc: func() (interface{}, error) {
						return "foo", nil
					},
				},
			},

			Config: nil,
		},

		// Required but has DefaultFunc return nil
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Required: true,
					DefaultFunc: func() (interface{}, error) {
						return nil, nil
					},
				},
			},

			Config: nil,

			Err: true,
		},

		// Optional sub-resource
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type: TypeList,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			Config: map[string]interface{}{},

			Err: false,
		},

		// Not a list
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type: TypeList,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			Config: map[string]interface{}{
				"ingress": "foo",
			},

			Err: true,
		},

		// Required sub-resource field
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type: TypeList,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			Config: map[string]interface{}{
				"ingress": []interface{}{
					map[string]interface{}{},
				},
			},

			Err: true,
		},

		// Good sub-resource
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Optional: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			Config: map[string]interface{}{
				"ingress": []interface{}{
					map[string]interface{}{
						"from": 80,
					},
				},
			},

			Err: false,
		},

		// Invalid/unknown field
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			Config: map[string]interface{}{
				"foo": "bar",
			},

			Err: true,
		},

		// Computed field set
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Computed: true,
				},
			},

			Config: map[string]interface{}{
				"availability_zone": "bar",
			},

			Err: true,
		},

		// Not a set
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			Config: map[string]interface{}{
				"ports": "foo",
			},

			Err: true,
		},

		// Maps
		{
			Schema: map[string]*Schema{
				"user_data": &Schema{
					Type:     TypeMap,
					Optional: true,
				},
			},

			Config: map[string]interface{}{
				"user_data": "foo",
			},

			Err: true,
		},

		{
			Schema: map[string]*Schema{
				"user_data": &Schema{
					Type:     TypeMap,
					Optional: true,
				},
			},

			Config: map[string]interface{}{
				"user_data": []interface{}{
					map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		},

		{
			Schema: map[string]*Schema{
				"user_data": &Schema{
					Type:     TypeMap,
					Optional: true,
				},
			},

			Config: map[string]interface{}{
				"user_data": map[string]interface{}{
					"foo": "bar",
				},
			},
		},

		{
			Schema: map[string]*Schema{
				"user_data": &Schema{
					Type:     TypeMap,
					Optional: true,
				},
			},

			Config: map[string]interface{}{
				"user_data": []interface{}{
					"foo",
				},
			},

			Err: true,
		},
	}

	for i, tc := range cases {
		c, err := config.NewRawConfig(tc.Config)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if tc.Vars != nil {
			if err := c.Interpolate(tc.Vars); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		ws, es := schemaMap(tc.Schema).Validate(terraform.NewResourceConfig(c))
		if (len(es) > 0) != tc.Err {
			if len(es) == 0 {
				t.Errorf("%d: no errors", i)
			}

			for _, e := range es {
				t.Errorf("%d: err: %s", i, e)
			}

			t.FailNow()
		}

		if (len(ws) > 0) != tc.Warn {
			t.Fatalf("%d: ws: %#v", i, ws)
		}
	}
}
