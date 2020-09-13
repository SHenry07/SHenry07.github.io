[go-redis-gihub](https://github.com/go-redis/redis)

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
# cli 命令
```shell
get KEYNAME // 取值
keys *     // 模糊查找/全部key KEYS y* // 2.8 scan会更好
// 因为redis是单线程的, 数据集过大会造成 服务器堵塞
SCAN 0 MATCH l* // COUNT 默认是10
SMEMBERS key // 返回集合中的所有成员
TYPE key  // 类型
TTL key  // 过期时间
INFO replication //集群状态
```
# Redis Persistence 持久化

[官方文档](https://redis.io/topics/persistence)

```
$ redis-cli config get 'append*'
1) "appendfsync"
2) "always"
3) "appendonly"
4) "yes"
```

Redis 提供了两种数据持久化的方式，分别是快照和追加文件。

- **快照方式**，会按照指定的时间间隔，生成数据的快照，并且保存到磁盘文件中。为了避免阻塞主进程，Redis 还会 fork 出一个子进程，来负责快照的保存。这种方式的性能好，无论是备份还是恢复，都比追加文件好很多。不过，它的缺点也很明显。在数据量大时，fork 子进程需要用到比较大的内存，保存数据也很耗时。所以，你需要设置一个比较长的时间间隔来应对，比如至少 5 分钟。这样，如果发生故障，你丢失的就是几分钟的数据。
- **追加文件**，则是用在文件末尾追加记录的方式，对 Redis 写入的数据，依次进行持久化，所以它的持久化也更安全。此外，它还提供了一个用 appendfsync 选项设置 fsync 的策略，确保写入的数据都落到磁盘中，具体选项包括 always、everysec、no 等。
  - always 表示，每个操作都会执行一次 fsync，是最为安全的方式；
  - everysec 表示，每秒钟调用一次 fsync ，这样可以保证即使是最坏情况下，也只丢失 1 秒的数据；
  - 而 no 表示交给操作系统来处理。



