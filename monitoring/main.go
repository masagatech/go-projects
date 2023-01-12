package main

// #include <unistd.h>
import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"prnaik/go/tcpclient"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

var clClient *tcpclient.TCPClient

func main() {

	start()
	// fmt.Println("Total system memory: %s ", human.Bytes(memory.TotalMemory()))
	// fmt.Println("Free memory: %s", human.Bytes(memory.FreeMemory()))
	// percent, err := cpu.Percent(time.Second, false)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("CPU USAGE %d", int(math.Ceil(percent[0])))
	// cpus, _ := cpu.Counts(true)
	// fmt.Println("cpu count %d", cpus)

	// info, _ := cpu.Info()

	// fmt.Println("model ", info[0].ModelName)
	// fmt.Println("Core ", info[0].Cores)

	// formatter := "%-14s %7s %7s %7s %4s %s\n"
	// fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

	// parts, _ := disk.Partitions(true)
	// for _, p := range parts {
	// 	device := p.Mountpoint
	// 	s, _ := disk.Usage(device)

	// 	if s.Total == 0 {
	// 		continue
	// 	}

	// 	percent := fmt.Sprintf("%2.f%%", s.UsedPercent)

	// 	fmt.Printf(formatter,
	// 		s.Fstype,
	// 		human.Bytes(s.Total),
	// 		human.Bytes(s.Used),
	// 		human.Bytes(s.Free),
	// 		percent,
	// 		p.Mountpoint,
	// 	)
	// }

	// processes, err := process.Processes()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, p := range processes {

	// 	n, err := p.Name()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if strings.Contains(n, "Postman") {
	// 		mm, _ := p.CPUPercent()
	// 		usr, _ := p.Username()
	// 		fmt.Println("user ", usr)
	// 		fmt.Println("cpu ", mm)
	// 		sts, _ := p.Status()
	// 		fmt.Println("status ", sts)
	// 		break
	// 	} else {
	// 		// fmt.Println(n)
	// 	}
	// }

	// server := tcpserver.New(&tcpserver.Config{
	// 	Host: "127.0.0.1",
	// 	Port: "9005",
	// })

	// go server.Run()

	// time.Sleep(time.Second * 3)
	// // for i := 0; i < 10; i++ {
	// // 	fmt.Println("i", i)
	// clClient = &tcpclient.TCPClient{Host: "127.0.0.1", Port: "9005",
	// 	OnMessageReceived: onMessageReceived,
	// 	OnConnect:         onConnect}
	// clClient.Connect()
	//}

	// var client = tcpclient.TCPClient{Host: "127.0.0.1", Port: "9005", OnMessageReceived: onMessageReceived}
	// client.Connect()

}

func onMessageReceived(conn *net.TCPConn, message string) {
	fmt.Println(conn.LocalAddr().String(), message)
}

func onConnect(conn *net.TCPConn) {
	fmt.Println("connected")
	go startDataSending()
}

func startDataSending() {
	// fmt.Println("Total system memory: %s ", human.Bytes(memory.TotalMemory()))
	// fmt.Println("Free memory: %s", human.Bytes(memory.FreeMemory()))
	// percent, err := cpu.Percent(time.Second, false)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	//fmt.Println("CPU USAGE %d", int(math.Ceil(percent[0])))
	for {
		fmt.Println("sending")
		time.Sleep(time.Second * 10)
		b, err := json.Marshal(fillData())
		if err != nil {
			fmt.Println(err)
			return
		}
		clClient.SendData(string(b))
	}

}

func fillData() *Props {

	props := Props{}

	// ram data
	ram := Ram{}
	ram.TotalMemory = memory.TotalMemory()
	ram.FreeMemory = memory.FreeMemory()
	props.RAM = ram

	info, _ := cpu.Info()
	percent, err := cpu.Percent(time.Second, false)

	if err != nil {
		log.Fatal(err)
	}

	cpus, _ := cpu.Counts(true)

	cpup := CPU{}
	cpup.info = info[0]
	cpup.Count = cpus
	cpup.Percent = int(math.Ceil(percent[0]))

	props.CPU = cpup

	// formatter := "%-14s %7s %7s %7s %4s %s\n"
	// fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

	parts, _ := disk.Partitions(true)
	for _, p := range parts {
		device := p.Mountpoint
		s, _ := disk.Usage(device)

		if s.Total == 0 {
			continue
		}

		// percent := fmt.Sprintf("%2.f%%", s.UsedPercent)
		props.DISK = append(props.DISK, *s)
		// disks.Filesystem = s.Fstype
		// s.InodesFree
		// // fmt.Printf(formatter,
		// // 	s.Fstype,
		// // 	human.Bytes(s.Total),
		// // 	human.Bytes(s.Used),
		// // 	human.Bytes(s.Free),
		// // 	percent,
		// // 	p.Mountpoint,
		// // )
	}
	return &props

	// processes, err := process.Processes()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, p := range processes {

	// 	n, err := p.Name()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if strings.Contains(n, "Postman") {
	// 		mm, _ := p.CPUPercent()
	// 		usr, _ := p.Username()
	// 		fmt.Println("user ", usr)
	// 		fmt.Println("cpu ", mm)
	// 		sts, _ := p.Status()
	// 		fmt.Println("status ", sts)
	// 		break
	// 	} else {
	// 		// fmt.Println(n)
	// 	}
	// }
}

type Props struct {
	RAM  Ram
	CPU  CPU
	DISK []disk.UsageStat
}

type Ram struct {
	TotalMemory uint64
	FreeMemory  uint64
}

type CPU struct {
	info    cpu.InfoStat
	Count   int
	Percent int
}

var wg sync.WaitGroup

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
	client.Subscribe("hello", 1, func(c mqtt.Client, m mqtt.Message) {

		fmt.Println(string(m.Payload()))
	})
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
	//wg.Done()
}

func start() {

	wg.Add(1)
	var broker = "test.mosquitto.org"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))

	opts.SetClientID("go_mqtt_clients")
	// opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnReconnecting = func(c mqtt.Client, co *mqtt.ClientOptions) {
		fmt.Printf("Reconnecting....")
	}
	opts.AutoReconnect = true
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		time.Sleep(time.Second * 5)

	}
	wg.Wait()

	// for i := 0; i < 5; i++ {
	// 	time.Sleep(time.Second * 5)
	// 	b, err := json.Marshal(fillData())
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(client.Publish("hello", 1, true, string(b)))
	// }

}
