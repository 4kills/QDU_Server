package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/4kills/qdn"
)

type configuration struct {
	Domain        string
	DirectoryPics string
	DirectoryWeb  string
	PortTCP       string
	PortWeb       string
}

func structToQDNFile(path string, stru interface{}) error {
	raw, err := qdn.Marshal(stru)
	if err != nil {
		return err
	}

	raw, err = qdn.Format(raw)
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

func qdnFileToStruct(path string, stru interface{}) error {
	if path[len(path)-4:] != ".qdn" {
		return errors.New("No .qdn-file provided")
	}

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

	err = qdn.Unmarshal(stru, raw)
	if err != nil {
		return err
	}

	return nil
}
