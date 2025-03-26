package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

func runServer(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleServerConnection(conn)
	}
}

func handleServerConnection(conn net.Conn) {
	defer conn.Close()

	// Receive timestamp from client
	var clientTime int64
	err := binary.Read(conn, binary.BigEndian, &clientTime)
	if err != nil {
		log.Printf("Error reading client time: %v", err)
		return
	}

	// Get current server time
	serverTime := time.Now().UnixNano()

	// Send server time back to client
	err = binary.Write(conn, binary.BigEndian, &serverTime)
	if err != nil {
		log.Printf("Error sending server time: %v", err)
		return
	}
}

func runClient(address string, count int, interval int) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < count; i++ {
		// Get current time
		clientTime := time.Now().UnixNano()

		// Send timestamp to server
		err := binary.Write(conn, binary.BigEndian, &clientTime)
		if err != nil {
			return err
		}

		// Receive server time
		var serverTime int64
		err = binary.Read(conn, binary.BigEndian, &serverTime)
		if err != nil {
			return err
		}

		// Calculate round-trip time and offset
		currentTime := time.Now().UnixNano()
		rtt := currentTime - clientTime
		offset := (serverTime - clientTime - rtt/2)

		fmt.Printf("Round-trip time: %.3f ms, Offset: %.3f ms\n",
			float64(rtt)/1e6, float64(offset)/1e6)

		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	return nil
}
