package proxy

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/client"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/proto/network-service"
)

var (
	portMap   = map[int]string{}
	listeners = map[int]net.Listener{}
	mu        sync.RWMutex
)

func StartProxy(_ []int) {
	go func() {
		for {
			syncAndUpdate()
			time.Sleep(3 * time.Second)
		}
	}()
}

func syncAndUpdate() {
	hostClient, err := client.NewHostClient()
	if err != nil {
		return
	}
	defer hostClient.Close()

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return
	}
	defer networkClient.Close()

	hosts, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return
	}

	newMap := map[int]string{}

	for _, h := range hosts.GetHosts() {
		ipResp, err := networkClient.GetServerMetadata(
			&network.GetServerMetadataRequest{ServerId: h.GetId()},
		)
		if err != nil {
			continue
		}

		ip := ipResp.GetIp()

		// ✅ FIX HERE
		for _, port := range h.GetServers() {
			p := int(port)
			newMap[p] = ip + ":" + strconv.Itoa(p)
		}
	}

	mu.Lock()
	portMap = newMap
	updateListenersLocked()
	mu.Unlock()
}

func updateListenersLocked() {
	for port := range portMap {
		if _, exists := listeners[port]; !exists {
			go startListener(port)
		}
	}

	for port, ln := range listeners {
		if _, exists := portMap[port]; !exists {
			ln.Close()
			delete(listeners, port)
			log.Println("closed listener on port", port)
		}
	}
}

func startListener(port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	mu.Lock()
	listeners[port] = ln
	mu.Unlock()

	log.Println("listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		go handleConn(conn, port)
	}
}

func handleConn(c net.Conn, port int) {
	defer c.Close()

	mu.RLock()
	target, ok := portMap[port]
	mu.RUnlock()

	if !ok {
		log.Println("no backend for port", port)
		return
	}

	serverConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer serverConn.Close()

	go io.Copy(serverConn, c)
	io.Copy(c, serverConn)
}