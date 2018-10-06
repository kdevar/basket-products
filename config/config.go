package config

import (
	"flag"
	"github.com/olivere/elastic"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Env string

func (e *Env) Set(s string) error {
	*e = Env(s)
	return nil
}

func (e *Env) String() string {
	return string(*e)
}

const (
	DEV  Env = "dev"
	TEST Env = "test"
	PROD Env = "prod"
)

type Config struct {
	ReindexedClusterPath string `yaml:"ReindexedClusterPath"`
	BasketBaseApiPath    string `yaml:"BasketBaseApiPath"`
	TypeAheadContextPath string `yaml:"TypeAheadContextPath"`
	AreaContextPath      string `yaml:"AreaContextPath"`
	TypeAheadToken       string `yaml:"TypeAheadToken"`
	Env                  string `yaml:"Env"`
}

type EnvConfig struct {
	Dev  Config
	Test Config
	Prod Config
}

func (ec *EnvConfig) GetConfig(e Env) Config {
	switch e {
	case DEV:
		return ec.Dev
	case TEST:
		return ec.Test
	case PROD:
		return ec.Prod
	default:
		return ec.Dev
	}
}

func NewConfig() *Config {
	yamlFile, err := ioutil.ReadFile("config.yml")

	var env Env
	flag.Var(&env, "env", "what is the env")

	flag.Parse()

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var envConfig EnvConfig
	err = yaml.Unmarshal(yamlFile, &envConfig)

	if err != nil {
		log.Printf("yamlFile unmarshal err   #%v \n", err)
	}

	currentConfig := envConfig.GetConfig(env)

	log.Printf("Config constructed env=%v config=%v \n", env, currentConfig)

	return &currentConfig
}

func NewElasticClient(c *Config) *elastic.Client {
	client, _ := elastic.NewSimpleClient(elastic.SetURL(c.ReindexedClusterPath))
	log.Println("Elastic Client constructed")
	return client
}
