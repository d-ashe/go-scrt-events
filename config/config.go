package config

type Configurations struct {
	Node    NodeConfigurations
	Database DataBaseConfigurations
	Enrichment EnrichmentConfigurations
}

type NodeConfigurations struct {
	Host string
	Path string
}

type DataBaseConfigurations struct {
	Conn string
}

type EnrichmentConfigurations struct {
	Run []string
}