package fsyn

import (
	"github.com/Unknwon/goconfig"
	"log"
	"strconv"
)

type Config struct {
	parseobj     *goconfig.ConfigFile
	redis_ip     string
	redis_port   int64
	mysql_ip     string
	mysql_port   int64
	mysql_user   string
	mysql_passwd string
	mysql_db     string
	log_dir      string
}

func NewConfig(filename string) *Config {
	var err error

	configobj := &Config{}
	configobj.parseobj, err = goconfig.LoadConfigFile(filename)

	if err != nil {
		log.Println("init configobj fail")
		return nil
	}

	return configobj
}

func (configobj Config) GetLogdir() string {
	return configobj.log_dir
}

func (configobj Config) GetRedisIp() string {
	return configobj.redis_ip
}

func (configobj Config) GetRedisPort() int64 {
	return configobj.redis_port
}

func (configobj *Config) GetMysqlIp() string {
	return configobj.mysql_ip
}

func (configobj *Config) GetMysqlPort() int64 {
	return configobj.mysql_port
}

func (configobj *Config) GetMysqlUser() string {
	return configobj.mysql_user
}

func (configobj *Config) GetMysqlPasswd() string {
	return configobj.mysql_passwd
}

func (configobj *Config) GetMysqlDb() string {
	return configobj.mysql_db
}

func (configobj *Config) InitLogSection() (err error) {

	if configobj.parseobj != nil {
		configobj.log_dir, err = configobj.parseobj.GetValue("log", "log_dir")
		if err != nil {
			log.Println("read log dir item from log section fail")
			return err
		}
	}

	return err
}

func (configobj *Config) InitRedisSection() (err error) {

	if configobj.parseobj != nil {

		configobj.redis_ip, err = configobj.parseobj.GetValue("redis", "ip")
		if err != nil {
			log.Println("read ip item from redis section fail")
			return err
		}

		var port string

		port, err = configobj.parseobj.GetValue("redis", "port")

		if err != nil {
			log.Println("read port item from redis section fail")
			return err
		}

		configobj.redis_port, err = strconv.ParseInt(port, 10, 0)
	}

	return err
}

func (configobj *Config) InitMysqlSection() (err error) {

	if configobj.parseobj != nil {
		configobj.mysql_ip, err = configobj.parseobj.GetValue("mysql", "ip")
		if err != nil {
			log.Println("read ip item from mysql section fail")
			return err
		}

		var port string
		port, err = configobj.parseobj.GetValue("mysql", "port")

		if err != nil {
			log.Println("read port item from mysql section fail")
			return err

		}

		configobj.mysql_port, _ = strconv.ParseInt(port, 10, 0)

		configobj.mysql_user, err = configobj.parseobj.GetValue("mysql", "user")
		if err != nil {
			log.Println("read user item from mysql section fail")
			return err
		}

		configobj.mysql_passwd, err = configobj.parseobj.GetValue("mysql", "passwd")

		if err != nil {
			log.Println("read passwd item from mysql section fail")
			return err
		}

		configobj.mysql_db, err = configobj.parseobj.GetValue("mysql", "db")
		if err != nil {
			log.Println("read db item from mysql section fail")
			return err
		}
	}

	return err
}
