package kube

import (
	"context"
	"fmt"
	"slices"

	"github.com/pumahawk/simont/libs/core"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clients = make(map[string]*kubernetes.Clientset)

func GetClusterState(ctx context.Context, c *core.Cluster) (*core.ClusterState, error) {
	client, err := getClient(c.ConfigPath)
	if err != nil {
		return nil, err
	}
	podsClient := client.CoreV1().Pods(apiv1.NamespaceAll)
	podsns := make(map[string][]*apiv1.Pod)
	if pods, err := podsClient.List(ctx, metav1.ListOptions{}); err != nil {
		return nil, fmt.Errorf("get cluster state %q: %w", c.Name, err)
	} else {
		for _, pod := range pods.Items {
			if slices.ContainsFunc(c.Namespaces, func(ns core.Namespace) bool {
				return ns.Name == pod.Namespace
			}) {
				podsns[pod.Namespace] = append(podsns[pod.Namespace], &pod)
			}
		}
	}

	nsss := make([]core.NamespaceState, 0, len(c.Namespaces))
	for _, ns := range c.Namespaces {
		nss := core.NamespaceState{
			Namespace: ns,
		}
		for _, pod := range podsns[ns.Name] {
			service := core.Service{
				Name: pod.Name,
				Pod:  pod.Name,
			}
			service.State = core.Ok
			for _, st := range pod.Status.ContainerStatuses {
				if !st.Ready {
					nss.State = core.Warning
					service.State = core.Warning
				}
			}
			nss.Services = append(nss.Services, service)
		}
		nsss = append(nsss, nss)
	}
	return &core.ClusterState{
		State:           core.Ok,
		Cluster:         *c,
		NamespacesState: nsss,
	}, nil
}

func getClient(kubeconfig string) (*kubernetes.Clientset, error) {
	client := clients[kubeconfig]
	if client != nil {
		return client, nil
	}
	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("get cluster state: load kubeconfig %q: %w", kubeconfig, err)
	}
	client, err = kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err)
	}
	clients[kubeconfig] = client
	return client, nil
}
