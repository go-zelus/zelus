package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Load(file ...string) {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	configName, configPath, configType := resolve(file...)
	configPath = filepath.Join(workPath, configPath)
	configType = strings.ToLower(configType)
	if configType != "yaml" && configType != "yml" && configType != "json" && configType != "toml" && configType != "properties" {
		panic("unsupported file")
	}
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType(configType)
	err = viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
}

func resolve(args ...string) (fName, fPath, fType string) {
	fName = "config"
	fPath = "."
	fType = "yml"
	if len(args) == 0 {
		return
	}
	//兼容windows
	arg := strings.TrimSpace(args[0])
	arg = strings.ReplaceAll(arg, "\\\\", "/")
	index1 := strings.LastIndex(arg, "/")
	index2 := strings.LastIndex(arg, ".")
	if index1 != -1 {
		fPath = arg[0:index1]
	}
	if index2 != -1 {
		fName = arg[index1+1 : index2]
		fType = arg[index2+1:]
	}
	return
}

// Get 获取所有类型的配置
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString 获取字符串类型的配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetBool 获取布尔类型的配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetInt 获取整数类型的配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetStringSlice 获取字符串数组类型的配置
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetStringMap 获取map接口配置
func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// AllSettings 所有配置
func AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

// UnmarshalKey 根据key解析配置到指定的结构中
func UnmarshalKey(key string, config interface{}) error {
	err := viper.UnmarshalKey(key, config)
	if err != nil {
		fmt.Printf("failed to resolve configuration key:[%s], %v \n", key, err)
	}
	return err
}

// UnmarshalFile 解析整个配置到指定的结构中
func UnmarshalFile(configFile string, config interface{}) error {
	Load(configFile)
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("failed to resolve configuration: %v \n", err)
	}
	return err
}
