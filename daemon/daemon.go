package daemon

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/takama/daemon"
	"github.com/toasterson/opencloud/common"
)

type Server struct {
	daemon.Daemon
	Port        string
	Name        string
	Description string
}

// Manage by daemon commands or run the daemon
func RunService(service Server, clientFn func() (string, error)) (string, error) {

	usage := "Usage: imaged install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		case "client":
			return clientFn()
		default:
			return usage, nil
		}
	}

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", service.Port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)
	defer func() {
		listener.Close()
		fmt.Println("Listener closed")
	}()

	rpc.Register(service)
	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			common.Stdlog.Println("Got signal:", killSignal)
			common.Stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	common.Stdlog.Println("Got connection from: " + client.RemoteAddr().String())

	// Close connection when this function ends
	defer func() {
		common.Stdlog.Println("Closing connection...")
		client.Close()
	}()
	rpc.ServeConn(client)
}
