package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func sendMessage(g *gocui.Gui, v *gocui.View) error {
	message := v.Buffer()
	feedView, _ := g.View("feed")
	fmt.Fprint(feedView, message)
	v.Clear()
	v.SetCursor(0, 0)

	return nil
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
	if _, err := g.SetView("feed", margin, margin, maxX-margin, maxY-margin-textBoxHeight-margin); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if v, err := g.SetView("icon", margin, maxY-margin-textBoxHeight, margin+textBoxHeight, maxY-margin); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, ">")
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
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	keybindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
