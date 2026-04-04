package core

type State int16

const (
	Ok State = iota
	Error
)

type Cluster struct {
	Name       string
	ConfigPath string
	Namespaces []Namespace
}

type ClusterState struct {
	Cluster
	State
}

type Namespace struct {
	Name        string
	IsAuthority bool
}

type NamespaceState struct {
}
