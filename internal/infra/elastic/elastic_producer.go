package elastic

import (
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"os"
)

func NewESClient() (*elasticsearch7.Client, error) {
	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Username: os.Getenv("ELASTICSEARCH_USERNAME"),
		Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
	})
	return es7, err
}
