package main

import "net"
import "fmt"
import "bufio"
import "os"

func recv(conn net.Conn, c chan string) {

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// Handle error
			fmt.Print("Error: Reading")
		}
		c <- msg
	}
}

func send(conn net.Conn, c chan string) {

	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// Handle error
			fmt.Print("Error: Sending")
		}
		fmt.Fprintf(conn, msg)
		c <- "> " + msg
	}
	return
}

func main() {

	// Check arguments
	if len(os.Args) != 4 {
		fmt.Print("Usage: go run main.go <NICK> <PASS> <USER>")
		return
	}

	// Connect to server
	conn, err := net.Dial("tcp", "irc.freenode.net:6667")

	if err != nil {
		// Handle Error
		fmt.Print("Error: Setting up connection\n")
		return
	}

	pass := "PASS " + os.Args[2] + "\n"
	nick := "NICK " + os.Args[1] + "\n"
	user := "USER 0 guest 0 * :" + os.Args[3] + "\n"

	// Send auth calls to socket
	fmt.Fprintf(conn, pass)
	fmt.Fprintf(conn, nick)
	fmt.Fprintf(conn, user)

	// Create a channel to communicate with goroutines
	c := make(chan string)
	// Start goroutines
	go recv(conn, c)
	go send(conn, c)

	// Read from the channel
	for {
		msg := <-c
		fmt.Println(msg)
	}

}
