package main

import (
	"github.com/go-redis/redis"
	"fmt"
)

func NewRClient() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:    "10.50.115.17:7003",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}

func NewRClusterClient() *redis.ClusterClient{
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"10.50.115.17:7006", "10.50.115.17:7001", "10.50.115.17:7002", "10.50.115.17:7003", "10.50.115.17:7004", "10.50.115.17:7005"},
	})
	pong, err := client.Ping().Result();
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, err)
	return client
}

func SimpleRedisClient() {
	client := NewRClient()
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key20170405").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		//panic(err)
		fmt.Println("error info :", err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exists
}

func ClusterRedisClient(){
	client := NewRClusterClient()

	err := client.Set("123456", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := client.Get("123456").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("123456", val)
}

//go get -u github.com/go-redis/redis
func main(){
	//SimpleRedisClient()				//single node
	ClusterRedisClient()				//cluster nodes
}
