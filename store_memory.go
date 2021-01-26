package dns

import (
	"fmt"
	"strings"
)

type MemoryStore struct {
	AddressOfAMap map[string]AddressOfA
}

func NewMemoryStore() *MemoryStore {
	store := MemoryStore{}
	store.FillByDefault()
	return &store
}

func (ms *MemoryStore) FillByDefault() {
	if ms.AddressOfAMap == nil {
		ms.AddressOfAMap = map[string]AddressOfA{}
	}
}

func (ms *MemoryStore) GetAddressOfA(host string) (AddressOfA, error) {
	if !strings.HasSuffix(host, ".") {
		host = fmt.Sprintf("%s.", host)
	}
	addressOfA, ok := ms.AddressOfAMap[host]
	if !ok {
		return addressOfA, fmt.Errorf("host[%s] not found.", host)
	}
	return addressOfA, nil
}

func (ms *MemoryStore) SetAddressOfA(host string, ips ...string) error {
	if !strings.HasSuffix(host, ".") {
		host = fmt.Sprintf("%s.", host)
	}
	addressOfA, ok := ms.AddressOfAMap[host]
	if !ok {
		addressOfA = AddressOfA{Host: host, IPPool: []string{}}
	}
	addressOfA.IPPool = append(addressOfA.IPPool, ips...)
	ms.AddressOfAMap[host] = addressOfA
	return nil
}
