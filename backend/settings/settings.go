package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify" //注意这个要是github那个fsnotify
	"github.com/spf13/viper"
)

var Conf AppConfig

// 结构体配置文件
type AppConfig struct {
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     int    `mapstructure:"password"`
	DB           string `mapstructure:"db"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

func Init() {
	viper.SetConfigFile("./conf/config.yaml")
	// 动态监视配置的改变
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("糟糕!配置文件被修改了!")
		viper.Unmarshal(&Conf)
	})

	// 把配置信息反序列化到AppConfig对象中
	viper.ReadInConfig()
	err := viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(Conf.MySQLConfig)
	fmt.Println(Conf.RedisConfig)
	fmt.Println(Conf.LogConfig)
}
