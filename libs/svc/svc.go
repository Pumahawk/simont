package svc

import (
	"context"
	"fmt"
	"slices"

	"github.com/pumahawk/simont/libs/core"
	"github.com/pumahawk/simont/libs/kube"
	apiv1 "k8s.io/api/core/v1"
)

func GetClusterState(ctx context.Context, c *core.Cluster) (*core.ClusterState, error) {
	client, err := kube.GetClient(c.ConfigPath)
	if err != nil {
		return nil, err
	}

	podsns := make(map[string][]*apiv1.Pod)
	if pods, err := client.Pods(ctx, ""); err != nil {
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
