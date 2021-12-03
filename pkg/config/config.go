package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
)

// LoadConfigFromFile - loads the config from the given file path, but it will prioritize the env variable, so you can
// override the config file with the env variables. This is for helpfulness. In c interface{} you should insert an empty
// struct. c argument will most likely be replaced by generics in the future (when go generics releases).
func LoadConfigFromFile(filePath string, c interface{}) (interface{}, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the configuration file")
	}

	typeOfStruct := reflect.TypeOf(c)
	v := reflect.New(typeOfStruct)
	newP := v.Interface()

	err = yaml.Unmarshal(yamlFile, newP)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode the configuration file")
	}

	v = reflect.ValueOf(newP).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsValid() {
			envKey := typeOfStruct.Field(i).Tag.Get("env")
			value := v.Field(i).Interface()
			envValue, err:= setEnvVariable(envKey, value)
			if err != nil {
				return nil, fmt.Errorf("couldn't set the environment variable %s", envKey)
			}
			if v.Field(i).CanSet() {
				v.Field(i).SetString(envValue)
			}
		}
	}


	return v.Interface(), nil
}

func setEnvVariable(key string, defaultVal interface{}) (string, error) {
	if os.Getenv(key) == "" {
		err := os.Setenv(key, fmt.Sprintf("%v", defaultVal))
		if err != nil {
			return "nil", fmt.Errorf("couldn't set the env variable %s", key)
		}
	}

	return os.Getenv(key), nil
}
