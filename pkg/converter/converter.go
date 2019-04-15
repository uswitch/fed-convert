package converter

import (
	"fmt"

	"github.com/kubernetes-sigs/federation-v2/pkg/kubefed2/enable"
	"github.com/kubernetes-sigs/federation-v2/pkg/kubefed2/federate"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

//Convert turns objects into federated objects
func Convert(config *rest.Config, objects []runtime.Object) ([]*unstructured.Unstructured, error) {
	var fedResources []*unstructured.Unstructured

	for _, object := range objects {

		typeKind := object.GetObjectKind().GroupVersionKind().Kind
		if typeKind == "Deployment" {
			typeKind = "deployments.apps"
		}

		apiResource, err := enable.LookupAPIResource(config, typeKind, "")
		if err != nil {
			return nil, fmt.Errorf("error looking up api resource: %v", err)
		}

		typeConfig := enable.GenerateTypeConfigForTarget(*apiResource, enable.NewEnableTypeDirective())
		fedResource, err := federate.FederatedResourceFromTargetResource(typeConfig, object.(*unstructured.Unstructured))
		if err != nil {
			return nil, fmt.Errorf("error generating federated resource: %v", err)
		}
		fedResources = append(fedResources, fedResource)
	}

	return fedResources, nil
}
