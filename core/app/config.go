package app

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

const (
	CONFIG_FILENAME = "config.yml"
)

var (
	cachedConfig  *Config
	defaultConfig = Config{
		Port:           8080,
		ReadTimeout:    60,
		WriteTimeout:   60,
		MaxHeaderBytes: 4096,
	}
)

type Config struct {
	Port           int     `yaml:"Port"`
	ReadTimeout    int     `yaml:"ReadTimeout"`
	WriteTimeout   int     `yaml:"WriteTimeout"`
	MaxHeaderBytes int     `yaml:"MaxHeaderBytes"`
	Name           *string `yaml:"Name"`
}

// Reads parameters from `config.yml` and from env vars. The first time this
// function is called, the Config struct is cached and subsequent calls will
// returned the cached copy. There are 3 sources of config info and they are
// prioritized as follows:
//
// 1. (high priority) config from env vars
// 2. (medium) config from `config.yml`
// 3. (low) default config, hardcoded into this file
func GetConfig() (Config, error) {
	if cachedConfig == nil {
		cfg, err := getConfig()
		if err != nil {
			return defaultConfig, err
		}
		cachedConfig = &cfg
	}
	return *cachedConfig, nil
}

func ReadConfig(filename string) (Config, error) {
	var cfg map[string]interface{}
	chugged, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(chugged, &cfg)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Port:           extractIntOrMinusOne(cfg, "Port"),
		ReadTimeout:    extractIntOrMinusOne(cfg, "ReadTimeout"),
		WriteTimeout:   extractIntOrMinusOne(cfg, "WriteTimeout"),
		MaxHeaderBytes: extractIntOrMinusOne(cfg, "MaxHeaderBytes"),
		Name:           extractString(cfg, "Name"),
	}, nil
}

func ReadEnvVarsConfig() Config {
	return Config{
		Port:           getEnvInt("PORT", -1),
		ReadTimeout:    getEnvInt("READ_TIMEOUT", -1),
		WriteTimeout:   getEnvInt("WRITE_TIMEOUT", -1),
		MaxHeaderBytes: getEnvInt("MAX_HEADER_BYTES", -1),
		Name:           getEnvString("NAME", nil),
	}
}

func MergeConfigs(cfg1 Config, cfg2 Config) Config {
	return Config{
		Port:           takeNonNegative(cfg1.Port, cfg2.Port),
		ReadTimeout:    takeNonNegative(cfg1.ReadTimeout, cfg2.ReadTimeout),
		WriteTimeout:   takeNonNegative(cfg1.WriteTimeout, cfg2.WriteTimeout),
		MaxHeaderBytes: takeNonNegative(cfg1.MaxHeaderBytes, cfg2.MaxHeaderBytes),
		Name:           takeNonNil(cfg1.Name, cfg2.Name),
	}
}

func CheckConfig(cfg Config) error {
	return nil
}

func takeNonNil(str1, str2 *string) *string {
	if str2 != nil {
		return str2
	}
	return str1
}

func takeNonNegative(int1, int2 int) int {
	if int2 < 0 {
		return int1
	}
	return int2
}

func getConfig() (Config, error) {
	configFromEnvVars := ReadEnvVarsConfig()
	configFromFile, err := ReadConfig(CONFIG_FILENAME)
	if err != nil {
		log.Printf("Error reading config file: %v\n", err)
		log.Println("Ignoring config file...")
		configFromFile = defaultConfig
	}

	finalConfig := MergeConfigs(
		MergeConfigs(defaultConfig, configFromFile),
		configFromEnvVars,
	)

	if err = CheckConfig(finalConfig); err != nil {
		return Config{}, err
	}

	return finalConfig, nil
}

func getEnvInt(param string, defaultValue int) int {
	p, err := strconv.Atoi(os.Getenv(param))
	if err != nil {
		return defaultValue
	}
	return p
}

func getEnvString(param string, defaultValue *string) *string {
	p := defaultValue
	if got, ok := os.LookupEnv(param); ok {
		p = &got
	}
	return p
}

func extractIntOrMinusOne(configMap map[string]interface{}, param string) int {
	got, ok := configMap[param]
	if !ok {
		return -1
	}
	converted, ok := got.(int)
	if !ok {
		return -1
	}
	return converted
}

func extractString(configMap map[string]interface{}, param string) *string {
	var p *string
	if got, ok := configMap[param]; ok {
		if cast, ok := got.(string); ok {
			p = &cast
		}
	}
	return p
}
