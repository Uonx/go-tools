package es

import (
	"strings"

	"github.com/olivere/elastic/v7"
)

type EsOpts struct {
	Endpoint string
	Username string
	Password string
}

var esClient *elastic.Client

func Es() *elastic.Client {
	if esClient == nil {
		panic("es 模块没有初始化")
	}

	return esClient
}

func NewEsClient(esOpts *EsOpts) error {
	client, err := elastic.NewClient(elastic.SetBasicAuth(esOpts.Username, esOpts.Password),
		elastic.SetURL(strings.Split(esOpts.Endpoint, ";")...),
		elastic.SetSniff(false))
	if err != nil {
		return err
	}
	esClient = client
	return nil
}
