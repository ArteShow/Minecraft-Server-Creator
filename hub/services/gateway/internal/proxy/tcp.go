package proxy

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/client"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/proto/network-service"
)

type LookupFunc func(port int) (hostAddr string, ok bool)

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
			hostAddr, ok := lookup(clientPort)
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

func LookupHostByPort(port int) (string, bool) {
	hostClient, err := client.NewHostClient()
	if err != nil {
		log.Println("Failed to create host client:", err)
		return "", false
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		log.Println("Failed to create network client:", err)
		return "", false
	}

	hostsResp, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		log.Println("Failed to get hosts:", err)
		return "", false
	}

	for _, h := range hostsResp.GetHosts() {
		for serverID, ports := range h.GetServers() {
			for _, p := range ports.GetPorts() {
				if int(p) == port {
					meta, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: serverID})
					if err != nil {
						log.Println("Failed to get server metadata:", err)
						return "", false
					}
					return meta.GetIp() + ":" + strconv.Itoa(port), true
				}
			}
		}
	}

	return "", false
}