// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/gvcgo/gocui"
)

func main() {
	opt := gocui.NewGuiOpts{
		OutputMode: gocui.OutputNormal,
	}
	g, err := gocui.NewGui(opt)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("size", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		if _, err := g.SetCurrentView("size"); err != nil {
			return err
		}
	}
	v.Clear()
	fmt.Fprintf(v, "%d, %d", maxX, maxY)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
