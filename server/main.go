package main

import (
	"log"
	"net"
	"os"
	"syscall"

	"github.com/flyaways/netpoll"
	reuseport "github.com/libp2p/go-reuseport"
)

// unix.Socket  创建一个socket文件描述符
// unix.Bind    绑定一个本机IP:port到socket文件描述符上
// unix.Listen  监听是否有连接请求
// unix.Accept  获取一个建立好的连接
// unix.Connect 发起连接请求
// unix.Close   关闭连接

// unix.EPOLLIN  表示对应的文件描述符可以读
// unix.EPOLLOUT 表示对应的文件描述符可以写
// unix.EPOLLPRI 表示对应的文件描述符有紧急的数据可读
// unix.EPOLLERR 表示对应的文件描述符发生错误
// unix.EPOLLHUP 表示对应的文件描述符被挂断
// unix.EPOLLET  表示对应的文件描述符设定为edge模式

// unix.EPOLL_CTL_ADD 注册
// unix.EPOLL_CTL_MOD 修改
// unix.EPOLL_CTL_DEL 删除

func main() {
	p := netpoll.OpenPoll()
	defer p.Close()

	var l net.PacketConn
	l, err := reuseport.ListenPacket("udp", "127.0.0.1:8421")
	if err != nil {
		log.Println(err)
		return
	}

	var f *os.File
	switch pconn := l.(type) {
	case *net.UDPConn:
		f, err = pconn.File()
		if err != nil {
			log.Println(err)
			return
		}
	}

	fd := int(f.Fd())

	//syscall.SetNonblock(fd, true)

	p.AddRead(int(fd))

	log.Println(l.LocalAddr().String())

	log.Println("Wait", p.Wait(func(fd int) {
		buf := make([]byte, 1460)
		n, addr, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(SockaddrToAddr(addr).String(), string(buf[:n]))
		log.Println(SockaddrToAddr(addr).String(), syscall.Sendto(fd, buf[:n], 0, addr))

		syscall.CloseOnExec(fd)
	}))

	log.Println("done")
}

// SockaddrToAddr returns a go/net friendly address
func SockaddrToAddr(sa syscall.Sockaddr) net.Addr {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return &net.TCPAddr{
			IP:   append([]byte{}, sa.Addr[:]...),
			Port: sa.Port,
		}

	case *syscall.SockaddrInet6:
		var zone string
		if sa.ZoneId != 0 {
			if ifi, err := net.InterfaceByIndex(int(sa.ZoneId)); err == nil {
				zone = ifi.Name
			}
		}

		return &net.TCPAddr{
			IP:   append([]byte{}, sa.Addr[:]...),
			Port: sa.Port,
			Zone: zone,
		}
	case *syscall.SockaddrUnix:
		return &net.UnixAddr{Net: "unix", Name: sa.Name}
	}

	return nil
}
