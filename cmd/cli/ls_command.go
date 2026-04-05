package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pumahawk/simont/libs/conf"
	"github.com/pumahawk/simont/libs/core"
	"github.com/pumahawk/simont/libs/svc"
)

var LsCommand = &Command{
	Name:  "ls",
	Descr: "Retrieve status of all clusters",
	Run: func(c *Command, args []string) int {
		conf, err := conf.LoadConf()
		if err != nil {
			log.Fatalf("ERROR - invalid configuration: %s", err)
		}
		clusters := conf.Clusters()
		type res struct {
			cs  *core.ClusterState
			err error
		}
		ch := make(chan res)
		for _, c := range clusters {
			go func(c core.Cluster) {
				r, err := svc.GetClusterState(context.TODO(), &c)
				ch <- res{r, err}
			}(c)
		}
		gstate := state(true)
		for _, c := range clusters {
			r := <-ch
			if cs, err := r.cs, r.err; err != nil {
				log.Printf("ERROR - info from cluster %q: %s", c.Name, err)
			} else {
				for _, ns := range cs.NamespacesState {
					state := state(true)
					gstate = state
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
		if gstate {
			return 0
		} else {
			return 1
		}
	},
}

type state bool

func (s state) String() string {
	if s {
		return "[X]"
	} else {
		return "[_]"
	}
}
