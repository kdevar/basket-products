package config

import "github.com/olivere/elastic"

var ElasticClient, e = elastic.NewSimpleClient(elastic.SetURL("https://vpc-reindex-nxvfoonqh3jbcz37uu6b4zfov4.us-east-1.es.amazonaws.com:443/"))