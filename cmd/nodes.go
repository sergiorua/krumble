package cmd

import (
	"log"
	"os"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func LoadKubeconf() *rest.Config {
	//kcfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	kcfg, err := buildConfigFromFlags(config.Kops.Name, kubeconfig)
	if err != nil {
		log.Printf("Cannot open kubectl config: %v\n", err)
		os.Exit(1)
	}

	return kcfg
}

func isPodRunning(podName string, namespace string) bool {
	if getPodStatus(podName, namespace) == "Running" {
		log.Printf("%s:%s is Running", namespace, podName)
		return true
	}
	return false
}

func isPodStarting(podName string, namespace string) bool {
	if getPodStatus(podName, namespace) == "Pending" {
		log.Printf("%s:%s is Pending", namespace, podName)
		return true
	}
	return false
}

func getPodStatus(podName string, namespace string) string {
	cfg := LoadKubeconf()
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return "Unknown"
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.GetName(), podName) {
			return string(pod.Status.Phase)
		}
	}

	return "Unknown"
}

func getNamespaces() []string {
	var nss []string
	if dryrun {
		return []string{"kube-system", "kube-public", "default"}
	}
	cfg := LoadKubeconf()
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nss
	}
	nameSpaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	for _, ns := range nameSpaces.Items {
		nss = append(nss, ns.Name)
	}
	if debug {
		log.Println(nss)
	}
	return nss
}

func getServiceAccount(serAccount string, namespace string) corev1.ServiceAccount {
	if debug {
		log.Printf("Trying to locate %s:%s", namespace, serAccount)
	}
	cfg := LoadKubeconf()
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return corev1.ServiceAccount{}
	}
	serviceAccounts, err := clientset.CoreV1().ServiceAccounts(namespace).List(metav1.ListOptions{})
	for _, serviceAccount := range serviceAccounts.Items {
		if serviceAccount.Name == serAccount {
			return serviceAccount
		}
	}
	return corev1.ServiceAccount{}
}

func KopsNodesUp() bool {
	var nodeCount int = 0
	var masterCount int = 0
	var timewait int = 0
	const timeout = 900

	cfg := LoadKubeconf()
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Waiting for kops to build %d nodes and %d masters\n", config.Kops.NodeCount, config.Kops.MasterCount)
	for {
		timewait += 10
		nodeCount = 0
		masterCount = 0
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			log.Printf("%v\n", err.Error())
			if timewait >= timeout {
				log.Printf("TIMEOUT\n")
				return false
			}
			time.Sleep(10 * time.Second)
		}
		for i := range nodes.Items {
			node := nodes.Items[i]
			for x := range node.Status.Conditions {
				if node.Status.Conditions[x].Type == "Ready" {
					if node.Status.Conditions[x].Status == "True" && node.ObjectMeta.Labels["kubernetes.io/role"] == "node" {
						nodeCount++
					}
					if node.Status.Conditions[x].Status == "True" && node.ObjectMeta.Labels["kubernetes.io/role"] == "master" {
						masterCount++
					}
				}
			}
		}
		if debug {
			log.Printf("HAVE: Nodes=%d, Masters=%d\n", nodeCount, masterCount)
			log.Printf("WANT: Nodes=%d, Masters=%d\n", config.Kops.NodeCount, config.Kops.MasterCount)
		}
		if nodeCount >= config.Kops.NodeCount && masterCount >= config.Kops.MasterCount {
			break
		}
		if timewait >= timeout {
			log.Printf("TIMEOUT\n")
			return false
		}
		log.Print(".")
		time.Sleep(10 * time.Second)
	}

	log.Printf("Build complete\n")
	return true
}
