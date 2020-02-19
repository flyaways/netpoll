// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build darwin netbsd freebsd openbsd dragonfly

package netpoll

import (
	"syscall"
)

type Poll struct {
	fd      int
	changes []syscall.Kevent_t
	//connections *Shards
}

// OpenPoll ...
func OpenPoll() *Poll {
	p := new(Poll)
	//p.connections = NewShards(1024)
	fd, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}
	p.fd = fd

	return p
}

// Close ...
func (p *Poll) Close() error {
	return syscall.Close(p.fd)
}

// Wait ...
func (p *Poll) Wait(iter func(int)) error {
	events := make([]syscall.Kevent_t, 128)
	for {
		n, err := syscall.Kevent(p.fd, p.changes, events, nil)
		if err != nil {
			return err
		}
		p.changes = p.changes[:0]

		for i := 0; i < n; i++ {
			if fd := int(events[i].Ident); fd != 0 {
				iter(fd)
			}
		}
	}
}

// AddRead ...
func (p *Poll) AddRead(fd int) {
	p.changes = append(p.changes,
		syscall.Kevent_t{
			Ident: uint64(fd), Flags: syscall.EV_ADD, Filter: syscall.EVFILT_READ,
		},
	)
}

// AddReadWrite ...
func (p *Poll) AddReadWrite(fd int) {
	p.changes = append(p.changes,
		syscall.Kevent_t{
			Ident: uint64(fd), Flags: syscall.EV_ADD, Filter: syscall.EVFILT_READ,
		},
		syscall.Kevent_t{
			Ident: uint64(fd), Flags: syscall.EV_ADD, Filter: syscall.EVFILT_WRITE,
		},
	)
}

// ModRead ...
func (p *Poll) ModRead(fd int) {
	p.changes = append(p.changes, syscall.Kevent_t{
		Ident: uint64(fd), Flags: syscall.EV_DELETE, Filter: syscall.EVFILT_WRITE,
	})
}

// ModReadWrite ...
func (p *Poll) ModReadWrite(fd int) {
	p.changes = append(p.changes, syscall.Kevent_t{
		Ident: uint64(fd), Flags: syscall.EV_ADD, Filter: syscall.EVFILT_WRITE,
	})
}

// ModDetach ...
func (p *Poll) ModDetach(fd int) {
	p.changes = append(p.changes,
		syscall.Kevent_t{
			Ident: uint64(fd), Flags: syscall.EV_DELETE, Filter: syscall.EVFILT_READ,
		},
		syscall.Kevent_t{
			Ident: uint64(fd), Flags: syscall.EV_DELETE, Filter: syscall.EVFILT_WRITE,
		},
	)
}
