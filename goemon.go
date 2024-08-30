package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
)

func main() {
	var config Config
	var wg sync.WaitGroup
	LoadConfigs("config.yaml", &config)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			log.Panicln(err)
		}
		v.Title = " " + config.Title + " "
		v.SetOrigin(0, 0)
		fmt.Fprintf(v, "%s", config.Map)
	}
	if err != nil {
		log.Panicln(err)
	}

	devices := []*Device{}
	for _, device := range config.Devices {
		intfs := []*deviceInterface{}
		for _, intf := range device.Interfaces {
			intfs = append(intfs, NewDeviceInterface(g, intf.Name, intf.PositionX, intf.PositionY))
		}
		devices = append(devices, NewDevice(device.IpAddress, device.Community, device.Port, intfs))
	}

	for _, device := range devices {
		wg.Add(1)
		go func(d *Device) {
			defer wg.Done()
			d.Connect()
		}(device)
		defer device.Close()
	}
	wg.Wait()

	go func() {
		for {
			g.Update(func(g *gocui.Gui) error {
				for _, device := range devices {
					wg.Add(1)
					go func(d *Device) {
						defer wg.Done()
						d.Update()
					}(device)
				}
				wg.Wait()
				return nil
			})
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
