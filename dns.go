package dns

import (
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

type DNS struct {
	Store   IStore
	UDPConn *net.UDPConn
}

func (dns *DNS) ServeAndStart() {
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
				return
			}
			for _, ip := range addressOfA.IPPool {
				resource := dnsmessage.Resource{
					Header: dnsmessage.ResourceHeader{
						Name:  question.Name,
						Class: dnsmessage.ClassINET,
						TTL:   600,
					},
					Body: &dnsmessage.AResource{
						A: ip,
					},
				}
				msg.Answers = append(msg.Answers, resource)
			}
		default:
			continue
		}
	}

	question := msg.Questions[0]
	var (
		queryNameStr = question.Name.String()
		queryName, _ = dnsmessage.NewName(queryNameStr)
	)
	var resource dnsmessage.Resource
	switch question.Type {
	case dnsmessage.TypeA:
		if rst, ok := addressBookOfA[queryNameStr]; ok {
			resource = NewAResource(queryName, rst)
		}
	default:
		return
	}
	msg.Response = true
	msg.Answers = append(msg.Answers, resource)
	packed, err := msg.Pack()
	if err != nil {
		return
	}
	_, _ = dns.UDPConn.WriteToUDP(packed, addr)
}

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Listing ...")
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)

		var msg dnsmessage.Message
		if err := msg.Unpack(buf); err != nil {
			fmt.Println(err)
			continue
		}
		go ServerDNS(addr, conn, msg)
	}
}

// address books
var (
	addressBookOfA = map[string][4]byte{
		"www.baidu.com.": [4]byte{220, 181, 38, 150},
	}
)

// ServerDNS serve
func ServerDNS(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	// query info
	if len(msg.Questions) < 1 {
		return
	}
	question := msg.Questions[0]
	var (
		queryTypeStr = question.Type.String()
		queryNameStr = question.Name.String()
		queryType    = question.Type
		queryName, _ = dnsmessage.NewName(queryNameStr)
	)
	fmt.Printf("[%s] queryName: [%s]\n", queryTypeStr, queryNameStr)

	// find record
	var resource dnsmessage.Resource
	switch queryType {
	case dnsmessage.TypeA:
		if rst, ok := addressBookOfA[queryNameStr]; ok {
			resource = NewAResource(queryName, rst)
		} else {
			fmt.Printf("not fount A record queryName: [%s] \n", queryNameStr)
			Response(addr, conn, msg)
			return
		}
	default:
		fmt.Printf("not support dns queryType: [%s] \n", queryTypeStr)
		return
	}

	// send response
	msg.Response = true
	msg.Answers = append(msg.Answers, resource)
	Response(addr, conn, msg)
}

// Response return
func Response(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	packed, err := msg.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := conn.WriteToUDP(packed, addr); err != nil {
		fmt.Println(err)
	}
}

// NewAResource A record
func NewAResource(query dnsmessage.Name, a [4]byte) dnsmessage.Resource {
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  query,
			Class: dnsmessage.ClassINET,
			TTL:   600,
		},
		Body: &dnsmessage.AResource{
			A: a,
		},
	}
}
