# KOM - Kubernetes Object Metrics CLI tool

Kom is a command-line tool written in Go that allows you to display Kubernetes metrics for nodes and pods. It provides an easy way to monitor resource usage and quickly identify performance issues in your Kubernetes cluster.

## Features

- View CPU and memory usage metrics for all nodes in the cluster.
- View CPU and memory usage metrics for all pods in the cluster.
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

To display metrics for all nodes in the Kubernetes cluster, use the following command:

```bash
kom nodes
```

This will show a table with information about each node's CPU usage, memory usage, and IP addresses.

To display metrics for all pods in the Kubernetes cluster, use the following command:

```bash
kom pods
```

This will show a table with information about each pod's namespace, name, container, CPU usage, memory usage, and IP.
