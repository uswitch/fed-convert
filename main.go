package main

import (
	"flag"
	"fmt"
	"log"

	yaml2 "github.com/ghodss/yaml"
	"github.com/uswitch/fed-convert/pkg/converter"
	"github.com/uswitch/fed-convert/pkg/reader"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfig string
	file       string
)

func main() {
	flag.StringVar(&kubeConfig, "kube-config", "", "Path to kubeconfig file")
	flag.StringVar(&file, "file", "", "Path to file to be converted")
	flag.Parse()

	objects, err := reader.ReadFile(file)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	hostConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("error creating kube config: %v", err)
	}

	fedResources, err := converter.Convert(hostConfig, objects)
	if err != nil {
		log.Fatalf("error converting resources: %v", err)
	}

	bytes, err := yaml2.Marshal(fedResources)
	if err != nil {
		log.Fatalf("oh noes: %v", err)
	}
	fmt.Println(string(bytes))
}
