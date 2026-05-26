package ratelimiter

import (
	"log"
	"net"
	"sync"
)

const NbRateLimitAuthorizedByIp = 5

type IPRateLimiter struct {
	clients map[string]int
	mutex   sync.Mutex
}

func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		clients: make(map[string]int),
	}
}

func (ipRateLimiter *IPRateLimiter) Allow(ip string) bool {
	ipRateLimiter.mutex.Lock()
	defer ipRateLimiter.mutex.Unlock()

	ip, _, err := net.SplitHostPort(ip)
	if err != nil {
		log.Println(err)
		return false
	}

	if ipRateLimiter.clients[ip] >= NbRateLimitAuthorizedByIp {
		return false
	}

	ipRateLimiter.clients[ip]++
	log.Printf("[%s] has now %d active connections", ip, ipRateLimiter.clients[ip])
	return true
}

func (ipRateLimiter *IPRateLimiter) Release(ip string) {
	ipRateLimiter.mutex.Lock()
	defer ipRateLimiter.mutex.Unlock()

	ip, _, _ = net.SplitHostPort(ip)

	if _, ok := ipRateLimiter.clients[ip]; ok {
		ipRateLimiter.clients[ip]--
		log.Printf("[%s] has been released : %d connections", ip, ipRateLimiter.clients[ip])
	}
}
