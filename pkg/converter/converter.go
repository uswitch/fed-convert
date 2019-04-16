package converter

import (
	"fmt"
	"strings"

	ctlutil "github.com/kubernetes-sigs/federation-v2/pkg/controller/util"
	"github.com/kubernetes-sigs/federation-v2/pkg/kubefed2/enable"
	"github.com/kubernetes-sigs/federation-v2/pkg/kubefed2/federate"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

//Convert turns objects into federated objects
func Convert(config *rest.Config, clusters []string, objects []runtime.Object) ([]*unstructured.Unstructured, error) {
	var fedResources []*unstructured.Unstructured

	for _, object := range objects {

		typeKind := pluralise(object)

		apiResource, err := enable.LookupAPIResource(config, typeKind, "")
		if err != nil {
			return nil, fmt.Errorf("error looking up api resource: %v", err)
		}

		typeConfig := enable.GenerateTypeConfigForTarget(*apiResource, enable.NewEnableTypeDirective())
		fedResource, err := federate.FederatedResourceFromTargetResource(typeConfig, object.(*unstructured.Unstructured))
		if err != nil {
			return nil, fmt.Errorf("error generating federated resource: %v", err)
		}

		clusterList := make([]interface{}, len(clusters))
		for i, v := range clusters {
			clusterList[i] = v
		}
		err = unstructured.SetNestedSlice(fedResource.Object, clusterList, ctlutil.SpecField, ctlutil.PlacementField, ctlutil.ClusterNamesField)
		if err != nil {
			return nil, fmt.Errorf("error generating federated resource: %v", err)
		}
		unstructured.RemoveNestedField(fedResource.Object, ctlutil.SpecField, ctlutil.PlacementField, ctlutil.ClusterSelectorField)
		fedResources = append(fedResources, fedResource)
	}

	return fedResources, nil
}

//This is because you need to have the plural kind and group when referring to objects such as Deployment which has both apps and extensions as its group
func pluralise(object runtime.Object) string {
	kind := object.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%ss.%s", strings.ToLower(kind.Kind), kind.Group)
}
