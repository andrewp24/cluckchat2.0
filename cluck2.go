package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

/*
 These are constants:

 ConnIPEn - The IP address for Enis' server.

 ConnPort - We are using port 8888

 ConnType - This client uses tcp.

 byteArraySize - The standard array size

 bufferSize - 2^16

 requestMLength - When the request doesn't contain a message, length of 0

 MOTDCode - Type uint16. The code is 2.

 registerCode - Type uint16. The code is 100.
*/
const (
	ConnIPEn              = "130.85.70.132"
	ConnPort              = "8888"
	ConnType              = "tcp"
	byteArraySize         = 2
	bufferSize            = 65536
	requestMLength uint16 = 0
	MOTDCode       uint16 = 2
	registerCode   uint16 = 100
)

func main() {
	// TODO an option for the user to pick if they want the MOTD or not?
	tcpAddr, err := net.ResolveTCPAddr(ConnType, ConnIPEn+":"+ConnPort)
	conn, err := net.DialTCP(ConnType, nil, tcpAddr) // Connect to server with tcp and handle any errors

	fmt.Println("Connecting to the server now.")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// look at this: https://youtu.be/icbFEmh7Ym0?t=708
	go printOutput(conn)
	registerUser(conn)
	time.Sleep(1 * time.Millisecond)
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
	mCodeArray := make([]byte, byteArraySize)   // Creates a byte type array of size 2 for the message type
	mLengthArray := make([]byte, byteArraySize) // Creates a byte type array of size 0 for the message length

	fmt.Println("Getting the MOTD")
	binary.BigEndian.PutUint16(mCodeArray, MOTDCode)         // Converts bytes to size 16 unsigned Int, saves to the code array
	binary.BigEndian.PutUint16(mLengthArray, requestMLength) // Converts bytes to size 16 unsigned Int, saves to the length array
	conn.Write(append(mCodeArray, mLengthArray...))          // Writes the bytes to the server
}

// Sends the server a request to register a user
func registerUser(conn *net.TCPConn) {
	rCodeArray := make([]byte, byteArraySize) // Creates a byte type array of size 2
	fmt.Println("Enter a username:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')                  // Takes in the user's input
	registerMLength := uint16(len(text) - 1)            // Takes the length of the user's input
	registerMLengthArray := make([]byte, byteArraySize) // creates an array the size of the length portion of the header
	username := []byte(text)                            // Creates the username byte array with the bytes of the input string
	username = username[:len(username)-1]               // Slices off the last character (new line character)

	binary.BigEndian.PutUint16(rCodeArray, registerCode)              // Converts bytes to size 16 unsigned Int
	binary.BigEndian.PutUint16(registerMLengthArray, registerMLength) // Converts bytes to size 16 unsigned Int, saves to the length array
	x := []byte{}                                                     // Creating a temp array
	x = append(rCodeArray, registerMLengthArray...)                   // Appends the header arrays to the temp array
	conn.Write(append(x, username...))                                // Appends the message array to the temp array. Writes the bytes to the server
}

// Creates and receives the buffer from the server
//
// Slices the different parts of the received packet
//
// The packet sent from the server is in this format:
//
// <unsigned_short: message_code> <unsigned_short: message_length> <char* string (not null terminated)>
func printOutput(conn *net.TCPConn) {
	buffer := make([]byte, bufferSize) // Creates the buffer
	inMessage := bufio.NewReader(conn) // Creates a new reader buffer for the connection

	// Endless loop to print the received data to the terminal
	for {
		n, err := inMessage.Read(buffer) // n is the number of bytes in the buffer
		mCode := buffer[0:2]             // Message code
		mLength := buffer[2:4]           // Message length
		mString := buffer[4:n]           // Message string

		fmt.Println(getPacketMessage(mString))
		getPacketCode(mCode)
		getPacketLength(mLength)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

// Takes the bytes and converts them to string
func getPacketMessage(mString []byte) string {
	s := string(mString[:])
	return s
}

// Takes the bytes and converts them to an unsigned short (uint16)
func getPacketCode(mCode []byte) uint16 {
	dCode := binary.BigEndian.Uint16(mCode)
	return dCode
}

// Takes the bytes and converts them to an unsigned short (uint16)
func getPacketLength(mLength []byte) uint16 {
	dLength := binary.BigEndian.Uint16(mLength)
	return dLength
}
