```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

var redisdb *redis.Client

// NewClient ...
func NewClient() {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.100.135:6379",
		Password: "mysteel123",
		DB:       0,
	})
	pong, err := redisdb.Ping().Result()

	if err != nil {
		log.Println(pong, err)
		return
	}
	// log.Println(pong, err)
}

func set(key string) {
	err := redisdb.Set(key, 24, time.Second*50).Err()
	if err != nil {
		defer log.Println("set key failed")
		panic(err)
	}
	ret, err := redisdb.Get(key).Result()
	if err != nil {
		panic(err)
	}
	log.Printf("%s'value is %s", key, ret)

	ret, err = redisdb.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", ret)
	}
}

func sortedset(zsetKey string, set []*redis.Z) {
	num, err := redisdb.ZAdd(zsetKey, set...).Result()
	if err != nil {
		log.Println("zadd failed, err: %v\n", err)
		return
	}
	log.Printf("zadd %d success \n", num)

}

func main() {
	NewClient()
	// set("sunheng")
	zsetKey := "language_rank"
	languages := []*redis.Z{
		&redis.Z{Score: 90.0, Member: "Golang"},
		&redis.Z{Score: 98.0, Member: "Java"},
		&redis.Z{Score: 95.0, Member: "Python"},
		&redis.Z{Score: 97.0, Member: "JavaScript"},
		&redis.Z{Score: 99.0, Member: "C"},
		&redis.Z{Score: 98.0, Member: "C++"},
	}
	sortedset(zsetKey, languages)

	newScore, err := redisdb.ZIncrBy(zsetKey, 5.0, "Golang").Result()
	if err != nil {
		log.Println("zincrby failed, err: %v\n", err)
		return
	}
	log.Printf("Golang's score is %f now. \n", newScore)

	// 前三名
	ret, err := redisdb.ZRevRangeWithScores(zsetKey, 0, 2).Result()

	for _, value := range ret {
		log.Println(value.Member, value.Score)
	}

	log.Println("区间分数")
	// 区间分数
	op := &redis.ZRangeBy{
		Min:    "96",
		Max:    "100",
		Offset: 0,
		Count:  2,
	}

	vals, err := redisdb.ZRangeByScoreWithScores(zsetKey, op).Result()
	if err != nil {
		log.Println("zrangeByScore failed, err:", err)
		return
	}

	for _, z := range vals {
		log.Println(z.Member, z.Score)
	}

}

```

# sqlx

bilibi在用