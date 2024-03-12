package config

import (
	"fmt"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"os"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Server   *Server             `yaml:"service"`
	MySQL    *MySQL              `yaml:"mysql"`
	Redis    *Redis              `yaml:"redis"`
	Etcd     *Etcd               `yaml:"etcd"`
	Services map[string]*Service `yaml:"services"`
	Domain   map[string]*Domain  `yaml:"domain"`
	SeConfig *SeConfig           `yaml:"SeConfig"`
	Kafka    *Kafka              `yaml:"kafka"`
}

type StarRocks struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	LoadUrl  string `yaml:"load_url"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Charset  string `yaml:"charset"`
}

type SeConfig struct {
	StoragePath      string   `yaml:"StoragePath"`
	SourceFiles      []string `yaml:"SourceFiles"`
	MergeChannelSize int64    `yaml:"MergeChannelSize"`
	Version          string   `yaml:"Version"`
	SourceWuKoFile   string   `yaml:"SourceWuKoFile"`
	MetaPath         string   `yaml:"MetaPath"`
}

type Server struct {
	Port      string `yaml:"port"`
	Version   string `yaml:"version"`
	JwtSecret string `yaml:"jwtSecret"`
}

type MySQL struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
	Charset    string `yaml:"charset"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisUsername string `yaml:"redisUsername"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

type Service struct {
	Name        string   `yaml:"name"`
	LoadBalance bool     `yaml:"loadBalance"`
	Addr        []string `yaml:"addr"`
}

type Kafka struct {
	Address []string `yaml:"address"`
}

type Domain struct {
	Name string `yaml:"name"`
}

func InitConfig() {
	zap.S().Info(logs.RunFuncName())
	workDir, _ := os.Getwd()
	fmt.Println(workDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}
