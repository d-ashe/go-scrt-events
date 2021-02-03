package config

type Configurations struct {
	Node    NodeConfigurations
	Database DataBaseConfigurations
}

type NodeConfigurations struct {
	Host string
	Path string
}

type DataBaseConfigurations struct {
	Conn string
}