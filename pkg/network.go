package pkg

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func RunServer(address string) error {
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

	for {
		// Receive timestamp from client
		var clientTime int64
		err := binary.Read(conn, binary.BigEndian, &clientTime)
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed connection")
				return
			}
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
}

func connectWithRetry(address string, maxRetries int, retryInterval time.Duration) (net.Conn, error) {
	for attempt := range maxRetries {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			return conn, nil
		}
		if attempt < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}
	return nil, fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

func isConnectionError(err error) bool {
	if err == io.EOF {
		return true
	}
	_, ok := err.(net.Error)
	return ok
}

func RunClient(address string, count int, interval int, maxRetries int, retryInterval time.Duration) error {
	// Load env
	ddAgent := os.Getenv("DD_AGENT")
	ddEnv := os.Getenv("DD_ENV")
	ddService := os.Getenv("DD_Service")
	InitStatsD(ddAgent, ddEnv, ddService)

	var err error
	conn, err := connectWithRetry(address, maxRetries, retryInterval)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; count == 0 || i < count; i++ {
		clientTime := time.Now().UnixNano()

		err := binary.Write(conn, binary.BigEndian, &clientTime)
		if err != nil {
			if isConnectionError(err) {
				conn, err = connectWithRetry(address, maxRetries, retryInterval)
				if err != nil {
					return err
				}
				i--
				continue
			}
			return err
		}

		var serverTime int64
		err = binary.Read(conn, binary.BigEndian, &serverTime)
		if err != nil {
			if isConnectionError(err) {
				conn, err = connectWithRetry(address, maxRetries, retryInterval)
				if err != nil {
					return err
				}
				i--
				continue
			}
			return fmt.Errorf("error reading server response: %v", err)
		}

		// Calculate round-trip time and offset
		currentTime := time.Now().UnixNano()
		rtt := currentTime - clientTime
		offset := (serverTime - clientTime - rtt/2)

		rtt_f := float64(rtt) / 1e6
		offset_f := float64(offset) / 1e6

		fmt.Printf("Round-trip time: %.3f ms, Offset: %.3f ms\n",
			rtt_f, offset_f)

		SendLatency(rtt_f, offset_f)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	return nil
}
