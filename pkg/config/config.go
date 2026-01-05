// Package config lets you retrieve config variables
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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

// GetDurationS returns the value of the key as a duration
func GetDurationS(key string) time.Duration {
	bindEnv(key)
	return viper.GetDuration(key) * time.Second
}

// GetDefaultDurationS returns the value of the in time.Duration or a default value
// The default value should be the amount of seconds and gets transformed to a time.Duration
func GetDefaultDurationS(key string, defaultVal int) time.Duration {
	viper.SetDefault(key, defaultVal)
	return GetDurationS(key)
}

func IsDev() bool {
	return strings.ToLower(GetDefaultString("app.env", "development")) == "development"
}
