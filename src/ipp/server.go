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
		stopChannel:make(chan bool, 0),
	}
}

func (obj *TServer) RunServer() {
	raddr, err := net.ResolveUDPAddr("udp", ":9595")
	if err != nil {
		Logger.Print(err)
		return
	}

	conn, err := net.ListenUDP("udp", raddr)
	if err != nil {
		Logger.Print(err)
		return
	}

	defer conn.Close()

	go func() {
		for {
			select {
			case <-obj.stopChannel:
				Logger.Print("server mode catch stop signal")
				conn.SetDeadline(time.Now().Add(time.Second * 3))
				return
			}
		}
	}()

	data := make([]byte, 1024)

	for {
		n,c,err := conn.ReadFromUDP(data)
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
}

func (obj *TServer) RunClient()  {
	raddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9595")
	if err != nil {
		Logger.Print(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		Logger.Print(err)
		return
	}

	defer conn.Close()

	interval := time.NewTicker(time.Second * 5)
	defer func() {
		interval.Stop()
	}()

E:
	for {
		select {
		case <- obj.stopChannel:
			Logger.Print("client mode catch stop signal")
			break E
		case <- interval.C:
			var currentMachineIps string = strings.Join(GetCurrentMachineIps(), "|")
			Logger.Printf("client send packet %s", currentMachineIps)
			conn.Write([]byte(currentMachineIps))
		}
	}

	Logger.Print("server exit client mode")
}

func (obj *TServer) Stop() {
	obj.stopChannel <- true
}

func Run(runAsServer bool)  {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	server := NewTServer()

	go func() {
		<-sigs
		Logger.Print("server catch exit signal")
		server.Stop()
	}()

	Logger.Print("server start running")
	if runAsServer {
		server.RunServer()
	} else {
		server.RunClient()
	}

	Logger.Print("server stop running")
}