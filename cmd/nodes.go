package cmd

import (
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func LoadKubeconf() *rest.Config {
	kcfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Cannot open kubectl config: %v\n", err)
		os.Exit(1)
	}

	return kcfg
}

func KopsNodesUp() bool {
	var node_count int = 0
	var master_count int = 0
	var timewait int = 0
	const timeout = 300

	cfg := LoadKubeconf()
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	for {
		timewait += 10
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			log.Printf("%v\n", err.Error())
			if timewait >= timeout {
				log.Printf("TIMEOUT\n")
				return false
			}
			time.Sleep(10)
		}
		for i := range nodes.Items {
			node := nodes.Items[i]
			for x := range node.Status.Conditions {
				if node.Status.Conditions[x].Type == "Ready" {
					if node.Status.Conditions[x].Status == "True" {
						if node.ObjectMeta.Labels["kubernetes.io/role"] == "node" {
							node_count++
						}
						if node.ObjectMeta.Labels["kubernetes.io/role"] == "master" {
							master_count++
						}
					}
				}
			}
		}
		log.Printf("HAVE: Nodes=%d, Masters=%d\n", node_count, master_count)
		log.Printf("WANT: Nodes=%d, Masters=%d\n", config.Kops.NodeCount, config.Kops.MasterCount)
		if node_count == config.Kops.NodeCount && master_count == config.Kops.MasterCount {
			break
		}
		if timewait >= timeout {
			log.Printf("TIMEOUT\n")
			return false
		}
		time.Sleep(10)
	}

	log.Printf("Nodes=%d, Masters=%d\n", node_count, master_count)
	return true
}
