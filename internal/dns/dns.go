package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Params struct {
	TestType    string
	Namespace   string
	Label       string
	Server      string
	Port        int
	Addresses   []string
	FailureType string
}

type TestType interface {
	FetchDnsServers() []string
}

type FailureType interface {
	OnFailure()
}

type ServerTest struct {
	server string
}

func (s ServerTest) FetchDnsServers() []string {
	return []string{s.server}
}

func TestDns(params *Params) (bool, error) {
	dnsServers := []string{}
	if params.TestType == "server" {
		test := ServerTest{server: params.Server}
		dnsServers = test.FetchDnsServers()
	} else {
		return false, errors.New("Unknown test type")
	}

	for _, s := range dnsServers {
		_, err := QueryDns(params.Addresses, s, params.Port)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func QueryDns(addresses []string, server string, port int) (bool, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, server+":"+strconv.Itoa(port))
		},
	}
	for _, address := range addresses {
		ip, err := resolver.LookupHost(context.Background(), address)
		if err != nil {
			fmt.Printf("Error in dns lookup: %v\n", err)
			return false, err
		} else {
			fmt.Printf("'%s' resolves to %v\n", address, ip)
		}
	}
	return true, nil
}
