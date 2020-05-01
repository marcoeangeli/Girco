package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"net"
	"os"
)

var conn net.Conn

func sendMessage(g *gocui.Gui, v *gocui.View) error {
	message := v.Buffer()
	feedView, _ := g.View("feed")

	fmt.Fprint(feedView, message)
	fmt.Fprintf(conn, message)

	v.Clear()
	v.SetCursor(0, 0)

	return nil
}

func recv(g *gocui.Gui, conn net.Conn) error {

	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("feed")
			if err != nil {
				// handle error
			}
			fmt.Fprintln(v, msg[:len(msg)-2])
			return nil
		})
	}
}

func keybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("textBox", gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	return nil
}

func layout(g *gocui.Gui) error {
	margin := 1
	textBoxHeight := 2
	maxX, maxY := g.Size()
	if v, err := g.SetView("feed", margin, margin, maxX-margin, maxY-margin-textBoxHeight-margin); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}
	if v, err := g.SetView("icon", margin, maxY-margin-textBoxHeight, margin+textBoxHeight, maxY-margin); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, ">")
	}
	if v, err := g.SetView("textBox", margin+textBoxHeight+margin, maxY-margin-textBoxHeight, maxX-margin, maxY-margin); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		//fmt.Fprintln(v, " > ")
		if _, err := g.SetCurrentView("textBox"); err != nil {
			return err
		}
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {

	if len(os.Args) != 3 {
		fmt.Print("Usage: go run main.go <NICK> <PASS>")
		return
	}

	network := "irc.freenode.net"
	port := "6667"
	pass := os.Args[2]
	nick := os.Args[1]
	user := "unknown"

	// Set up interface
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.SetManagerFunc(layout)
	keybindings(g)

	// Connect to the server
	tmp, err := net.Dial("tcp", network+":"+port)
	conn = tmp

	if err != nil {
		// Handle Error
		fmt.Print("Error: Setting up connection\n")
		return
	}

	pass_msg := "PASS " + pass + "\n"
	nick_msg := "NICK " + nick + "\n"
	user_msg := "USER 0 guest 0 * :" + user + "\n"

	// Send auth calls to socket
	fmt.Fprintf(conn, pass_msg)
	fmt.Fprintf(conn, nick_msg)
	fmt.Fprintf(conn, user_msg)

	// Start the goroutines
	go recv(g, conn)

	// Start the interface
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
