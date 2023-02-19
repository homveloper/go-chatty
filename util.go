package main

import (
	"encoding/json"
	"io/ioutil"
)

func ReadConfig(filename string) (interface{}, error) {
	var conf interface{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(file), &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
