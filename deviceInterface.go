package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	hmnize "github.com/dustin/go-humanize"
	"github.com/jroimartin/gocui"
)

const trafficValueLength = 6

type deviceInterface struct {
	name             string
	OutputTrafficBps float64
	lastOctets       uint64
	lastGetDateime   time.Time
	view             *gocui.View
	positionX        int
	positionY        int
}

func NewDeviceInterface(gui *gocui.Gui, name string, positionX int, positionY int) *deviceInterface {
	v, err := gui.SetView(name+"_"+strconv.Itoa(positionX)+"_"+strconv.Itoa(positionY), positionX, positionY, positionX+trafficValueLength+1, positionY+2)
	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
	v.Frame = false
	return &deviceInterface{
		name:       name,
		lastOctets: 0,
		view:       v,
		positionX:  positionX,
		positionY:  positionY,
	}
}

func (i *deviceInterface) Update(newOutOctets uint64, getTime time.Time) {
	i.view.Clear()
	if i.lastOctets != 0 {
		if i.lastOctets == newOutOctets {
			if getTime.Sub(i.lastGetDateime).Seconds() > 10 {
				i.OutputTrafficBps = 0
			}
		} else {
			i.OutputTrafficBps = float64(newOutOctets-i.lastOctets) * 8 / getTime.Sub(i.lastGetDateime).Seconds()
		}
		hOutTraffic, unit := hmnize.ComputeSI(i.OutputTrafficBps)
		fmt.Fprintf(i.view, "%5.1f%s", hOutTraffic, unit)
	} else {
		// init state
		fmt.Fprintf(i.view, " - ")
	}

	i.lastOctets = newOutOctets
	i.lastGetDateime = getTime
}
