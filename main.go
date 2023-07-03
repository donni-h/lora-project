package main

import (
	"bufio"
	"fmt"
	"go.bug.st/serial"
	"io"
	"log"
	"lora-project/protocol/messages"
	"lora-project/protocol/serial_handlers"
	"math/rand"
	"os"
	"time"
)

var commands = []string{
	"AT",
	"AT+SEND=3",
	"AT+SEND=A",
	"jasdjkjasjdjkasd",
}

func write(port serial.Port) {
	rnd := rand.Intn(len(commands))
	lineSep := "\r\n"
	msg := []byte(fmt.Sprintf("%s%s", commands[rnd], lineSep))
	n, err := port.Write(msg)
	if err != nil {
		fmt.Println("Error writing to serial port: ", err)
		return
	}
	fmt.Printf("Sent %v bytes: '%q'\n", n, msg)
	if rnd == 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("enter text: ")
		buf := make([]byte, 3)
		_, _ = io.ReadFull(reader, buf)
		fmt.Printf("AT,SENDING%q\n", lineSep)
		_, _ = port.Write(buf)
		fmt.Printf("AT,SENDED%q\n", lineSep)
	}
	time.Sleep(time.Second)
}

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

	data := []byte{'0', '0', '2', '1', '1', '1', '1', '0', '0', '0', '0', '4', '7', '6', '1', '3', '3', '3', '3', '2', '2', '2', '2', '0'}
	message, err := messages.Unmarshal(data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", message)
	byteSlice, _ := message.Marshal()
	fmt.Println(string(byteSlice))
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open("/home/Hannes/dev/ttyS21", mode)
	if err != nil {
		log.Fatal(err)
	}
	handler := serial_handlers.NewATHandler(port)
	handler.Run()
}
