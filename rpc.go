package main

import (
	"net"
	"net/rpc"
)

// ServeRPC exposes the methods of a DX7 over RPC.
func ServeRPC(dx7 *DX7, ip string) error {
	server := rpc.NewServer()
	laddr, err := net.ResolveTCPAddr("tcp", ip+":0")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	defer l.Close()

	if err := server.RegisterName("YamahaDX7", dx7); err != nil {
		return err
	}
	logger.Printf("serving rpc on %s\n", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		server.ServeConn(conn)
	}
	return nil
}
