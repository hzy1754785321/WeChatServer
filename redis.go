package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var pool = &redis.Pool{
	MaxIdle:   3, /*最大的空闲连接数*/
	MaxActive: 8, /*最大的激活连接数*/
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "localhost:6379")
		if err != nil {
			fmt.Println("redis数据库连接出错,", err)
			return nil, err
		}
		return c, nil
	},
}

func SetRedis(key string, value string) {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
}

func GetRedis(key string) (value string) {
	conn := pool.Get()
	defer conn.Close()
	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	return data
}

func CheckRedis(key string) (exists bool) {
	conn := pool.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		fmt.Println("illegal exception:", err)
	}
	return exists
}

func SetKeyExpire(key string, time int) {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("expire", key, time)
	if err != nil {
		fmt.Println("redis SetExpire failed:", err)
	}
}

func DeleteRedis(key string) {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		fmt.Println("redis delelte value failed >>>", err)
	}
}

func GetRedisListValue(key string, index int) (value int) {
	conn := pool.Get()
	defer conn.Close()
	data, err := redis.Int(conn.Do("lindex", key, index))
	if err != nil {
		fmt.Println("redis get list value failed:", err)
	}
	return data
}

func CheckRedisSet(key string, value int) (exists int) {
	conn := pool.Get()
	defer conn.Close()
	exists, err := redis.Int(conn.Do("SISMEMBER", key, value))
	if err != nil {
		fmt.Println("Value already exists :", err)
	}
	return exists
}

func AddRedisSet(key string, value string) (result int) {
	conn := pool.Get()
	defer conn.Close()
	exists, err := redis.Int(conn.Do("SADD", key, value))
	if err != nil {
		fmt.Println("RedisSet add value failed :", err)
	}
	return exists
}

func GetRedisSet(key string) (result []string) {
	conn := pool.Get()
	defer conn.Close()
	result, err := redis.Strings(conn.Do("smembers", key))
	if err != nil {
		fmt.Println("RedisSet get value failed :", err)
	}
	return result
}
