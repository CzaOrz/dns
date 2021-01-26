package dns

type IStore interface {
	GetAddressOfA(host string) (AddressOfA, error)
}

type AddressOfA struct {
	Host   string    `json:"host"`
	IPPool [][4]byte `json:"ip_pool"`
}
