package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/viper"
)

type Mysql struct {
	Addr     string
	Username string
	Password string
	DBName   string
}

func (mysql Mysql) DataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysql.Username, mysql.Password, mysql.Addr, mysql.DBName)
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Minio struct {
	Host              string
	Port              int
	Region            string
	Secure            bool
	AccessKey         string
	SecretAccessKey   string
	BucketName        string
	BasePath          string
	Multipart         bool
	AdvertisedHost    string
	ConnectionTimeout time.Duration
}

type Config struct {
	Mysql Mysql `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
	Minio Minio `yaml:"minio"`
}

func LoadConf(filepath string) (error, *Config) {
	_, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err, nil
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath)

	err = viper.ReadInConfig()
	if err != nil {
		return err, nil
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return err, nil
	}
	return nil, conf
}
