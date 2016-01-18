package main

import (
	"net"
	"net/rpc"
)

// ServeRPC exposes the methods of a DX7 over RPC.
func ServeRPC(dx7 *DX7, ip string) error {
	server := rpc.NewServer()
	laddr, err := net.ResolveUDPAddr("udp", ip+":0")
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	if err := server.Register(dx7); err != nil {
		return err
	}
	logger.Printf("serving rpc on %s\n", laddr)
	server.ServeConn(conn)
	return nil
}
