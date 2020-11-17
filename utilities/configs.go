package utilities

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfigure() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")    // 設定檔名稱（無後綴）
	v.SetConfigType("yaml")      // 設定字尾名 {"1.6以後的版本可以不設定該字尾"}
	v.AddConfigPath("./inputs/") // 設定檔案所在路徑
	v.Set("verbose", true)       // 設定預設引數

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(" Config file not found; ignore error if desired")
		} else {
			panic("Config file was found but another error was produced")
		}
	}
	// 監控配置和重新獲取配置
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	return v
}
