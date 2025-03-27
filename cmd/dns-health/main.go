package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"os"
	"strings"

	"github.com/agxs/k8s-dns-health/internal/dns"

	"k8s.io/client-go/util/homedir"
)

func main() {
	// read startup/env variables
	// - read dns server mode
	//   - read k8s namespace, label
	//   - read dns server
	// - read test addresses (comma separated)
	// - read failure action type (teams or restart or both)
	// - read test frequency

	// choose direct dns or kube label listing
	// query dns items, 1 or more
	// send teams message if failure
	// future, restart pod if failure

	testType := flag.String("testType", "server", "The test type to use, either 'server' or 'k8s'")
	namespace := flag.String("namespace", "kube-system", "The Kubernetes namespace to scan")
	label := flag.String("label", "k8s-app=kube-dns", "The pod label to scan for")
	server := flag.String("server", "192.168.0.1", "The server IP to query")
	port := flag.Int("port", 53, "The server port to query")
	addresses := flag.String("addresses", "kubernetes.default,google.com", "The DNS test queries")
	failureType := flag.String("failureType", "teams", "Options for failure actions, 'teams', 'restart', 'both'")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	if *testType != "server" && *testType != "k8s" {
		fmt.Printf("Unknown testType: %s\n", *testType)
		os.Exit(1)
	}

	addressesSplit := strings.Split(*addresses, ",")

	params := dns.Params{
		TestType:    *testType,
		Namespace:   *namespace,
		Label:       *label,
		Server:      *server,
		Port:        *port,
		Addresses:   addressesSplit,
		FailureType: *failureType,
		KubeConfig:  *kubeconfig,
	}

	test, err := dns.GetTestType(&params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	_, err = dns.TestDns(&params, test)
	if err != nil {
		//todo handle failure type
	}
}
