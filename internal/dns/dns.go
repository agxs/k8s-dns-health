package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Params struct {
	TestType      string
	Namespace     string
	Label         string
	Server        string
	Port          int
	Addresses     []string
	FailureAction string
	KubeConfig    string
}

type TestType interface {
	FetchDnsServers() ([]string, error)
}

type ServerTest struct {
	server string
}

func (s ServerTest) FetchDnsServers() ([]string, error) {
	return []string{s.server}, nil
}

type K8sTest struct {
	Clientset    *kubernetes.Clientset
	Namespace    string
	LabelMatcher string
}

type DnsError struct {
	Server string
	Error  error
}

func (k K8sTest) FetchDnsServers() ([]string, error) {
	pods, err := k.Clientset.CoreV1().Pods(k.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: k.LabelMatcher,
	})
	if err != nil {
		return []string{}, err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	dnsServers := []string{}
	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s, Pod IP: %s\n", pod.Name, pod.Status.PodIP)
		dnsServers = append(dnsServers, pod.Status.PodIP)
	}

	return dnsServers, nil
}

func GetTestType(params *Params) (TestType, error) {
	var test TestType

	if params.TestType == "server" {
		test = ServerTest{server: params.Server}
	} else if params.TestType == "k8s" {
		// First assume in cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Printf("In cluster k8s config unavailable, trying a .kube config file: %v\n", err)
			// if no incluster config then try and use a .kubeconfig file
			config, err = clientcmd.BuildConfigFromFlags("", params.KubeConfig)
			if err != nil {
				return nil, err
			}
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		test = K8sTest{Clientset: clientset, Namespace: params.Namespace, LabelMatcher: params.Label}
	} else {
		return nil, errors.New("Unknown test type")
	}

	return test, nil
}

func TestDns(params *Params, test TestType) (bool, string, error) {
	dnsServers, err := test.FetchDnsServers()
	if err != nil {
		return false, "", err
	}

	for _, s := range dnsServers {
		_, err := QueryDns(params.Addresses, s, params.Port)
		if err != nil {
			return false, s, err
		}
	}

	return true, "", nil
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
