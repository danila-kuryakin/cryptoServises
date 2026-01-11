package config

import (
	"bufio"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config — корневая конфигурация приложения
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}
type DatabaseConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Name    string `yaml:"dbname"`
	SSLMode string `yaml:"sslmode"`
}

type KafkaConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

// LoadConfig загружает конфигурацию из YAML файла
func LoadConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Println(err)
	}

	return &cfg
}

// LoadEnv загружает .env файл в окружение. Строки .env файла:
// DB_USER - Имя пользователя
// DB_PASSWORD - Пароль
// API_KEY - api ключь телеграм бота
func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing file", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// пропускаем комментарии и пустые строки
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		err := os.Setenv(key, value)
		if err != nil {
			log.Println("Error setting env var", err)
			return
		}
	}
}
