package kube

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clients = make(map[string]*Client)

func GetClient(kubeconfig string) (*Client, error) {
	if client := clients[kubeconfig]; client != nil {
		return client, nil
	}
	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("get cluster state: load kubeconfig %q: %w", kubeconfig, err)
	}
	kc, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err)
	}
	clients[kubeconfig] = &Client{kc}
	return clients[kubeconfig], nil
}

type PodOptions struct {
	Namespace string
}

type Client struct {
	kclient *kubernetes.Clientset
}

func (c *Client) Pods(ctx context.Context, ns string) (*apiv1.PodList, error) {
	podsClient := c.kclient.CoreV1().Pods(ns)
	return podsClient.List(ctx, metav1.ListOptions{})
}
