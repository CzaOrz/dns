package dns

type IStore interface {
	GetAddressOfA(host string) (AddressOfA, error)
	SetAddressOfA(host string, ips ...string) error
}

type AddressOfA struct {
	Host   string   `json:"host"`
	IPPool []string `json:"ip_pool"`
}
