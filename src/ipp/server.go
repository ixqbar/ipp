package ipp

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type TServer struct {
	stopChannel chan bool
}

func NewTServer() *TServer {
	return &TServer{
		stopChannel: make(chan bool),
	}
}

func (obj *TServer) RunServer() error {
	raddr, err := net.ResolveUDPAddr("udp", ":9595")
	if err != nil {
		Logger.Print(err)
		return err
	}

	conn, err := net.ListenUDP("udp", raddr)
	if err != nil {
		Logger.Print(err)
		return err
	}

	defer conn.Close()

	go func() {
		select {
		case <-obj.stopChannel:
			Logger.Print("server mode catch stop signal")
			conn.SetDeadline(time.Now().Add(time.Second * 3))
		}
	}()

	data := make([]byte, 1024)

	for {
		n, c, err := conn.ReadFromUDP(data)
		if err != nil {
			Logger.Print(err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				break
			}
			continue
		}

		Logger.Printf("receive client %s packet %s", c.String(), string(data[:n]))

		fmt.Printf("%s\n\n", bytes.Replace(data[:n], []byte("|"), []byte("\n"), -1))
	}

	Logger.Print("server exit server mode")
	return nil
}

func (obj *TServer) RunClient() error {
	raddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9595")
	if err != nil {
		Logger.Print(err)
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		Logger.Print(err)
		return err
	}

	defer conn.Close()

	interval := time.NewTicker(time.Second * 5)
	defer func() {
		interval.Stop()
	}()

E:
	for {
		select {
		case <-obj.stopChannel:
			Logger.Print("client mode catch stop signal")
			break E
		case <-interval.C:
			var currentMachineIps string = strings.Join(GetCurrentMachineIps(), "|")
			Logger.Printf("client send packet %s", currentMachineIps)
			conn.Write([]byte(currentMachineIps))
		}
	}

	Logger.Print("server exit client mode")
	return nil
}

func (obj *TServer) Stop() {
	obj.stopChannel <- true
}

func Run(runAsServer bool) {
	errorChan := make(chan bool)
	signalChan := make(chan os.Signal)

	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	server := NewTServer()

	go func() {
		Logger.Printf("signal goroutine running")

		select {
		case <-signalChan:
			Logger.Print("signal goroutine catch exit signal")
			server.Stop()
		case <-errorChan:
			Logger.Print("signal goroutine catch error signal")
		}

		Logger.Printf("signal goroutine exit")
	}()

	Logger.Print("server start running")
	if runAsServer {
		err := server.RunServer()
		if err != nil {
			errorChan <- true
		}
	} else {
		err := server.RunClient()
		if err != nil {
			errorChan <- true
		}
	}

	Logger.Print("server stop running")
}
