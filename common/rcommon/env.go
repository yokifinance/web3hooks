package rcommon

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var parsedSecrets map[string]string

// Loads default env files (if any)
func LoadDefaultEnvFiles(inTests bool) {
	var envFilePath string
	if inTests {
		envFilePath = "../../test.env"
	} else {
		envFilePath = "../.env"
	}

	Println("loading default env from '%s'", envFilePath)
	err := godotenv.Load(envFilePath)
	if err != nil {
		Println("cannot load default env secrets from %s. Using values set by shell", envFilePath)
	}
}

// Loads secrets and !!!does not!! load them into ENV for this process (since they're secrets, not to be visible by others)
func LoadSecrets(filename string) {
	Println("Loading secrets %s", filename)

	file, err := os.Open(filename)
	if err != nil {
		Println("loading secrets: %s:", err)
		return
	}
	defer file.Close()

	if parsedSecrets == nil {
		parsedSecrets = map[string]string{}
	}

	nowSecrets, err := godotenv.Parse(file)
	if err != nil {
		Println("parsing secrets file: %s:", err)
		return
	}

	parsedSecrets = mergeMaps(parsedSecrets, nowSecrets)
}

// Get all params (from env and secrets)
func GetAllParams() map[string]string {
	res := map[string]string{}
	allEnvVars := os.Environ()

	for _, variable := range allEnvVars {
		pair := strings.SplitN(variable, "=", 2)
		res[pair[0]] = pair[1]
	}

	res = mergeMaps(res, parsedSecrets)
	return res
}

func SetParamStrOrLeaveDefault(value *string, env string) {
	env_val := getParam(env)
	if env_val == "" {
		return
	}
	*value = env_val
}

func GetParamStrOrDefault(env, def string) string {
	value := getParam(env)
	if value == "" {
		return def
	} else {
		return value
	}
}

func GetParamStrOrFail(env string) string {
	value := strings.ToLower(getParam(env))
	if value == "" {
		log.Fatalf("No '%s' specified in env config or secrets", env)
	}
	return value
}

func GetParamIntOrDefault(env string, def int) int {
	value := getParam(env)
	if value == "" {
		return def
	} else {
		res, err := strconv.Atoi(value)
		if err != nil {
			return def
		}
		return res
	}
}

func getParam(param string) string {
	value, ok := parsedSecrets[param]
	if !ok {
		return os.Getenv(param)
	}
	return value
}

/*
Merges two maps. If key is present in both maps, the value from highPriorityMap is used.
Returns a new map, leaving two original maps intact.
*/
func mergeMaps(lowPriorityMap, highPriorityMap map[string]string) map[string]string {
	res := make(map[string]string)
	for key, value := range lowPriorityMap {
		res[key] = value
	}
	for k, v := range highPriorityMap {
		res[k] = v
	}
	return res
}
