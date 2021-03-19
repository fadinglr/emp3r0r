package cc

import (
	"log"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/jm33-m0/emp3r0r/core/internal/agent"
)

const SocketName = "/tmp/emp3r0r.socket"

var APIConn net.Conn

func HeadlessMain() {
	log.Printf("%s", color.CyanString("Starting emp3r0r API server"))
	APIListen()
}

// listen on a unix socket
// users can send commands to this socket as if they were
// using a console
func APIListen() {
	// if socket file exists
	if agent.IsFileExist(SocketName) {
		err := os.Remove(SocketName)
		if err != nil {
			CliPrintError("Failed to delete socket: %v", err)
			return
		}
	}

	l, err := net.Listen("unix", SocketName)
	if err != nil {
		CliPrintError("listen error:", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			CliPrintError("emp3r0r API: accept error:", err)
			return
		}
		APIConn = conn
		log.Printf("%s: %s", color.BlueString("emp3r0r got an API connection"), conn.RemoteAddr().String())
		processAPIReq(conn)
	}
}

// handle connections to our socket: echo whatever we get
func processAPIReq(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		CliPrintInfo("emp3r0r received \"%s\"", data)

		// deal with the command
		cmd := string(data)
		err = CmdHandler(cmd)
		if err != nil {
			CliPrintError("Command failed: %v", err)
		}
	}
}
