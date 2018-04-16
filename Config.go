package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type configuration struct {
	Domain        string `json:"domain"`
	DirectoryPics string `json:"directoryPics"`
	DirectoryWeb  string `json:"directoryWeb"`
	PortTCP       string `json:"portTCP"`
	PortWeb       string `json:"portWeb"`
}

func structToJSONFile(path string, stru interface{}) error {
	raw, err := json.Marshal(stru)
	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(raw)
	if err != nil {
		return err
	}
	return nil
}

func jsonFileToStruct(path string, stru interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var raw []byte
	raw, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, stru)
	if err != nil {
		return err
	}
	return nil
}
