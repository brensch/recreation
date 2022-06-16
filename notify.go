package main

import (
	"github.com/xconstruct/go-pushbullet"
)

func Yeet() {
	pb := pushbullet.New("o.gMQajO2iSWtnlJpUBdXDl55CRO9UNhLw")
	devs, err := pb.Devices()
	if err != nil {
		panic(err)
	}

	err = pb.PushNote(devs[0].Iden, "yeet!", "")
	if err != nil {
		panic(err)
	}

}
