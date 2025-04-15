package core

import (
	"bufio"
	"log"
	"net"
)

type TCP struct {
	Network string
	Ip      net.IP
	Port    int
}

func (t *TCP) GetTcpServer() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   t.Ip,
		Port: t.Port,
	})
	if err != nil {
		log.Printf("Listene failed，error：%s", err.Error())
		return
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Printf("TCP server close failed，error：%s", err.Error())
			return
		}
	}(listen)
	for {
		log.Println("等待客户端连接")
		// 监听客户端
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Accept failed，error：%s", err.Error())
			continue
		} else {
			log.Printf("Accept connect，Conn：%s，Client IP：%s", conn, conn.RemoteAddr().String())
		}
		go t.HandleTcpConnect(conn)
	}
}

func (t *TCP) HandleTcpConnect(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Close failed，error：%s", err.Error())
		}
	}(conn)
	for {
		buf := make([]byte, 1024)
		// 缓冲区
		reader := bufio.NewReader(conn)
		n, err := reader.Read(buf)
		if err != nil {
			log.Printf("Read from client failed，error：%s", err.Error())
			err := conn.Close()
			if err != nil {
				return
			}
			return
		}
		str := string(buf[:n])
		log.Printf("Recive message：%s", str)
		// 转发
		_, err = conn.Write([]byte(str))
		if err != nil {
			return
		}
	}
}
