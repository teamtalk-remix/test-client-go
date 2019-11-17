package dispatcher

/*
import (
	"fmt"
	"github.com/shiywang/GoTalk/pkg/pduhandler"
	"net"
	"os"
)

type Dispatcher struct {
	recvConn []net.Conn
	sendConn []net.Conn
	l        net.Listener
	recvSeq  int
	sendSeq  int
}

const (
	CONN_HOST     = "localhost"
	CONN_PORT     = "3333"
	CONN_TYPE     = "tcp"
	READ_BUF_SIZE = 2048
)

func (d *Dispatcher) Listen() {
	// Listen for incoming connections.
	var err error
	d.l, err = net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
}

func (d *Dispatcher) Start(p pduhandler.PduHandler) {

	for {
		// Listen for an incoming connection.
		conn, err := d.l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		d.recvConn = append(d.recvConn, conn)
		// Handle connections in a new goroutine.
		go d.handleRequest(p)
	}
}

// Handles incoming requests.
func (d *Dispatcher) handleRequest(p pduhandler.PduHandler) {
	// Make a buffer to hold incoming data.
	recvBuf := make([]byte, READ_BUF_SIZE)

	for {
		// Read the incoming connection into the buffer.
		n, err := d.recvConn[d.recvSeq].Read(recvBuf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		if n == READ_BUF_SIZE {
			go p.HandlePdu(recvBuf, d)
		}
	}

}

func (d *Dispatcher) handleResponse() {
	sendBuf := make([]byte, READ_BUF_SIZE)

	for {
		n, err := d.sendConn[d.sendSeq].Write(sendBuf)
		if err != nil {
			fmt.Println("Error writing:", err.Error())
		}

		if n == READ_BUF_SIZE {
			break
		}
	}

}
*/
