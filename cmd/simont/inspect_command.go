package main

import (
	"context"
	"fmt"

	"github.com/pumahawk/simont/libs/conf"
	"github.com/pumahawk/simont/libs/core"
	"github.com/pumahawk/simont/libs/svc"
)

var InspectCommand = &Command{
	Name:  "inspect",
	Descr: "retrieve details by cluster and namespace",
	Run: func(c *Command, args []string) int {
		if len(args) < 2 {
			fmt.Println("Invalid arguments. Required [cluster name] [namespace]")
			return 1
		}
		clusterName := args[0]
		namespace := args[1]
		ac, err := conf.LoadConf()
		if err != nil {
			fmt.Printf("inspect read config: %s\n", err)
			return 1
		}
		cluster := GetClusterByName(ac.Clusters(), clusterName)
		if cluster == nil {
			fmt.Printf("Not found cluster %q %q\n", clusterName, namespace)
			return 1
		}
		cstate, err := svc.GetClusterState(context.TODO(), cluster)
		if err != nil {
			fmt.Printf("inspect get state: %s\n", err)
			return 1
		}
		gstate := state(true)
		for _, nss := range cstate.NamespacesState {
			if nss.Name == namespace {
				for _, svc := range nss.Services {
					state := state(svc.State == core.Ok)
					gstate = gstate && state
					fmt.Printf("%s %s %s %s\n", state, cstate.Name, nss.Name, svc.Pod)
				}
			}
		}
		if gstate {
			return 0
		} else {
			return 1
		}
	},
}

func GetClusterByName(clusters []core.Cluster, name string) *core.Cluster {
	var cl *core.Cluster
	for _, c := range clusters {
		if c.Name == name {
			cl = &c
		}
	}
	return cl
}
