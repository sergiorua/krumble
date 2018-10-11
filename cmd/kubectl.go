package cmd

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"strings"
)

func LoadKubeconf() *rest.Config {
	kcfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Cannot open kubectl config: %v\n", err)
		os.Exit(1)
	}

	return kcfg
}

func LoadKubeYamlFromUrl(url string) []byte {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(url)
	buf := new(bytes.Buffer)

	if err != nil {
		log.Println("Could not download from " + url)
		return make([]byte, 0)
	}
	defer response.Body.Close()
	buf.ReadFrom(response.Body)
	respByte := buf.Bytes()

	return respByte
}

func LoadKubeYamlFromFile(fpath string) []byte {
	dat, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Printf("Cannot open file %s: %v\n", fpath, err)
		return make([]byte, 0)
	}

	return dat
}

type kubeDocs struct {
	Deployments         []*v1.Deployment
	Pods                []*apiv1.Pod
	ClusterRoleBindings []*rbacv1beta1.ClusterRoleBinding
	ClusterRoles        []*rbacv1.ClusterRole
	ServiceAccounts     []*apiv1.ServiceAccount
}

func LoadKubeYaml(url string) kubeDocs {
	var kd kubeDocs
	var respByte = make([]byte, 5)

	if strings.HasPrefix(url, "http") {
		respByte = LoadKubeYamlFromUrl(url)
	} else {
		respByte = LoadKubeYamlFromFile(url)
	}
	fileAsString := string(respByte[:])
	sepYamlfiles := strings.Split(fileAsString, "---")
	for _, f := range sepYamlfiles {
		if f == "\n" || f == "" {
			continue
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode

		obj, _, err := decode([]byte(f), nil, nil)
		if err != nil {
			log.Printf("%#v", err)
		}

		switch x := obj.(type) {
		case *v1.Deployment:
			deployment := obj.(*v1.Deployment)
			kd.Deployments = append(kd.Deployments, deployment)
		case *rbacv1beta1.ClusterRoleBinding:
			krb := obj.(*rbacv1beta1.ClusterRoleBinding)
			kd.ClusterRoleBindings = append(kd.ClusterRoleBindings, krb)
		case *rbacv1.ClusterRole:
			cr := obj.(*rbacv1.ClusterRole)
			kd.ClusterRoles = append(kd.ClusterRoles, cr)
		case *apiv1.ServiceAccount:
			sa := obj.(*apiv1.ServiceAccount)
			kd.ServiceAccounts = append(kd.ServiceAccounts, sa)
		case *apiv1.Pod:
			p := obj.(*apiv1.Pod)
			kd.Pods = append(kd.Pods, p)
		default:
			log.Printf("Unknown type %T\n", x)
		}
	}

	return kd
}

func ProcessKubectl() {
	var kd kubeDocs
	kcfg := LoadKubeconf()

	clientset, err := kubernetes.NewForConfig(kcfg)
	if err != nil {
		log.Printf("Cannot create clientset: %v\n", err)
		return
	}
	for i := range config.Kubectl {
		log.Printf("Processing entry %v\n", config.Kubectl[i])
		kd = LoadKubeYaml(config.Kubectl[i].URL)

		/* install deployments */
		deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
		for d := range kd.Deployments {
			result, err := deploymentsClient.Create(kd.Deployments[d])
			if err != nil {
				log.Printf("Cannot create deployment: %v\n", err)
				return
			}
			log.Printf("%v\n", result)
		}

	}
}
