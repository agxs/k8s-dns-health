package actions

import (
	"context"
	"testing"

	"github.com/agxs/k8s-dns-health/internal/dns"
	testclient "k8s.io/client-go/kubernetes/fake"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDoAction(t *testing.T) {
	testNamespace := "kube-system"
	client := testclient.NewClientset()
	k8sTest := dns.K8sTest{
		Client:       client,
		Namespace:    testNamespace,
		LabelMatcher: "k8s-app=kube-dns",
	}
	action := RestartFailureAction{
		K8sTest: k8sTest,
	}

	// Create a mock pod
	_, err := client.CoreV1().Pods(testNamespace).Create(
		context.Background(),
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod-1",
				Namespace: testNamespace,
				Labels: map[string]string{
					"k8s-app": "kube-dns",
				},
			},
			Status: corev1.PodStatus{
				PodIP: "10.0.0.1",
			},
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		t.Errorf("Somehow got an error mocking a pod: %v", err)
	}

	err = action.DoAction("10.0.0.1", nil)
	if err != nil {
		t.Errorf("Error running Restart DoAction: %v", err)
	}
}

func TestDoActionNoPods(t *testing.T) {
	testNamespace := "kube-system"
	client := testclient.NewClientset()
	k8sTest := dns.K8sTest{
		Client:       client,
		Namespace:    testNamespace,
		LabelMatcher: "k8s-app=kube-dns",
	}
	action := RestartFailureAction{
		K8sTest: k8sTest,
	}
	err := action.DoAction("10.0.0.1", nil)
	if err == nil {
		t.Errorf("Should generate an error")
	}
}
