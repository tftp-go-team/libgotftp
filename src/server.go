package tftp

import (
	"fmt"
	"net"
	"regexp"
)

type Server struct {
	conn      *net.UDPConn
	buffer    []byte
	localAddr string
}

func (server *Server) Accept() (*RRQresponse, error) {

	written, addr, err := server.conn.ReadFrom(server.buffer)
	if err != nil {
		return nil, fmt.Errorf("Failed to read data from client: %v", err)
	}

	request, err := ParseRequest(server.buffer[:written])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse request: %v", err)
	}
	request.Addr = &addr
	request.LocalAddr = server.localAddr

	if request.Opcode != RRQ {
		return nil, fmt.Errorf("Unkown opcode %v", request.Opcode)
	}

	raddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve client address: %v", err)
	}

	response, err := NewRRQresponse(raddr, request, false)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewTFTPServer(addr *net.UDPAddr) (*Server, error) {
	re := regexp.MustCompile(":[^:]*$")

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed listen UDP %v", err)
	}

	return &Server{
		conn,
		make([]byte, 2048),
		re.ReplaceAllString(addr.String(), ":0"),
	}, nil

}
