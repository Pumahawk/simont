package kube

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const kapplabel = "app.kubernetes.io/name"

var wclients = make(chan int, 1)
var clients = make(clientStore)

type clientStore map[string]*Client

func (c clientStore) Get(key string) *Client {
	wclients <- 1
	defer func() { <-wclients }()
	return clients[key]
}

func (c clientStore) Set(key string, client *Client) {
	wclients <- 1
	defer func() { <-wclients }()
	clients[key] = client
}

func GetClient(kubeconfig string) (*Client, error) {
	if client := clients.Get(kubeconfig); client != nil {
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
	client := &Client{kc}
	clients.Set(kubeconfig, client)
	return client, nil
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

func AppName(pod *apiv1.Pod) string {
	return pod.Labels[kapplabel]
}
