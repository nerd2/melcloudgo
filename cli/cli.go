package main

import (
	"github.com/nerd2/melcloudgo"
	"log"
	"os"
)

func main() {
	username := os.Args[1]
	password := os.Args[2]

	nh := melcloudgo.NewMelCloud(&melcloudgo.Options{Username: username, Password: password})
	_, err := nh.Login()
	if err != nil {
		log.Fatalln(err.Error())
	}

	data, err := nh.ListDevices()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if len(data) == 0 {
		log.Fatalln("No data returned")
	}

	for _, dev := range data[0].Structure.Devices {
		log.Printf("%s: %f/%f/%f/%f", dev.DeviceName, dev.Device.FlowTemperature, dev.Device.ReturnTemperature, dev.Device.TankWaterTemperature, dev.Device.OutdoorTemperature)
	}
}
