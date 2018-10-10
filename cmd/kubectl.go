package cmd

import (
	"crypto/tls"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/kubernetes/pkg/api"
	_ "k8s.io/kubernetes/pkg/api/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"

	"log"
	"net/http"
	"os"
)

func LoadKubeconf() *rest.Config {
	kcfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Cannot open kubectl config: %v\n", err)
		os.Exit(1)
	}

	return kcfg
}

func LoadKubeYaml(url string) *yaml.YAMLOrJSONDecoder {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(url)
	if err != nil {
		log.Println("Could not download from " + url)
		return nil
	}
	defer response.Body.Close()
	decode := api.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(json), nil, nil)
	if err != nil {
		log.Printf("%#v", err)
	}
	deployment := obj.(*v1beta1.Deployment)

	return d
}

func ProcessKubectl() {
	kcfg := LoadKubeconf()
	var deployment *v1beta1.Deployment

	for i := range config.Kubectl {
		log.Printf("Processing entry %v\n", config.Kubectl[i])
		d := LoadKubeYaml(config.Kubectl[i].URL)

		log.Printf("%v\n", d)
		_ = d.Decode(&deployment)
		log.Printf("%#v\n", deployment)

		clientset, err := kubernetes.NewForConfig(kcfg)
		if err != nil {
			log.Printf("Error creating kubernetes clientset: %v\n", err)
			return
		}

		dd := clientset.Discovery()
		log.Println(dd)
	}
}
