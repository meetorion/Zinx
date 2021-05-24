package znet

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	Name string
	IPversion string
	IP string
	Port int

}

//将实例化一个服务器封装成函数,参数为服务器的名字；返回值是一个IServer接口
func NewServer(name string) *Server {
	s := &Server {
		Name:name,
		IPversion: "tcp4",
		IP:"0.0.0.0",
		Port:7777,
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	//用一个匿名协程去做监听服务
	go func() {
		//1. 获取TCP的地址addr
		addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp address err: ", err)
			return
		}
		//2. 监听服务器地址
		listenner, err := net.ListenTCP(s.IPversion, addr)
		if err != nil {
			fmt.Println("listen", s.IPversion,err)
			return
		}
		fmt.Println("Start Zinx Server ", s.Name, " succ, now listenning...")
		//3. 启动网络连接服务
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err",err)
				continue
			}
			//TODO 设置服务器最大连接控制，如果超过最大连接，则关闭新的连接
			//TODO 处理新连接请求的“业务”方法

			//暂时做一个最大512字节的回显服务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt,err := conn.Read(buf)
					if err != nil {
						fmt.Println("rec buf err ", err)
						continue
					}
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write buf back err ", err)
						continue
					}
				}

			}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Server,name ",s.Name)
}

func (s *Server) Server() {
	s.Start()

	//一直运行，否则主协程结束，匿名协程也会结束
	for {
		time.Sleep(time.Second * 10)
	}
}