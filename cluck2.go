package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// These are constants
// ConnIPEn - the IP address for Enis' server. This is the main server.
// ConnIPRob - the IP address for Robert's server
// ConnPort - we are using port 8888
// ConnType - this client uses tcp. can be udp as well.
const (
	ConnIPEn  = "130.85.70.132" // Enis' IP address
	ConnIPRob = "130.85.70.193" // Robert's IP address
	ConnPort  = "8888"
	ConnType  = "tcp"
)

func main() {

	// Checks where the user is trying to connect to.
	// Correct := false
	for {
		// Asks user for which server they want to connect to.
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Which server would you like to log in to?\n Type 'en' for Enis or 'ro' for Robert..")
		hostOption, err := reader.ReadString('\n')
		hostOption = strings.Replace(hostOption, "\n", "", -1) // Makes this work on unix machines
		//hostOption = strings.Replace(hostOption, "\r\n", "", -1) // Makes this work on windows machines

		if err != nil {
			log.Fatal(err)
		}

		// This if else statement checks the user's input for which server they are connecting to.
		// If they give an invalid input, they will have to enter it again.
		if hostOption == "en" {
			fmt.Println("You chose Enis' server. Connecting now...")
			tcpAddr, err := net.ResolveTCPAddr(ConnType, ConnIPEn+":"+ConnPort)

			// Connect to server with tcp
			conn, err := net.DialTCP(ConnType, nil, tcpAddr)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			go printOutput(conn)

			for {
				writeInput(conn)
			}

		} else if hostOption == "ro" {
			fmt.Println("You chose Robert's server. Connecting now...")
			tcpAddr, err := net.ResolveTCPAddr(ConnType, ConnIPRob+":"+ConnPort)

			// Connect to server with tcp
			conn, err := net.DialTCP(ConnType, nil, tcpAddr)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			// First time entry

			for {
				writeInput(conn)
				go printOutput(conn)
			}

		} else {
			fmt.Println("You didn't pick either 'en' or 'ro'. Try again.\n ")
		}
	}
}

func writeInput(conn *net.TCPConn) {
	// Read from standard input
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Fprintf(conn, text+"\n")
}

func printOutput(conn *net.TCPConn) {

	inMessage := bufio.NewScanner(conn)

	for inMessage.Scan() {
		fmt.Println(inMessage.Text())
	}
}
