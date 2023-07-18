package main

import (
	"bufio"
	"fmt"
	"go.bug.st/serial"
	"log"
	"lora-project/protocol/messages"
	"lora-project/protocol/routing"
	"lora-project/serial_handlers"
	"os"
)

func main() {
	/*	rand.Seed(time.Now().UnixNano())
		app := &cli.App{
			Name:  "boom",
			Usage: "make an explosive entrance",
			Action: func(*cli.Context) error {
				fmt.Println("boom! I say!")
				return nil
			},
		}

		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}
		mode := &serial.Mode{
			BaudRate: 115200,
		}
		port, err := serial.Open("/home/Hannes/dev/ttyS21", mode)
		if err != nil {
			log.Fatal(err)
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sig
			fmt.Println("Interrupt signal received, closing serial port...")
			port.Close()
			os.Exit(0)
		}()
		write(port)
		scanner := bufio.NewScanner(port)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("received: " + line)
			write(port)
		}*/

	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err := serial.Open("/dev/ttyUSB0", mode)
	defer func(port serial.Port) {
		err := port.Close()
		if err != nil {

		}
	}(port)
	if err != nil {
		log.Fatal(err)
	}
	handler := serial_handlers.NewATHandler(port)
	aodv := routing.NewAODV(handler)
	aodv.Run()

	go func() {
		for data := range aodv.IncomingDataQueue {
			fmt.Printf("Received message: %s from %s\n", string(data.Payload), data.OriginatorAddress.String())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter destination address: ")
		scanned := scanner.Scan()
		if !scanned {
			fmt.Println("Failed to scan address!")
			return
		}

		destination := scanner.Text()
		var address messages.Address
		if address.UnmarshalText([]byte(destination)) != nil {
			fmt.Println("Invalid address..")
			continue
		}

		fmt.Print("Enter message: ")
		scanned = scanner.Scan()

		if !scanned {
			fmt.Println("Failed to scan message!")
			return
		}

		message := scanner.Text()

		aodv.SendData(message, address)

	}
}
