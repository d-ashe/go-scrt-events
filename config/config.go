package config

type Configurations struct {
	Node    NodeConfigurations
	Database DataBaseConfigurations
	Enrichment EnrichmentConfigurations
	Pubsub PubSubConfigurations
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

type PubSubConfigurations struct {
	ProjectId string
	TopicName string
}