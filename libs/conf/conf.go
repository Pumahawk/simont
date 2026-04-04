package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pumahawk/simont/libs/core"
)

var confPath *string = flag.String("config", "", "Configuration path. Default: $HOME/.simont")

func ConfigPath() string {
	return *confPath
}

type clusterj struct {
	Name       string
	ConfigPath string
}

type namespacej struct {
	Cluster     string
	Name        string
	IsAuthority bool
}

type appconfj struct {
	Clusters   []clusterj
	Namespaces []namespacej
}

type AppConfig struct{ appconfj }

func (a *AppConfig) Clusters() []core.Cluster {
	nsgroup := make(map[string][]core.Namespace)
	for _, ns := range a.Namespaces {
		nsgroup[ns.Cluster] = append(nsgroup[ns.Cluster], core.Namespace{
			Name:        ns.Name,
			IsAuthority: ns.IsAuthority,
		})
	}
	cls := make([]core.Cluster, 0, len(a.appconfj.Clusters))
	for _, cl := range a.appconfj.Clusters {
		cls = append(cls, core.Cluster{
			Name:       cl.Name,
			ConfigPath: cl.ConfigPath,
			Namespaces: nsgroup[cl.Name],
		})
	}
	return cls
}

func LoadPath() (string, error) {
	confPath := *confPath
	if confPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("load path conf: %w", err)
		}
		return filepath.Join(home, ".simont.json"), nil
	} else {
		return confPath, nil
	}
}

func LoadConf() (*AppConfig, error) {
	path, err := LoadPath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", path, err)
	}
	defer file.Close()
	var confj appconfj
	err = json.NewDecoder(file).Decode(&confj)
	if err != nil {
		return nil, fmt.Errorf("decode config file %q: %w", path, err)
	}
	return &AppConfig{confj}, nil
}
