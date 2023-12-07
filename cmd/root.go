package cmd

import (
	"context"
	"flag"
)

var (
	configure string
	delete    bool
	printer   bool
	cluster   map[string][]string
	conf      *config
	ctx       context.Context
)

type config struct {
	Redis struct {
		Addr     string   `yaml:"Addr"`
		Password string   `yaml:"Password"`
		DB       int      `yaml:"DB"`
		Match    []string `yaml:"Match"`
	} `yaml:"Redis"`
	Cluster struct {
		Addrs    []string `yaml:"Addrs"`
		Username string   `yaml:"Username"`
		Password string   `yaml:"Password"`
		Match    []string `yaml:"Match"`
	} `yaml:"Cluster"`
}

func init() {
	flag.StringVar(&configure, "conf", "config.yml", "指定配置文件")
	flag.BoolVar(&delete, "delete", false, "使用delete删除匹配key")
	flag.BoolVar(&printer, "print", false, "控制台输出要删除的key")
	cluster = map[string][]string{}
	ctx = context.Background()
	flag.Parse()
	loadconfiguration()
}

func Execute() {
	if conf.Redis.Addr != "" {
		deleteKeySingleInstance()
	}
	if len(conf.Cluster.Addrs) != 0 {
		deleteKeyCluster()
	}
}
