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

func StartDynamicTCPProxy(hubPort string) {
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

		go handleConnection(clientConn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	clientPort := clientConn.LocalAddr().(*net.TCPAddr).Port
	port := int32(clientPort) 

	hostClient, err := client.NewHostClient()
	if err != nil {
		log.Println("Failed to create host client:", err)
		return
	}
	defer hostClient.Close()

	hostsResp, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		log.Println("Failed to get hosts from gRPC:", err)
		return
	}

	var targetHostID, targetServerID string
FOUND:
	for _, h := range hostsResp.GetHosts() {
		for serverID, serverPorts := range h.GetServers() {
			for _, p := range serverPorts.GetPorts() {
				if p == port {
					targetHostID = h.GetId()
					targetServerID = serverID
					break FOUND
				}
			}
		}
	}

	if targetHostID == "" {
		log.Printf("No server found for port %d\n", port)
		return
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		log.Println("Failed to create network client:", err)
		return
	}

	metaResp, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: targetServerID})
	if err != nil {
		log.Println("Failed to get server metadata:", err)
		return
	}

	hostAddr := metaResp.GetIp() + ":" + strconv.Itoa(int(port))
	serverConn, err := net.Dial("tcp", hostAddr)
	if err != nil {
		log.Println("Failed to connect to host server:", err)
		return
	}
	defer serverConn.Close()

	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}