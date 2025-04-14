package main

import (
	"testing"

	"github.com/agxs/k8s-dns-health/internal/actions"
	"github.com/agxs/k8s-dns-health/internal/dns"
)

func TestGetFailureActionUnknownAction(t *testing.T) {
	testParams := dns.Params{
		TestType:      "unused",
		Namespace:     "ns",
		Label:         "l",
		Server:        "srv",
		Port:          53,
		Addresses:     []string{"a", "b"},
		FailureAction: "unknown",
		KubeConfig:    "config",
	}
	_, err := getFailureAction(testParams, dns.ServerTest{})
	if err == nil {
		t.Error("Should generate error when unknown failure action is specified")
	}
}

func TestGetFailureActionTeams(t *testing.T) {
	testParams := dns.Params{
		TestType:      "unused",
		Namespace:     "ns",
		Label:         "l",
		Server:        "srv",
		Port:          53,
		Addresses:     []string{"a", "b"},
		FailureAction: "teams",
		KubeConfig:    "config",
	}
	action, err := getFailureAction(testParams, dns.ServerTest{})
	if err != nil {
		t.Errorf("No error should be generated: %v", err.Error())
	}

	_, ok := action.(actions.TeamsFailureAction)
	if !ok {
		t.Error("Action should be a Teams action")
	}
}

func TestGetFailureActionRestart(t *testing.T) {
	testParams := dns.Params{
		TestType:      "unused",
		Namespace:     "ns",
		Label:         "l",
		Server:        "srv",
		Port:          53,
		Addresses:     []string{"a", "b"},
		FailureAction: "restart",
		KubeConfig:    "config",
	}
	action, err := getFailureAction(testParams, dns.K8sTest{})
	if err != nil {
		t.Errorf("No error should be generated: %v", err.Error())
	}

	_, ok := action.(actions.RestartFailureAction)
	if !ok {
		t.Error("Action should be a Restart action")
	}
}

func TestGetFailureActionRestartNonK8s(t *testing.T) {
	testParams := dns.Params{
		TestType:      "unused",
		Namespace:     "ns",
		Label:         "l",
		Server:        "srv",
		Port:          53,
		Addresses:     []string{"a", "b"},
		FailureAction: "restart",
		KubeConfig:    "config",
	}
	_, err := getFailureAction(testParams, dns.ServerTest{})
	if err == nil {
		t.Error("Restart actions should be requiring a K8s test type")
	}
}

func TestGetFailureActionBoth(t *testing.T) {
	testParams := dns.Params{
		TestType:      "unused",
		Namespace:     "ns",
		Label:         "l",
		Server:        "srv",
		Port:          53,
		Addresses:     []string{"a", "b"},
		FailureAction: "both",
		KubeConfig:    "config",
	}
	action, err := getFailureAction(testParams, dns.K8sTest{})
	if err != nil {
		t.Errorf("No error should be generated: %v", err.Error())
	}

	_, ok := action.(actions.CombinedFailureAction)
	if !ok {
		t.Error("Action should be a Combined action")
	}
}

func TestGetFailureActionBothNonK8s(t *testing.T) {
	testParams := dns.Params{
		TestType:      "unused",
		Namespace:     "ns",
		Label:         "l",
		Server:        "srv",
		Port:          53,
		Addresses:     []string{"a", "b"},
		FailureAction: "both",
		KubeConfig:    "config",
	}
	_, err := getFailureAction(testParams, dns.ServerTest{})
	if err == nil {
		t.Error("Both actions should be requiring a K8s test type")
	}
}
