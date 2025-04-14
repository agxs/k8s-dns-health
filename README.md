# DNS Health Check for Kubernetes

This project is a Go program designed to run within a Kubernetes cluster to
perform health checks on DNS pods, ensuring their functionality. If a DNS pod
is not working, the program can either send a Microsoft Teams message, restart
the pod, or perform both actions. The program also supports running outside of
a Kubernetes cluster for testing purposes.

## Features

- Health checks on DNS pods within a Kubernetes cluster.
- Failure handling with configurable actions: Teams notification, pod restart, or both.
- Run in test mode outside of a Kubernetes cluster.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/agxs/k8s-dns-health.git
   ```
2. Navigate to the project directory:
   ```sh
   cd k8s-dns-health
   ```
3. Build the program:
   ```sh
   make
   ```

By default this will build binaries for both Linux/x64 and MacOS/Arm64.

## Usage

```sh
./dns-health-check [flags]
```

### Command Line Parameters

- `-addresses string`: The DNS test queries (default: `"kubernetes.default,google.com"`).
- `-failureAction string`: Options for failure actions: `teams`, `restart`, `both` (default: `teams`).
- `-kubeconfig string`: Optional absolute path to the kubeconfig file (default: `"$HOME/.kube/config"`).
- `-label string`: The pod label to scan for (default: `"k8s-app=kube-dns"`).
- `-namespace string`: The Kubernetes namespace to scan (default: `"kube-system"`).
- `-port int`: The server port to query (default: `53`).
- `-server string`: The server IP to query (default: `"192.168.0.1"`).
- `-testType string`: The test type to use: `server` or `k8s` (default: `server`).

## Example

To run the program with specific configurations inside a cluster:

```sh
./dns-health-check -addresses "kubernetes.default,example.com" -failureAction "both" -namespace "custom-namespace" -testType "k8s"
```

To run in test mode using a local DNS server:

```sh
./dns-health-check -addresses "google.com"
```

### Kubernetes Configuration

The `setup-k8s.sh` script will create a service account with the correct
permissions for the `kube-system` namespace.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
