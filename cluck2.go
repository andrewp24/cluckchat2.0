package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

/*
 These are constants:

 ConnIPEn - The IP address for Enis' server.

 ConnPort - We are using port 8888

 ConnType - This client uses tcp. can be udp as well.
*/
const (
	ConnIPEn = "130.85.70.132"
	ConnPort = "8888"
	ConnType = "tcp"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, ConnIPEn+":"+ConnPort)
	conn, err := net.DialTCP(ConnType, nil, tcpAddr) // Connect to server with tcp and handle any errors

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go printOutput(conn)
	askForMotd(conn)
	for {
		writeInput(conn)
	}
}

// Read from standard input
func writeInput(conn *net.TCPConn) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	fmt.Fprintf(conn, text+"\n")
}

// Requests the message of the day from the server
//
// Currently, the code for the MOTD is 2
func askForMotd(conn *net.TCPConn) {
	var mCode uint16 = 2
	arraySize := 2
	mCodeArray := make([]byte, arraySize) // Creates a byte type array of size 2
	writer := bufio.NewWriter(conn)       // Creates a writer to right to the server

	binary.BigEndian.PutUint16(mCodeArray, mCode) // Converts bytes to size 16 unsigned Int
	writer.Write(mCodeArray)
	fmt.Println("done")
}

// Creates and receives the buffer from the server
//
// Slices the different parts of the received packet
//
// The packet sent from the server is in this format:
//
// <unsigned_short: message_code> <unsigned_short: message_length> <char* string (not null terminated)>
func printOutput(conn *net.TCPConn) {
	bufferSize := 65536                // 2^16
	buffer := make([]byte, bufferSize) // Creates the buffer5536 ==
	inMessage := bufio.NewReader(conn) // Creates a new reader buffer for the connection

	// Endless loop to print the received data to the terminal
	for {
		n, err := inMessage.Read(buffer) // n is the number of bytes in the buffer
		mCode := buffer[0:2]             // Message code
		mLength := buffer[2:4]           // Message length
		mString := buffer[4:n]           // Message string

		getPacketMessage(mString)
		getPacketCode(mCode)
		getPacketLength(mLength)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

// Takes the bytes and converts them to string
func getPacketMessage(mString []byte) {
	s := string(mString[:])
	fmt.Println(s)
}

// Takes the bytes and converts them to an unsigned short (uint16)
func getPacketCode(mCode []byte) {
	dCode := binary.BigEndian.Uint16(mCode)
	fmt.Println(dCode)
}

// Takes the bytes and converts them to an unsigned short (uint16)
func getPacketLength(mLength []byte) {
	dLength := binary.BigEndian.Uint16(mLength)
	fmt.Println(dLength)
}
