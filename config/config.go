package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

//
// 配置初始化
//
func Init() {
	// 当前运行路径
	/*dir, err := base.GetAppPathAbs() //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		panic(err)
	}*/
	currentDir, err := os.Getwd()
	//fmt.Println(dir)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	//viper.AddConfigPath(dir)      // path to look for the config file in
	viper.AddConfigPath(currentDir) // path to look for the config file in
	err = viper.ReadInConfig()      // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
