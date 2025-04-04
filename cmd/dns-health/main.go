package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agxs/k8s-dns-health/internal/dns"
	"github.com/agxs/k8s-dns-health/internal/notifications"

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
	notificationType := flag.String("notificationType", "teams", "Options for notification actions, 'teams', 'restart', 'both'")

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
		TestType:         *testType,
		Namespace:        *namespace,
		Label:            *label,
		Server:           *server,
		Port:             *port,
		Addresses:        addressesSplit,
		NotificationType: *notificationType,
		KubeConfig:       *kubeconfig,
	}

	test, err := dns.GetTestType(&params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var errorServer string
	_, errorServer, err = dns.TestDns(&params, test)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		notification, err := getNotificationType(params, test)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		notification.Notify(errorServer, err)
	}
}

func isK8sTest(testType dns.TestType) bool {
	_, ok := testType.(dns.K8sTest)
	return ok
}

func getNotificationType(params dns.Params, testType dns.TestType) (notifications.NotificationType, error) {
	if params.NotificationType == "teams" {
		return notifications.TeamsNotification{}, nil
	} else if params.NotificationType == "restart" {
		if !isK8sTest(testType) {
			return nil, fmt.Errorf("restart notification type only works with k8s")
		}

		return notifications.RestartNotification{K8sTest: testType.(dns.K8sTest)}, nil
	} else if params.NotificationType == "both" {
		if !isK8sTest(testType) {
			return nil, fmt.Errorf("both notification type only works with k8s")
		}

		teams := notifications.TeamsNotification{}
		restart := notifications.RestartNotification{K8sTest: testType.(dns.K8sTest)}
		return notifications.CombinedNotification{TeamsNotification: teams, RestartNotification: restart}, nil
	}

	return nil, fmt.Errorf("Unknown notification type: %s", params.NotificationType)
}
