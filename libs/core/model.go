package core

type State int16

const (
	Ok State = iota
	Warning
	ConnectionError
)

type Cluster struct {
	Name       string
	ConfigPath string
	Namespaces []Namespace
}

type ClusterState struct {
	Cluster
	State           State
	NamespacesState []NamespaceState
}

type Namespace struct {
	Name        string
	IsAuthority bool
}

type NamespaceState struct {
	State State
	Namespace
	Services []Service
}

type Service struct {
	State
	Message string
	Name    string
	Pod     string
}
