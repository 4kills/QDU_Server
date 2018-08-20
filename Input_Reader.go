package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"runtime"
	"strings"
)

// Haupteinstiegsfunktion für das Lesen von User-Input
func readInput() {
	log.Print("QDU-Server launched...\n",
		"Launch servers by entering the respetive settings\n\n")

	for {
		if err := manageSettings(); err != nil {
			log.Println(err)
			continue
		}
		break
	}

	go listenCommands()
}

func listenCommands() {
	log.Println("\nWrite '/help' for a list of commands")
	log.Print("-------------------------------------\n\n")
	for {
		s, err := readLine()
		if err != nil {
			log.Println(err)
		}
		checkCommands(s)
	}
}

func manageSettings() error {

	log.Print("Do you want to load a setting file \n",
		"or create and initialize a new one?\n",
		"Load [0] / Create [1]\n\n")
	s, err := readLine()
	if err != nil {
		return err
	}
	if s == "0" {
		// Load
		return loadSettingsFile()
	}
	if s == "1" {
		// Create
		return createSettingsFile()
	}
	return errors.New("input was not in a correct format")
}

func loadSettingsFile() error {
	log.Print("Please enter the absolute(!) file path to your config (qdn-file)\n")

	s, err := readLine()
	if err != nil {
		return err
	}

	if err = qdnFileToStruct(s, &config); err != nil {
		return err
	}

	log.Print("\nSucessfully loaded config\n")

	return nil
}

func createSettingsFile() error {
	log.Print("Please enter the path + (name).qdn to your config file\n\n")

	s, err := readLine()
	if err != nil {
		return err
	}
	path := s

	log.Print("Please enter your domain <[sub.]domain.toplvl>\n",
		"(this will be used for creating the URL that is sent to the connected client):\n\n")
	s, err = readLine()
	if err != nil {
		return err
	}
	config.Domain = s

	log.Print("Please enter your web directory\n",
		"(this will be used for creating the URL that is sent to the connected client):\n\n")
	s, err = readLine()
	if err != nil {
		return err
	}
	config.DirectoryWeb = s

	log.Print("Please enter the directory (absolute path) where the images should be saved:\n\n")
	s, err = readLine()
	if err != nil {
		return err
	}
	config.DirectoryPics = s

	log.Print("Please enter the desired TCP port:\n\n")
	s, err = readLine()
	if err != nil {
		return err
	}
	config.PortTCP = ":" + s

	log.Print("Please enter the desired web port on which this web server will run:\n\n")
	s, err = readLine()
	if err != nil {
		return err
	}
	config.PortWeb = ":" + s

	return structToQDNFile(path, config)
}

// Liest Konsolenzeile und lässt diese auf Befehle überprüfen
func readLine() (string, error) {
	// Liest Konsolen-Stream und bricht bei "Enter" ab zu lesen
	const inputDelimiter = '\n'
	r := bufio.NewReader(os.Stdin)
	input, err := r.ReadString(inputDelimiter)
	if err != nil {
		return "", err
	}

	// Formatiert input
	input = strings.Replace(input, "\n", "", -1)
	input = strings.Trim(input, " ")
	if runtime.GOOS == "windows" {
		input = input[:len(input)-1]
	}
	input = strings.Replace(input, " ", "", -1)

	// überprüft auf Befehle

	return input, nil
}

// Überprüft auf Befehle und setzt die Channel
func checkCommands(s string) {
	// TODO: reduce hard coded strings (via export?)
	switch {
	case s == "/getPortTCP":
		if len(config.PortTCP) < 1 {
			log.Print("\n\n")
			break
		}
		log.Print(config.PortTCP[1:], "\n\n")
	case s == "/getPortWeb":
		if len(config.PortWeb) < 1 {
			log.Print("\n\n")
			break
		}
		log.Print(config.PortWeb[1:], "\n\n")
	case s == "/getDirWeb":
		log.Print(config.DirectoryWeb)
	case s == "/help":
		log.Print("\n-------------------------------------------------------------------------------",
			"\n/help \t\t\t\t- returns a list of commands",
			"\n/getDomain \t\t\t- gets current domain for links",
			"\n/getDir \t\t\t- gets output directory of pics",
			"\n/getDirWeb \t\t\t- gets directory for web requests",
			"\n/getPortTCP \t\t\t- gets the admin defined port",
			"\n/getPortWeb \t\t\t- gets the web port",
			"\n/quit \t\t\t\t- exits the application w/ exit status 3",
			"\n/info \t\t\t\t- returns information about the programm",
			"\n\n- Be aware to port-forward to your server and to adjust firewall settings",
			"\n-------------------------------------------------------------------------------\n\n")
	case s == "/getDir":
		log.Println(config.DirectoryPics)

	case s == "/getDomain":
		log.Println(config.Domain)
	case s == "/quit":
		os.Exit(3)
	case s == "/info":
		log.Print("\n-------------------------------------------------------------------------------",
			"\nDevelopment Framework: Golang 1.10",
			"\nVersion: 1.1 (beta)",
			"\nCreator: Dominik Ochs",
			"\n-------------------------------------------------------------------------------\n\n")
	}
}
