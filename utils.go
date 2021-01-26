package dns

import (
	"errors"
	"strconv"
	"strings"
)

func IP2A(ip string) ([4]byte, error) {
	ips := strings.Split(ip, ".")
	if len(ips) != 4 {
		return [4]byte{}, errors.New("length is invalid")
	}
	ip1, err := strconv.ParseInt(ips[0], 10, 0)
	if err != nil {
		return [4]byte{}, errors.New("ip is invalid")
	}
	ip2, err := strconv.ParseInt(ips[1], 10, 0)
	if err != nil {
		return [4]byte{}, errors.New("ip is invalid")
	}
	ip3, err := strconv.ParseInt(ips[2], 10, 0)
	if err != nil {
		return [4]byte{}, errors.New("ip is invalid")
	}
	ip4, err := strconv.ParseInt(ips[3], 10, 0)
	if err != nil {
		return [4]byte{}, errors.New("ip is invalid")
	}
	return [4]byte{byte(ip1), byte(ip2), byte(ip3), byte(ip4)}, nil
}
