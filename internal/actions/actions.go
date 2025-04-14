package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/agxs/k8s-dns-health/internal/dns"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FailureAction interface {
	DoAction(dnsServer string, err error) error
}

type TeamsFailureAction struct{}

type RestartFailureAction struct {
	K8sTest dns.K8sTest
}

type CombinedFailureAction struct {
	TeamsFailureAction   TeamsFailureAction
	RestartFailureAction RestartFailureAction
}

func (n TeamsFailureAction) DoAction(dnsServer string, dnsError error) error {
	message := fmt.Sprintf("Server: %s, error: %s", dnsServer, dnsError.Error())

	webhookURL := os.Getenv("TEAMS_WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("environment variable TEAMS_WEBHOOK_URL is not set")
	}

	// Prepare the payload
	payload := map[string]string{
		"text": message,
	}

	// Encode the payload as JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the appropriate header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response (200 OK or 204 No Content)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	return nil
}

func (n RestartFailureAction) DoAction(dnsServer string, err error) error {
	podList, err := n.K8sTest.Clientset.CoreV1().
		Pods(n.K8sTest.Namespace).
		List(context.Background(), metav1.ListOptions{
			FieldSelector: "status.podIP=" + dnsServer,
		})
	if err != nil {
		return err
	}

	var podNameToDelete string
	if len(podList.Items) != 1 {
		return fmt.Errorf(
			"No pod found in namespace %s with IP %s\n",
			n.K8sTest.Namespace,
			dnsServer,
		)
	} else {
		for _, pod := range podList.Items {
			podNameToDelete = pod.Name
			fmt.Printf("Found Pod %q (namespace=%s) with IP: %s\n",
				pod.Name, pod.Namespace, pod.Status.PodIP)
			break
		}
	}

	err = n.K8sTest.Clientset.CoreV1().
		Pods(n.K8sTest.Namespace).
		Delete(context.Background(), podNameToDelete, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Error deleting the pod %s: %v\n", podNameToDelete, err)
	}

	fmt.Printf("Pod %s successfully deleted.\n", podNameToDelete)

	return nil
}

func (n CombinedFailureAction) DoAction(dnsServer string, err error) error {
	new_err := n.TeamsFailureAction.DoAction(dnsServer, err)
	if new_err != nil {
		return new_err
	}

	new_err = n.RestartFailureAction.DoAction(dnsServer, err)
	if new_err != nil {
		return new_err
	}
	return nil
}
