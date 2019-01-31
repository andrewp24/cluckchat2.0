package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
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

 whoAmICode - Type uint16. The code is 110

 registerCode - Type uint16. The code is 100.

 failureCode - Type uint16. The code is 11.
*/
const (
	connIPEn              = "130.85.70.132"
	connPort              = "8888"
	connType              = "tcp"
	byteArraySize         = 2
	bufferSize            = 65536
	requestMLength uint16 = 0
	mOTDCode       uint16 = 2
	whoAmICode     uint16 = 110
	registerCode   uint16 = 100
	failureCode    uint16 = 11
)

var m = sync.Mutex{}

func main() {
	// TODO an option for the user to pick if they want the MOTD or not?
	tcpAddr, err := net.ResolveTCPAddr(connType, connIPEn+":"+connPort)
	conn, err := net.DialTCP(connType, nil, tcpAddr) // Connect to server with tcp and handle any errors

	fmt.Println("Connecting to the server now.")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go printOutput(conn)
	registerUser(conn)
	whoAmI(conn)
	for {
		//	writeInput(conn)

	}
}

// Read from standard input
func writeInput(conn *net.TCPConn) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	m.Lock()
	fmt.Fprintf(conn, text+"\n")
	m.Unlock()
}

func whoAmI(conn *net.TCPConn) {
	wCodeArray := make([]byte, byteArraySize)   // Creates a byte type array of size 2 for the message type
	wLengthArray := make([]byte, byteArraySize) // Creates a byte type array of size 0 for the message length

	//fmt.Println("Getting the MOTD")
	binary.BigEndian.PutUint16(wCodeArray, whoAmICode)       // Converts bytes to size 16 unsigned Int, saves to the code array
	binary.BigEndian.PutUint16(wLengthArray, requestMLength) // Converts bytes to size 16 unsigned Int, saves to the length array
	m.Lock()
	conn.Write(append(wCodeArray, wLengthArray...)) // Writes the bytes to the server
	code, _, message := getPackage(conn)

	if code == failureCode {
		fmt.Println("You haven't been registered.")
	} else {
		fmt.Println("You are: " + message)
	}
	m.Unlock()

}

// Requests the message of the day from the server
//
// Currently, the code for the MOTD is 2
func askForMotd(conn *net.TCPConn) {
	mCodeArray := make([]byte, byteArraySize)   // Creates a byte type array of size 2 for the message type
	mLengthArray := make([]byte, byteArraySize) // Creates a byte type array of size 0 for the message length

	//fmt.Println("Getting the MOTD")
	binary.BigEndian.PutUint16(mCodeArray, mOTDCode)         // Converts bytes to size 16 unsigned Int, saves to the code array
	binary.BigEndian.PutUint16(mLengthArray, requestMLength) // Converts bytes to size 16 unsigned Int, saves to the length array
	conn.Write(append(mCodeArray, mLengthArray...))          // Writes the bytes to the server
	_, _, message := getPackage(conn)                        // Retrieves the message of the package
	fmt.Println(message)
}

// Sends the server a request to register a user
//
// If the user successfully registers, the MOTD is requested.
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
	m.Lock()
	conn.Write(append(x, username...)) // Appends the message array to the temp array. Writes the bytes to the server
	code, _, _ := getPackage(conn)     // Retreives the code of the package
	if code == failureCode {
		println("Failed to register the user.")
	} else {
		askForMotd(conn)
	}
	m.Unlock()
}

// getPackage takes the connection and returns the code, length, and message of the package
func getPackage(conn *net.TCPConn) (uint16, uint16, string) {
	buffer := make([]byte, bufferSize)
	recMessage := bufio.NewReader(conn)
	n, err := recMessage.Read(buffer) // n is the number of bytes in the buffer
	mCode := buffer[0:2]              // Message code
	mLength := buffer[2:4]            // Message length
	mString := buffer[4:n]            // Message string
	packetCode := getPacketCode(mCode)
	packetLength := getPacketLength(mLength)
	packetMessage := getPacketMessage(mString)
	if err != nil {
		fmt.Println(err)
	}
	return packetCode, packetLength, packetMessage
}

// Creates and receives the buffer from the server
//
// Slices the different parts of the received packet
//
// The packet sent from the server is in this format:
//
// <unsigned_short: message_code> <unsigned_short: message_length> <char* string (not null terminated)>
// look at how to do locks
func printOutput(conn *net.TCPConn) {
	buffer := make([]byte, bufferSize) // Creates the buffer
	inMessage := bufio.NewReader(conn) // Creates a new reader buffer for the connection

	// Endless loop to print the received data to the terminal
	for {
		m.Lock()
		if inMessage.Buffered() > 0 {
			fmt.Println(inMessage.Buffered())
			_, err := inMessage.Read(buffer) // n is the number of bytes in the buffer
			code, _, message := getPackage(conn)
			m.Unlock()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(code, message)
		} else {
			m.Unlock()
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
