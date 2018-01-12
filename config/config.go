package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const rootConfigLocation = "./config.json"

// Config for the ms-geo-data web service
var Config *config

func init() {
	log.Println("calling ms-geo-data.config.init()")
	Config = newConfig()
}

type config struct {
	Version          string
	Environment      string
	Port             int
	LogglyKey        string
	LogLevel         string
	ConfigSearchPath []string
}

func newConfig() *config {

	found := false
	var c *config
	tmpc := &config{}

	if _, err := os.Stat(rootConfigLocation); !os.IsNotExist(err) {
		fmt.Println("Found base config file.")
		raw, err := ioutil.ReadFile(rootConfigLocation)
		if err != nil {
			log.Panic(fmt.Sprintf("Base config file could not be read: \n\t%v", err.Error()))
		}
		err = json.Unmarshal(raw, tmpc)
		if err != nil {
			log.Panic(fmt.Sprintf("Base config file could not be parsed: \n\t%v", err.Error()))
		}

		c = tmpc
	} else {
		log.Panic("No base config file found!")
	}

	// This is a list of possible locations for the config file. The highest-priority location (the one
	//   that "wins" all conflicts) comes first. The "built-in" (dev) config is assumed to be at "./config.json"
	if c == nil {
		log.Panic("Base config is empty!")
	}
	if c.ConfigSearchPath == nil || len(c.ConfigSearchPath) == 0 {
		foo, _ := json.MarshalIndent(c, "", " ")
		log.Panic("Base config contains no configSearchPath: " + string(foo) + "\n")
	}
	configSearchPath := c.ConfigSearchPath

	for i := len(configSearchPath) - 1; i >= 0; i-- {
		thisConfig := configSearchPath[i]
		if _, err := os.Stat(thisConfig); !os.IsNotExist(err) {
			fmt.Printf("Found config file at %v \n", thisConfig)
			raw, err := ioutil.ReadFile(thisConfig)
			if err != nil {
				log.Panic(fmt.Sprintf("Config file %v could not be read: \n\t%v", thisConfig, err.Error()))
			}
			err = json.Unmarshal(raw, tmpc)
			if err != nil {
				log.Panic(fmt.Sprintf("Config file %v could not be parsed: \n\t%v", thisConfig, err.Error()))
			}

			if c == nil {
				// this is the first config we've found
				c = tmpc
			} else {
				c.overwriteFields(tmpc)
			}
			c.validate()
			found = true
		} else {
			fmt.Printf("No config file found at %v \n", thisConfig)
		}
	}
	if !found {
		fmt.Printf("Config file not found in any of: \n\t%v", strings.Join(configSearchPath, "\n\t"))
	}

	return c
}

func (c *config) validate() {
	// PANICS
	if c.Version == "" {
		log.Panic("ms-geo-data config error: Version cannot be blank")
	}
	if c.Environment == "" {
		log.Panic("ms-geo-data config error: Environment cannot be blank")
	}
	if c.Port == 0 {
		log.Panic("ms-geo-data config error: Port cannot be blank")
	}

	// WARNINGS
	if c.LogglyKey == "" {
		log.Warn("\n\n-- WARNING -- (config) Loggly key is blank\n\n ")
	}

	// verify & set default log level
	lev := strings.ToLower(c.LogLevel)

	if lev != "debug" && lev != "info" && lev != "warning" && lev != "error" && lev != "fatal" && lev != "panic" {
		log.Panic("ms-geo-data config error: LogLevel must be debug, info, warning, error, fatal, or panic.")
	} else {
		c.LogLevel = lev
	}
}

func (c *config) overwriteFields(src *config) {
	if src == nil {
		log.Panic("ms-geo-data config error: overwriteFields: source config must not be nil")
	}
	if c == nil {
		log.Panic("ms-geo-data config error: overwriteFields: destination config must not be nil")
	}

	if src.Version != "" && src.Version != c.Version {
		fmt.Printf("Overwriting Version: %v => %v\n", c.Version, src.Version)
		c.Version = src.Version
	}
	if src.Environment != "" && src.Environment != c.Environment {
		fmt.Printf("Overwriting Environment: %v => %v\n", c.Environment, src.Environment)
		c.Environment = src.Environment
	}
	if src.Port != 0 && src.Port != c.Port {
		fmt.Printf("Overwriting Port: %v => %v\n", c.Port, src.Port)
		c.Port = src.Port
	}
	if src.LogglyKey != "" && src.LogglyKey != c.LogglyKey {
		fmt.Printf("Overwriting LogglyKey: %v => %v\n", c.LogglyKey, src.LogglyKey)
		c.LogglyKey = src.LogglyKey
	}
	if src.LogLevel != "" && src.LogLevel != c.LogLevel {
		fmt.Printf("Overwriting LogLevel: %v => %v\n", c.LogLevel, src.LogLevel)
		c.LogLevel = src.LogLevel
	}
}
