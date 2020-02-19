// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build linux

package netpoll

import (
	"syscall"
	"unsafe"
)

// Poll ...
type Poll struct {
	fd int
	//connections *Shards
}

// OpenPoll ...
func OpenPoll() *Poll {
	p := new(Poll)
	//p.connections = NewShards(1024)
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	p.fd = fd
	return p
}

// Close ...
func (p *Poll) Close() error {
	if err := syscall.Close(p.wfd); err != nil {
		return err
	}
	return syscall.Close(p.fd)
}

// Wait ...
func (p *Poll) Wait(iter func(int)) error {
	events := make([]syscall.EpollEvent, 64)
	for {
		n, err := syscall.EpollWait(p.fd, events, 100)
		if err != nil && err != syscall.EINTR {
			return err
		}

		for i := 0; i < n; i++ {
			if fd := int(events[i].Fd); fd != 0 {
				iter(fd)
			}
		}
	}
}

// AddReadWrite ...
func (p *Poll) AddReadWrite(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}

// AddRead ...
func (p *Poll) AddRead(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		panic(err)
	}
}

// ModRead ...
func (p *Poll) ModRead(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_MOD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		panic(err)
	}
}

// ModReadWrite ...
func (p *Poll) ModReadWrite(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_MOD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}

// ModDetach ...
func (p *Poll) ModDetach(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_DEL, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}
