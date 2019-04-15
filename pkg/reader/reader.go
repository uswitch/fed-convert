package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/streaming"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

//ReadFile reads file of kubernetes resources and returns them as unstructured objects
func ReadFile(file string) ([]runtime.Object, error) {
	// Build a yaml decoder with the unstructured Scheme
	yamlDecoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	bytess, err := ioutil.ReadFile("deployment.yaml")
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	}
	// Parse the objects from the yaml
	var objects []runtime.Object
	reader := json.YAMLFramer.NewFrameReader(ioutil.NopCloser(bytes.NewReader(bytess)))
	d := streaming.NewDecoder(reader, yamlDecoder)
	for {
		obj, _, err := d.Decode(nil, nil)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error during parse: %v", err)
		}
		objects = append(objects, obj)
	}
	return objects, nil
}
