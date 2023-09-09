package main

import (
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"
)

var (
	ip         = flag.String("addr", "255.255.255.255", "ip")
	serverMode = flag.Bool("s", false, "server mode")
	port       = flag.Int("port", 9012, "port")
	key        = flag.String("key", "Knock", "key")
	knock      = new(Knock)
)

type Knock struct {
	servieLogger service.Logger
}

func (k *Knock) Start(s service.Service) error {
	k.servieLogger.Info("service start")
	go run()
	return nil
}

func (k *Knock) Stop(s service.Service) error {
	return nil
}

func main() {
	flag.Parse()
	config := service.Config{
		Name:        "knock-go",
		DisplayName: "knock-go service",
		Arguments:   []string{"-s"},
	}
	s, err := service.New(knock, &config)
	if err != nil {
		panic(err)
	}
	knock.servieLogger, _ = s.Logger(make(chan<- error))
	switch flag.Arg(0) {
	case "uninstall":
		err = s.Uninstall()
	case "install":
		err = s.Install()
	case "start":
		err = s.Start()
	case "stop":
		err = s.Stop()
	default:
		run()
	}
	if err != nil {
		slog.Error(err.Error())
	}

}
func run() {
	if *serverMode {
		Listener()
		return
	}
	Sender()
}
func Listener() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: *port,
	})
	if err != nil {
		panic(err)
	}
	var tmp = make([]byte, 512)
	slog.Info("running")
	for {
		n, addr, err := listen.ReadFrom(tmp)
		if err != nil {
			continue
		}
		if string(tmp[:n]) == *key {
			slog.Info("knock", "addr", addr.String())
			hostname, _ := os.Hostname()
			listen.WriteTo([]byte("Answer-"+hostname), addr)
		}
	}
}

func Sender() {
	listen, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}
	_, err = listen.WriteToUDP([]byte(*key), &net.UDPAddr{
		IP:   net.ParseIP(*ip),
		Port: *port,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		var tmp = make([]byte, 512)
		fmt.Println("IP", "\t\t", "Hostname")
		for {
			n, addr, err := listen.ReadFrom(tmp)
			if err != nil {
				continue
			}
			if s := strings.SplitN(string(tmp[:n]), "-", 2); len(s) == 2 && s[0] == "Answer" {
				fmt.Println(strings.Split(addr.String(), ":")[0], "\t", s[1])
			}
		}
	}()
	time.Sleep(5 * time.Second)
}
