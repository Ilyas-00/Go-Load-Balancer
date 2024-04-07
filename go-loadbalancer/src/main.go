package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)


type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
	alive bool
	mu    sync.RWMutex 
}

func NewSimpleServer(addr string) *SimpleServer {
	fmt.Printf("URL en cours d'analyse : %s\n", addr)
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("error parsing server URL: %v\n", err)
		os.Exit(1)
	}


	return &SimpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
		alive: true, 
	}
}

func (s *SimpleServer) Address() string {
	return s.addr
}

func (s *SimpleServer) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.alive
}

func (s *SimpleServer) SetAlive(alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alive = alive
}

func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

type LoadBalancer struct {
	port           string
	roundRobinCount int
	servers        []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:           port,
		roundRobinCount: 0,
		servers:        servers,
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	lb.roundRobinCount++
	totalServers := len(lb.servers)
	for i := 0; i < totalServers; i++ {
		server := lb.servers[lb.roundRobinCount%totalServers]
		if server.IsAlive() {
			return server
		}
		lb.roundRobinCount++
	}
	return nil 
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()
	if targetServer == nil {
		http.Error(rw, "No available servers", http.StatusServiceUnavailable)
		return
	}
	fmt.Printf("Forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, req)
}

func main() {
	servers := []Server{
		NewSimpleServer("http://www.facebook.com"),
		NewSimpleServer("http://www.bing.com"),
		NewSimpleServer("http://www.duckduckgo.com"),
	}

	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.ServeProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Serving requests at localhost:%s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
