package tcpclient

import (
	"fmt"
	"net"
	"os"
)

type TCPClient struct {
	Host              string
	Port              string
	OnMessageReceived func(conn *net.TCPConn, message string)
	OnConnect         func(conn *net.TCPConn)
	conn              *net.TCPConn
}

func (cl *TCPClient) Connect() {

	servAddr := fmt.Sprintf("%s:%s", cl.Host, cl.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	cl.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	if cl.OnConnect != nil {
		cl.OnConnect(cl.conn)
	}
	cl.listner()
}

func (cl *TCPClient) SendData(message string) {
	// if cl.conn != nil {
	_, err := cl.conn.Write([]byte(message + "\n"))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	//}
}

func (cl *TCPClient) Close() {
	if cl.conn != nil {
		cl.conn.Close()
	}
}

func (cl *TCPClient) listner() {

	received := make([]byte, 1024)
	for {
		_, err := cl.conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}
		if cl.OnMessageReceived != nil {
			cl.OnMessageReceived(cl.conn, string(received))
		}
	}
}

// func main() {
// 	strEcho := "Halo"
// 	servAddr := "localhost:9005"
// 	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
// 	if err != nil {
// 		println("ResolveTCPAddr failed:", err.Error())
// 		os.Exit(1)
// 	}

// 	conn, err := net.DialTCP("tcp", nil, tcpAddr)
// 	if err != nil {
// 		println("Dial failed:", err.Error())
// 		os.Exit(1)
// 	}

// 	_, err = conn.Write([]byte(strEcho))
// 	if err != nil {
// 		println("Write to server failed:", err.Error())
// 		os.Exit(1)
// 	}

// 	println("write to server = ", strEcho)

// 	reply := make([]byte, 1024)

// 	_, err = conn.Read(reply)
// 	if err != nil {
// 		println("Write to server failed:", err.Error())
// 		os.Exit(1)
// 	}

// 	println("reply from server=", string(reply))

// 	conn.Close()
// }
