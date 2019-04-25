package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

func main() {
	servoTest()
}

func ledTest() {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, "7")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}

func servoTest() {
	r := raspi.NewAdaptor()
	servo := gpio.NewServoDriver(r, "12")

	var angle uint8 = 0

	work := func() {
		gobot.Every(1*time.Second, func() {
			fmt.Printf("serve angle:%d\n", angle)
			if err := servo.Move(angle); err != nil {
				fmt.Println(err)
			}
			angle++
		})
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{r},
		[]gobot.Device{servo},
		work,
	)

	robot.Start()
}
