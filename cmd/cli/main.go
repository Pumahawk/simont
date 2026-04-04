package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/pumahawk/simont/libs/conf"
	"github.com/pumahawk/simont/libs/kube"
)

func main() {
	flag.Parse()
	log.Println("Load configuration")
	conf, err := conf.LoadConf()
	if err != nil {
		log.Fatalf("INFO - invalid configuration: %s", err)
	}
	clusters := conf.Clusters()
	for _, c := range clusters {
		log.Printf("Load configuration %q", c.Name)
		if cs, err := kube.GetClusterState(context.TODO(), &c); err != nil {
			log.Printf("ERROR - info from cluster %q: %s", c.Name, err)
		} else {
			for _, ns := range cs.NamespacesState {
				for _, service := range ns.Services {
					fmt.Printf("%s %s %s %s %s\n", cs.Name, ns.Name, service.Name, service.State, service.Message)
				}
			}
		}
	}
}
