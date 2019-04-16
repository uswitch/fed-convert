package writer

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml2 "github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//WriteFile writes the federated resources to a file
func WriteFile(output string, resources []*unstructured.Unstructured) error {

	var bytes []byte
	for _, resource := range resources {
		object, err := yaml2.Marshal(resource)
		if err != nil {
			return fmt.Errorf("error marshalling yaml: %v", err)
		}
		dashes := []byte("---\n")
		object = append(dashes, object...)
		bytes = append(bytes, object...)
	}

	err := ioutil.WriteFile(output, bytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println(string(bytes))

	return nil
}
