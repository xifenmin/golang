//auth:fengming.xi
//Golang firsts example
//common:mysql dcsp sync to redis tools

package main

import (
	"./fsyn"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/fatih/set.v0"
	"log"
	"os"
	"strings"
	"time"
)

const (
	LBS_REDIS_FSS_KEY = "lbs.fss.client"
	LBS_MYSQL_TABLE   = "BGTP_ACCESSCODE_BASEINFO"
)

type FssObj struct {
	configobj *fsyn.Config
	mysqlobj  *fsyn.Mysql
	redisobj  *fsyn.Redis
}

func md5Value(src string) string {
	md5obj := md5.New()
	md5obj.Write([]byte(src))
	return hex.EncodeToString(md5obj.Sum(nil))
}

func redisMysqlComparison(record_map map[string]string, redisclient *fsyn.Redis) {

	for k, v := range record_map {
		redis_vale := redisclient.Hget(LBS_REDIS_FSS_KEY, k)
		if strings.EqualFold(md5Value(v), md5Value(redis_vale)) {
			log.Printf("ac:%s not update", k)

		} else {
			redisclient.Hset(LBS_REDIS_FSS_KEY, k, v)
			log.Printf("ac:%s,add:%s", k, v)
		}
	}
}

func redisMysqlDifference(record_map map[string]string, record_array []string) []interface{} {

	mysql := set.New()
	redis := set.New()

	for ac, _ := range record_map {
		mysql.Add(ac)
	}

	for i := 0; i < len(record_array); i++ {
		redis.Add(record_array[i])
	}

	redis_no_mysql := set.Difference(redis, mysql) //求在redis,不在mysql中集合差集

	return redis_no_mysql.List()
}

func main() {

	fss_obj := &FssObj{}

	fss_obj.configobj = fsyn.NewConfig("fss_syn.conf")

	err := fss_obj.configobj.InitLogSection()

	if err != nil {
		log.Println("init log section fail")
		return
	}

	logFile, logErr := os.OpenFile(fss_obj.configobj.GetLogdir()+"/fss_syn.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "syn_tool start Failed")
		os.Exit(1)
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err = fss_obj.configobj.InitRedisSection()
	if err != nil {
		log.Println("init redis section fail")
		return
	}

	err = fss_obj.configobj.InitMysqlSection()

	if err != nil {
		log.Println("init mysql section fail")
		return
	}

	fss_obj.mysqlobj = fsyn.NewMysql(fss_obj.configobj.GetMysqlIp(), int(fss_obj.configobj.GetMysqlPort()),
		fss_obj.configobj.GetMysqlUser(), fss_obj.configobj.GetMysqlPasswd(),
		fss_obj.configobj.GetMysqlDb())

	record_map := fss_obj.mysqlobj.Query("SELECT * from " + LBS_MYSQL_TABLE)

	fss_obj.redisobj, err = fsyn.NewRedis(fss_obj.configobj.GetRedisIp(), int(fss_obj.configobj.GetRedisPort()))

	if err != nil {
		log.Println("redis connect fail")
		return
	}

	for {
		redisMysqlComparison(record_map, fss_obj.redisobj)

		items := redisMysqlDifference(record_map, fss_obj.redisobj.HgetAll(LBS_REDIS_FSS_KEY))

		for _, item := range items {
			result := fss_obj.redisobj.Hdel(LBS_REDIS_FSS_KEY, item.(string))
			log.Printf("hdel key:%s,result:%d", item, result)
		}

		time.Sleep(3 * time.Second)
	}
}
