package proxy

import (
	"io"
	"log"
	"net"
)

type LookupFunc func(port string) (hostAddr string, ok bool)

func StartDynamicTCPProxy(hubPort string, lookup LookupFunc) {
    ln, err := net.Listen("tcp", ":"+hubPort)
    if err != nil {
        log.Fatal("Failed to listen on hub port:", err)
    }
    log.Printf("Dynamic TCP proxy listening on hub port %s\n", hubPort)

    for {
        clientConn, err := ln.Accept()
        if err != nil {
            log.Println("Failed to accept connection:", err)
            continue
        }

        go func(c net.Conn) {
            defer c.Close()

            clientPort := c.LocalAddr().(*net.TCPAddr).Port
            hostAddr, ok := lookupInt(clientPort)
            if !ok {
                log.Printf("No host mapping for port %d\n", clientPort)
                return
            }

            serverConn, err := net.Dial("tcp", hostAddr)
            if err != nil {
                log.Println("Failed to connect to host:", err)
                return
            }
            defer serverConn.Close()

            go io.Copy(serverConn, c)
            io.Copy(c, serverConn)
        }(clientConn)
    }
}

func lookupInt(port int) (string, bool) {
    if port == 25566 {
        return "host1:25566", true
    }
    if port == 25567 {
        return "host2:25567", true
    }
    return "", false
}