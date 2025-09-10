package config

import (
	"fmt"
	// "path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Port int `mapstructure:"port"`
}

func LoadConfig(path string) (config *Config, err error) {
	// Конфигурируем Viper
	viper.SetConfigName("config") // имя файла без расширения
	viper.SetConfigType("yaml")   // тип конфигурационного файла
	viper.AddConfigPath(path)     // путь к директории с конфигом

	// Устанавливаем значения по умолчанию
	viper.SetDefault("port", 8282)

	// Пытаемся прочитать конфигурационный файл
	err = viper.ReadInConfig()
	if err != nil {
		// Если файл не найден - это не фатальная ошибка,
		// т.к. у нас есть значения по умолчанию
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Конфиг файл не найден, используем значения по умолчанию
			config = &Config{}
			viper.Unmarshal(config)
			return config, nil
		} else {
			// Другая ошибка (например, файл есть, но он некорректный)
			return nil, fmt.Errorf("fatal error config file: %w", err)
		}
	}

	// Если файл найден и прочитан успешно
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}

// Альтернативная версия функции с явным указанием пути к файлу
func LoadConfigExplicit(configPath string) (config *Config, err error) {
	viper.SetConfigFile(configPath) // Полный путь к файлу конфигурации
	viper.SetDefault("port", 8282)

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			config = &Config{}
			viper.Unmarshal(config)
			return config, nil
		}
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return config, nil
}
