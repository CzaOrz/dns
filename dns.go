package dns

import (
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

type DNS struct {
	Store   IStore
	UDPConn *net.UDPConn
}

func (dns *DNS) ServeAndStart() {
	if dns.Store == nil {
		dns.Store = NewMemoryStore()
	}
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()
	dns.UDPConn = udpConn
	for {
		buf := make([]byte, 512)
		_, addr, _ := udpConn.ReadFromUDP(buf)

		var msg dnsmessage.Message
		err := msg.Unpack(buf)
		if err != nil {
			continue
		}
		if len(msg.Questions) < 1 {
			continue
		}
		go dns.QueryAndFillDNS(addr, msg)
	}
}

func (dns *DNS) QueryAndFillDNS(addr *net.UDPAddr, msg dnsmessage.Message) {
	for _, question := range msg.Questions {
		switch question.Type {
		case dnsmessage.TypeA:
			addressOfA, err := dns.Store.GetAddressOfA(question.Name.String())
			if err != nil {
				continue
			}
			for _, ip := range addressOfA.IPPool {
				a, err := IP2A(ip)
				if err != nil {
					continue
				}
				resource := dnsmessage.Resource{
					Header: dnsmessage.ResourceHeader{
						Name:  question.Name,
						Class: dnsmessage.ClassINET,
						TTL:   600,
					},
					Body: &dnsmessage.AResource{
						A: a,
					},
				}
				msg.Answers = append(msg.Answers, resource)
			}
		default:
			continue
		}
	}
	msg.Response = true
	packed, err := msg.Pack()
	if err != nil {
		return
	}
	_, _ = dns.UDPConn.WriteToUDP(packed, addr)
}
