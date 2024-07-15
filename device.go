package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gosnmp/gosnmp"
)

type Device struct {
	snmpHandler      *gosnmp.GoSNMP
	deviceInterfaces []*deviceInterface
}

func NewDevice(ipaddress string, community string, port uint16, deviceInterfaces []*deviceInterface) *Device {
	target := &gosnmp.GoSNMP{
		Target:    ipaddress,
		Community: community,
		Port:      port,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(1) * time.Second,
	}
	return &Device{
		snmpHandler:      target,
		deviceInterfaces: deviceInterfaces,
	}
}

func (d *Device) Connect() error {
	err := d.snmpHandler.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) Close() error {
	err := d.snmpHandler.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) Update() {
	ifDescs, err := d.snmpHandler.BulkWalkAll(".1.3.6.1.2.1.2.2.1.2")
	if err != nil {
		log.Panicln(err)
	}

	ifOutOctets, err := d.snmpHandler.BulkWalkAll(".1.3.6.1.2.1.31.1.1.1.10")
	if err != nil {
		log.Panicln(err)
	}
	getTime := time.Now()

	for _, intf := range d.deviceInterfaces {
		var ifindex = 0
		for i := 0; i < len(ifDescs); i++ {
			if string(ifDescs[i].Value.([]uint8)) == intf.name {
				ifindex, err = strconv.Atoi(ifDescs[i].Name[21:])
				break
			}
		}

		for i := 0; i < len(ifOutOctets); i++ {
			if ifindex_l, _ := strconv.Atoi(ifOutOctets[i].Name[25:]); ifindex_l == ifindex {
				intf.Update(ifOutOctets[i].Value.(uint64), getTime)
				break
			}
		}
	}
}
