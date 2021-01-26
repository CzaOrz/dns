package main

import (
	"encoding/json"
	"github.com/czaorz/dns"
	"io/ioutil"
	"net/http"
)

var store dns.IStore

type HttpHandler struct{}

func (h HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.RequestURI {
	case "/add":
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		data := struct {
			Host string   `json:"host"`
			IPS  []string `json:"ips"`
		}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		err = store.SetAddressOfA(data.Host, data.IPS...)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write([]byte("ok"))
	}
}

func DNSServe(store dns.IStore) {
	server := dns.DNS{Store: store}
	server.ServeAndStart()
}

func main() {
	store = dns.NewMemoryStore()
	go DNSServe(store)
	err := http.ListenAndServe(":8080", HttpHandler{})
	if err != nil {
		panic(err)
	}
}
