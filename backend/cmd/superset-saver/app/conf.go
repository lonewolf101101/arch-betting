package app

type conf struct {
	DSN       string `yaml:"dsn"`
	KafkaHost string `yaml:"kafka_host"`
	GroupID   string `yaml:"group_id"`
	Topic     string `yaml:"topic"`
}
