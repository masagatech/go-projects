package tcpclient

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

type TCPRedisClient struct {
	Host              string
	Port              string
	OnMessageReceived func(channel string, message interface{})
	OnConnect         func()
	rdb               *redis.Client
	ctx               *context.Context
}

func (cl *TCPRedisClient) Connect() {

	cl.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cl.Host, cl.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// servAddr := fmt.Sprintf("%s:%s", cl.Host, cl.Port)
	// tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	// if err != nil {
	// 	println("ResolveTCPAddr failed:", err.Error())
	// 	os.Exit(1)
	// }

	// cl.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	// if err != nil {
	// 	println("Dial failed:", err.Error())
	// 	os.Exit(1)
	// }
	if cl.OnConnect != nil {
		cl.OnConnect()
	}
	// cl.listner()
}

func (cl *TCPRedisClient) SendData(channel string, message interface{}) {
	cl.rdb.Publish(*cl.ctx, channel, message)
}

func (cl *TCPRedisClient) Close() {
	cl.rdb.Close()
}

func (cl *TCPRedisClient) Listen(channel string) {
	pubsub := cl.rdb.Subscribe(*cl.ctx, channel)
	// Close the subscription when we are done.
	defer pubsub.Close()

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
