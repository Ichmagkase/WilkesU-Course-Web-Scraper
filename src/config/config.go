/*
 * file: config.go
 * Author: Zackery Drake
 * Last updated: 5/26/2025
 * Description:
 *   Loads config.json into a struct for use
 *   throughout the application. The path to
 *   config.json will be found in the environment variable:
 *
 *   WILKESUSCRAPY_CONFIG
 *
 *   If a configuration path is not found, it will attempt
 *   to create a new config.json with default values in
 *   createDefault.
 *
 *   NOTE: As new entries are added to config.json,
 *   entries must also be added to the Configuration
 *   struct, and default values should be updated under
 *   defaultConfig in createDefault.
*/
package config

import (
	"fmt"
	"encoding/json"
	"os"
	"log"
)

type Configuration struct {
	ConfigPath   string   `json:"ConfigPath"`
	MongoUri     string   `json:"MongoUri"`
}

/*
 * Create a default config if one does
 * not exist.
 * Arguments:
 *   path : Expected path to config.json file
 */
func createDefault(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("config.go: ", err)
	}

	defaultConfig, err := json.Marshal(Configuration{
		ConfigPath:  "config/config.json",
		MongoUri:    "mongodb://mongodb:27017",
	})
	if err != nil {
		log.Fatal("config.go: ", err)
	}

	fmt.Fprintf(file, string(defaultConfig))
}

/*
 * Load config.json into Configuration struct
 */
func LoadConfig() (Configuration) {
	path := os.Getenv("WILKESUSCRAPY_CONFIG")
	if path == "" {
		path = "config/config.json"
	}
	file, err := os.Open(path)
	if err != nil {
		createDefault(path)
	}

	file, err = os.Open(path)
	if err != nil {
		log.Fatal("config.go: ", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("config.go: ", err)
	}
	return configuration
}

// For testing purposes only
func main() {
	fmt.Println("Loading config")
}
