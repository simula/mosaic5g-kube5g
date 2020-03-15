// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1.Mosaic5g":       schema_pkg_apis_mosaic5g_v1alpha1_Mosaic5g(ref),
		"github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1.Mosaic5gSpec":   schema_pkg_apis_mosaic5g_v1alpha1_Mosaic5gSpec(ref),
		"github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1.Mosaic5gStatus": schema_pkg_apis_mosaic5g_v1alpha1_Mosaic5gStatus(ref),
	}
}

func schema_pkg_apis_mosaic5g_v1alpha1_Mosaic5g(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Mosaic5g is the Schema for the mosaic5gs API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1.Mosaic5gSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1.Mosaic5gStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1.Mosaic5gSpec", "github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1.Mosaic5gStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_mosaic5g_v1alpha1_Mosaic5gSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Mosaic5gSpec defines the desired state of Mosaic5g",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_mosaic5g_v1alpha1_Mosaic5gStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Mosaic5gStatus defines the observed state of Mosaic5g",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
