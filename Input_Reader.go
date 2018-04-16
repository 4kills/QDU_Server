package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

//---------------------------------------------------------
// Lesen und Auswerten von User-Input startet hier
//---------------------------------------------------------

// Haupteinstiegsfunktion für das Lesen von User-Input
func readInput() {
	fmt.Print("QDU-Server launched...\n",
		"Launch servers by entering the respetive settings\n\n")

	for {
		if err := manageSettings(); err != nil {
			fmt.Println(printTS, err)
			continue
		}
		break
	}

	fmt.Println("Write '/help' for a list of commands")
	fmt.Print("-------------------------------------\n\n")
	// liest konstant Konsolenzeile
	for {
		readLine()
	}
}

func manageSettings() error {

	fmt.Print("Do you want to load a setting file \n",
		"or create and initialize a new one?\n",
		"Load [0] / Create [1]\n")
	s, err := readLine()
	if err != nil {
		return err
	}
	if s == "0" {
		// Load
		return LoadSettingsFile()
	}
	if s == "1" {
		// Create
		return CreateSettingsFile()
	}
	return errors.New("input was not in a correct format")
}

func LoadSettingsFile() error {

}

func CreateSettingsFile() error {

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

	// switch in Go kann benutzt werden um lange if-else
	// ausdrücke leserlicher zu gestalten.
	switch {
	case len(s) > 11 && (s[:11] == "/setPortTCP"):
		c <- ":" + s[11:]
		fmt.Print("successfully set port\n\n")
	case s == "/getPortTCP":
		if len(port) < 1 {
			fmt.Print("\n\n")
			break
		}
		fmt.Print(port[1:], "\n\n")
	case len(s) > 11 && (s[:11] == "/setPortWeb"):
		tcpC <- ":" + s[11:]
		fmt.Print("successfully set portWeb\n\n")
	case s == "/getPortWeb":
		if len(portWeb) < 1 {
			fmt.Print("\n\n")
			break
		}
		fmt.Print(portWeb[1:], "\n\n")
	case s == "/getDirWeb":
		fmt.Print(directoryWeb)
	case len(s) > 10 && (s[:10] == "/setDirWeb"):
		dirC <- s[10:]
		fmt.Print("successfully set webDir\n\n")
	case s == "/help":
		fmt.Print("\n-------------------------------------------------------------------------------",
			"\n/help \t\t\t\t- returns a list of commands",
			"\n/getDomain \t\t\t- gets current domain for links",
			"\n/setDomain <domain>\t\t- sets domain for returning link",
			"\n/getDir \t\t\t- gets output directory of pics",
			"\n/setDir <abs path> \t\t- sets output directory of pics",
			"\n/getDirWeb \t\t\t- gets directory for web requests",
			"\n/setDirWeb <rel path>\t\t- sets directory for web requests; may be empty",
			"\n/getPortTCP \t\t\t- gets the admin defined port",
			"\n/setPortTCP <port> \t\t- sets the tcp port. Set this first",
			"\n/getPortWeb \t\t\t- gets the web port",
			"\n/setPortWeb <port>\t\t- sets the web port. Set this second",
			"\n/quit \t\t\t\t- exits the application w/ exit status 3",
			"\n/info \t\t\t\t- returns information about the programm",
			"\n\n- Be aware to port-forward to your server and to adjust firewall settings",
			"\n-------------------------------------------------------------------------------\n\n")
	case s == "/getDir":
		fmt.Println(directory)
	case len(s) > 7 && (s[:7] == "/setDir"):
		directory = s[7:]
		fmt.Print("successfully set output directory\n\n")
	case s == "/getDomain":
		fmt.Println(domain)
	case len(s) > 10 && (s[:10] == "/setDomain"):
		domain = s[10:]
		fmt.Print("successfully set domain\n\n")
	case s == "/quit":
		os.Exit(3)
	case s == "/info":
		fmt.Print("\n-------------------------------------------------------------------------------",
			"\nDevelopment Framework: Golang 1.10",
			"\nVersion: 1.0 (stable release)",
			"\nCreator: Dominik Ochs",
			"\n-------------------------------------------------------------------------------\n\n")
	}
}
