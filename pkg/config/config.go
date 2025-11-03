// Package config lets you retrieve config variables
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func bindEnv(key string) {
	envName := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	// nolint:errcheck // we do not care if it binds
	viper.BindEnv(key, envName)
}

func Init() error {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println("Failed to load .env file", err)
		}
	}

	viper.AutomaticEnv()
	env := GetDefaultString("app.env", "development")

	// Config variables should be set in Infisical
	// but you can overwrite them with a toml file if needed

	viper.SetConfigName(env + ".toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")

	return viper.ReadInConfig()
}

func GetString(key string) string {
	bindEnv(key)
	return viper.GetString(key)
}

func GetDefaultString(key, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return GetString(key)
}

func GetInt(key string) int {
	bindEnv(key)
	return viper.GetInt(key)
}

func GetDefaultInt(key string, defaultVal int) int {
	viper.SetDefault(key, defaultVal)
	return GetInt(key)
}

func GetUint16(key string) uint16 {
	bindEnv(key)
	return viper.GetUint16(key)
}

func GetDefaultUint16(key string, defaultVal uint16) uint16 {
	viper.SetDefault(key, defaultVal)
	return GetUint16(key)
}

func GetBool(key string) bool {
	bindEnv(key)
	return viper.GetBool(key)
}

func GetDefaultBool(key string, defaultVal bool) bool {
	viper.SetDefault(key, defaultVal)
	return GetBool(key)
}

func IsDev() bool {
	return strings.ToLower(GetDefaultString("app.env", "development")) == "development"
}
