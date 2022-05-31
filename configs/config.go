package configs

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/ethan256/books/pkg/log"
)

var config = new(Config)

type Config struct {
	Host string `toml:"host"`

	MySQL struct {
		Addr            string        `toml:"addr"`
		User            string        `toml:"user"`
		Pass            string        `toml:"pass"`
		Name            string        `toml:"name"`
		MaxOpenConn     int           `toml:"maxOpenConn"`
		MaxIdleConn     int           `toml:"maxIdleConn"`
		ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
	} `toml:"mysql"`

	Redis struct {
		Addr         string `toml:"addr"`
		Pass         string `toml:"pass"`
		Db           int    `toml:"db"`
		MaxRetries   int    `toml:"maxRetries"`
		PoolSize     int    `toml:"poolSize"`
		MinIdleConns int    `toml:"minIdleConns"`
	} `toml:"redis"`
}

func InitConfig() error {

	viper.SetConfigFile("./configs/configs.toml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "unmarshal config error")
	}

	// 监听配置文件的变化，并通过`unmarshal`更新至全局的config中
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			log.Logger.Error().Err(err).Msg("update config error")
		}
	})

	fmt.Printf("%+v\n", config)
	return nil
}

func Get() Config {
	return *config
}
