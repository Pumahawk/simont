package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/pumahawk/simont/libs/conf"
	"github.com/pumahawk/simont/libs/core"
	"github.com/pumahawk/simont/libs/svc"
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
		if cs, err := svc.GetClusterState(context.TODO(), &c); err != nil {
			log.Printf("ERROR - info from cluster %q: %s", c.Name, err)
		} else {
			for _, ns := range cs.NamespacesState {
				state := state(true)
				for _, svc := range ns.Services {
					if svc.State != core.Ok {
						state = false
						break
					}
				}
				fmt.Printf("%s %s %s\n", state, cs.Name, ns.Name)
			}
		}
	}
}

type state bool

func (s state) String() string {
	if s {
		return "[X]"
	} else {
		return "[_]"
	}
}
