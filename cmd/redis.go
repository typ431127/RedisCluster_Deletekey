package cmd

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func loadconfiguration() {
	configuration, err := os.ReadFile(configure)
	if err != nil {
		log.Fatalf("Failed to read configuration file %v!", err)
	}
	err = yaml.Unmarshal(configuration, &conf)
	if err != nil {
		log.Fatalf("Yaml configuration file error: %v", err)
	}
}

func instanceDelete(shard *redis.Client, match string) {
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = shard.Scan(ctx, cursor, match, 100).Result()
		if err != nil {
			log.Fatalln(err)
		}
		if printer {
			for _, key := range keys {
				log.Printf("%s Match-key:\033[1;31;40m%s\033[0m\n", shard.Options().Addr, key)
			}
		}
		if delete {
			for _, key := range keys {
				log.Printf("%s delete-key:\033[1;31;40m%s\033[0m\n", shard.Options().Addr, key)
				shard.Del(ctx, key)
			}
		}
		n += len(keys)
		if cursor == 0 {
			break
		}
	}
	log.Printf("%s found \033[1;31;1m%d\033[0m keys\n", shard.Options().Addr, n)
}
func deleteKeySingleInstance() {
	log.Printf("Redis-Instance :%s", conf.Redis.Addr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	for _, match := range conf.Redis.Match {
		log.Printf("Match: \033[1;36;1m%s\033[0m\n", match)
		instanceDelete(rdb, match)
	}
}

func deleteKeyCluster() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Username: conf.Cluster.Username,
		Password: conf.Cluster.Password,
		Addrs:    conf.Cluster.Addrs,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	//for i := 0; i < 330; i++ {
	//	err := rdb.Set(ctx, fmt.Sprintf("opstest:tmp%d", i), "value", 0).Err()
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	log.Printf("Redis-ClusterInstance :%s", conf.Cluster.Addrs)
	for _, match := range conf.Cluster.Match {
		log.Printf("Match:: \033[1;36;1m%s\033[0m\n", match)
		err := rdb.ForEachMaster(ctx, func(ctx context.Context, shard *redis.Client) error {
			instanceDelete(shard, match)
			return nil
		})
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
