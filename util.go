package main

import (
	"encoding/json"
	"io/ioutil"
)

func ReadConfig(filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		return err
	}

	return nil
}
