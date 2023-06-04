package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"go.bug.st/serial"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
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
	rand.Seed(time.Now().UnixNano())
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
	}
}
