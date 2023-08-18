[!["Buy Me A Coffee"](https://img.shields.io/badge/Buy_Me_A_Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/miltlima)

# KOM - Kubernetes Object Metrics CLI tool

Kom is a command-line tool written in Go that allows you to display Kubernetes metrics for nodes and pods. It provides an easy way to monitor resource usage and quickly identify performance issues in your Kubernetes cluster.

## Features

- View CPU and memory usage metrics for all nodes in the cluster.
- View CPU and memory usage metrics for all pods in the cluster.
- View Logs from pods and save inside folder komlogs
- Color-coded output for easy visualization of resource usage levels.

### Color-Coded Output

Kom uses color-coding to visualize resource usage levels:

- CPU usage above 80% is displayed in red.
- CPU usage between 50% and 80% is displayed in yellow.
- CPU usage below 50% is displayed in green.
- Memory usage above 80% is displayed in red.
- Memory usage between 50% and 80% is displayed in yellow.
- Memory usage below 50% is displayed in green.

## Installation

Assuming you have Git and Go installed on your system, here's what you can do:

Clone the repository:

```bash
git clone https://github.com/miltlima/kom.git
```

Change to the project directory:

```bash
cd kom
```

Build the binary:

```bash
go build
```

This will generate a binary file named "kom" in the project directory.

```bash
sudo mv kom /usr/local/bin
```

Now, the "kom" binary is located in your bin folder, and you can run it from any terminal window as long as the bin folder is in your system's PATH.

Save the file and either restart your terminal or run source to apply the changes:

```bash
source ~/.bashrc
```

Or

```bash
source ~/.bash_profile
```

Now, you should be able to run the "kom" command from anywhere in your terminal.

## Usage

## Kom nodes

To display metrics for all nodes in the Kubernetes cluster, use the following command:

```bash
kom nodes
```

This will show a table with information about each node's CPU usage, memory usage, and IP addresses.

```bash
❯ ./kom nodes
+-------+-----------+-------------+----------------+----+
| NODE  |    IP     | CPU USAGE % | MEMORY USAGE % | H  |
+-------+-----------+-------------+----------------+----+
| node1 | 10.0.0.11 | 5           | 32             | 🟩 |
| node2 | 10.0.0.12 | 2           | 25             | 🟩 |
| node3 | 10.0.0.13 | 0           | 27             | 🟩 |
| node4 | 10.0.0.14 | 1           | 27             | 🟩 |
+-------+-----------+-------------+----------------+----+
```

## Kom pods

To display metrics for all pods in the Kubernetes cluster, use the following command:

```bash
kom pods
```

This will show a table with information about each node, pod's namespace, name, CPU usage, memory usage, and IP.

```bash
❯ ./kom pods
+-------+-------------+--------------------------------------------------+-----------+-------------+----------------+----+
| NODE  |  NAMESPACE  |                       POD                        |  POD IP   | CPU USAGE % | MEMORY USAGE % | H  |
+-------+-------------+--------------------------------------------------+-----------+-------------+----------------+----+
| node2 | default     | build-code-deployment-68dd47875-4bt5p            | 10.36.0.1 | 0           | 0              | 🟩 |
| node2 | default     | health-check-deployment-59f4b679b-8k4pj          | 10.36.0.5 | 0           | 0              | 🟩 |
| node3 | default     | hidden-in-layers-qtst5                           | 10.44.0.1 | 0           | 0              | 🟩 |
| node2 | default     | internal-proxy-deployment-7699c5dd65-xm4tw       | 10.36.0.3 | 0           | 0              | 🟩 |
| node2 | default     | kubernetes-goat-home-deployment-7bf7785ff5-gghts | 10.36.0.2 | 0           | 0              | 🟩 |
| node4 | default     | metadata-db-648b64948f-vjvsg                     | 10.42.0.1 | 0           | 0              | 🟩 |
| node2 | default     | poor-registry-deployment-75f47d55dc-vhs9d        | 10.36.0.4 | 0           | 0              | 🟩 |
| node4 | default     | system-monitor-deployment-674bb4dc65-9wj4m       | 10.42.0.3 | 0           | 0              | 🟩 |
| node1 | kube-system | coredns-787d4945fb-4vnqj                         | 10.32.0.3 | 0           | 0              | 🟩 |
| node1 | kube-system | coredns-787d4945fb-q5r2h                         | 10.32.0.2 | 0           | 0              | 🟩 |
| node1 | kube-system | etcd-node1                                       | 10.0.0.11 | 1           | 1              | 🟩 |
| node1 | kube-system | kube-apiserver-node1                             | 10.0.0.11 | 3           | 8              | 🟩 |
| node1 | kube-system | kube-controller-manager-node1                    | 10.0.0.11 | 0           | 0              | 🟩 |
| node4 | kube-system | kube-proxy-5r278                                 | 10.0.0.14 | 0           | 0              | 🟩 |
| node3 | kube-system | kube-proxy-dzzrp                                 | 10.0.0.13 | 0           | 0              | 🟩 |
| node1 | kube-system | kube-proxy-h5wsb                                 | 10.0.0.11 | 0           | 0              | 🟩 |
| node2 | kube-system | kube-proxy-htlwv                                 | 10.0.0.12 | 0           | 0              | 🟩 |
| node1 | kube-system | kube-scheduler-node1                             | 10.0.0.11 | 0           | 0              | 🟩 |
| node4 | kube-system | metrics-server-75fcb88b7d-n2l7p                  | 10.0.0.14 | 0           | 0              | 🟩 |
| node4 | kube-system | weave-net-774l4                                  | 10.0.0.14 | 0           | 0              | 🟩 |
| node3 | kube-system | weave-net-gv6dn                                  | 10.0.0.13 | 0           | 0              | 🟩 |
| node2 | kube-system | weave-net-n6n8f                                  | 10.0.0.12 | 0           | 0              | 🟩 |
| node1 | kube-system | weave-net-sn6v9                                  | 10.0.0.11 | 0           | 0              | 🟩 |
+-------+-------------+--------------------------------------------------+-----------+-------------+----------------+----+
```

## Kom logs

The logs command in the kom CLI allows you to collect logs from a Kubernetes pod and optionally save them to a file.

```bash
kom logs <pod-name> [flags]
```

### Arguments

`<pod-name>`: The name of the pod from which you want to collect logs.

### Flags

```bash
-s, --save: (Optional) Save the output to a log file. If this flag is not provided, the logs will be displayed in the terminal.
-o, --output <file-path>: (Optional) Specify the file path to save the logs. Default is output.log in the komlogs folder.
-n, --namespace <namespace>: (Optional) Specify the namespace of the pod. Default is default.
-c, --container <container-name>: (Optional) Specify the name of the container in the pod. If not provided, logs will be collected from the first container in the pod.
```
