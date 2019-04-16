package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/uswitch/fed-convert/pkg/converter"
	"github.com/uswitch/fed-convert/pkg/reader"
	"github.com/uswitch/fed-convert/pkg/writer"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfig string
	file       string
	fileOut    string
	clusters   string
)

func main() {
	flag.StringVar(&kubeConfig, "kube-config", "", "Path to kubeconfig file")
	flag.StringVar(&file, "file", "", "Path to file to be converted")
	flag.StringVar(&fileOut, "ouput-file", "", "Output file path, defaults to filename.out")
	flag.StringVar(&clusters, "clusters", "blue,red,black", "clusters to deploy to")
	flag.Parse()

	if fileOut == "" {
		fileOut = fmt.Sprintf("%s.out", file)
	}

	clusterArray := strings.Split(clusters, ",")

	objects, err := reader.ReadFile(file)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	hostConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("error creating kube config: %v", err)
	}

	fedResources, err := converter.Convert(hostConfig, clusterArray, objects)
	if err != nil {
		log.Fatalf("error converting resources: %v", err)
	}

	err = writer.WriteFile(fileOut, fedResources)
	if err != nil {
		log.Fatalf("error writing resources: %v", err)
	}

}
