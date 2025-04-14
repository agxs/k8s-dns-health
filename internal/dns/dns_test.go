package dns

import (
	"slices"
	"testing"
)

func TestServerFetchDnsServers(t *testing.T) {
	s := ServerTest{server: "aserver"}
	actual, err := s.FetchDnsServers()
	if err != nil {
		t.Errorf("FetchDnsServers should return no errors: %v", err)
	}
	if !slices.Equal(actual, []string{"aserver"}) {
		t.Errorf("Returned server should be [aserver] but was %v", actual)
	}
}

func TestTestDns(t *testing.T) {
	params := Params{
		TestType:      "server",
		Namespace:     "kube-system",
		Label:         "k8s-app=kube-dns",
		Server:        "1.1.1.1",
		Port:          53,
		Addresses:     []string{"google.com"},
		FailureAction: "teams",
		KubeConfig:    ".kube/config",
	}
	testType := ServerTest{
		server: "1.1.1.1",
	}
	success, server, err := TestDns(&params, testType)
	if !success {
		t.Errorf("DNS test success for server %s is false: %v", server, err)
	}
}
