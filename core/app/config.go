package app

import (
	"log"
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"
)

const CONFIG_FILENAME = "config.yml"

var cache = struct {
	sync.Mutex
	config *Config
}{}

var defaultConfig = Config{
	Port:           8080,
	ReadTimeout:    60,
	WriteTimeout:   60,
	MaxHeaderBytes: 4096,
}

type Config struct {
	Port           int       `json:"Port"`
	ReadTimeout    int       `json:"ReadTimeout"`
	WriteTimeout   int       `json:"WriteTimeout"`
	MaxHeaderBytes int       `json:"MaxHeaderBytes"`
	Name           *string   `json:"Name"`
	CORS           *[]string `json:"CORS"`
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
	cache.Lock()
	defer cache.Unlock()
	if cache.config == nil {
		cfg, err := getConfig()
		if err != nil {
			return defaultConfig, err
		}
		cache.config = &cfg
	}
	return *cache.config, nil
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
		CORS:           extractSlice(cfg, "CORS"),
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
		Name:           takeNonNilStr(cfg1.Name, cfg2.Name),
		CORS:           takeNonNilSlice(cfg1.CORS, cfg2.CORS),
	}
}

func takeNonNilSlice(s1, s2 *[]string) *[]string {
	if s2 != nil {
		return s2
	}
	return s1
}

func takeNonNilStr(s1, s2 *string) *string {
	if s2 != nil {
		return s2
	}
	return s1
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

	if finalConfig.CORS == nil {
		finalConfig.CORS = &[]string{}
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

func extractSlice(configMap map[string]interface{}, param string) *[]string {
	s := []string{}

	got, ok := configMap[param]
	if !ok {
		return nil
	}

	gotSlice, ok := got.([]interface{})
	if !ok {
		return nil
	}

	for _, castMe := range gotSlice {
		cast, ok := castMe.(string)
		if !ok {
			return nil
		}
		s = append(s, cast)
	}

	return &s
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
