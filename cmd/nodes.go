package cmd

import (
	"context"
	"fmt"
	"kom/kube"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var nodes *corev1.NodeList

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Show metrics for all nodes",
	Run:   showNodesMetrics,
}

func showNodesMetrics(cmd *cobra.Command, args []string) {
	config, err := kube.GetConfig()
	if err != nil {
		fmt.Println("Error to get config", err)
		os.Exit(1)
	}

	clientset, err := kube.NewClientSet(config)
	if err != nil {
		fmt.Println("Error to create kubernetes client", err)
		os.Exit(1)
	}

	metricsClientset, err := kube.NewMetricsClientSet(config)
	if err != nil {
		fmt.Println("Error to list metrics of kubernetes", err)
		os.Exit(1)
	}
	nodes, err = clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error to list nodes", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Node", "IP", "CPU Usage %", "Memory Usage %", "h", "Status - Labels"})

	for _, node := range nodes.Items {
		nodeName := node.Name
		hasLabels := len(node.Labels) > 0

		var hasNOSchedule bool
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.TaintNodeUnschedulable && condition.Status == corev1.ConditionTrue {
				hasNOSchedule = true
				break
			}
		}

		metrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error to get metrics from node %s: %s\n", nodeName, err)
			continue
		}

		cpuUsage := metrics.Usage.Cpu().MilliValue() / 10
		coloredCPU := getColorValue(int(cpuUsage))

		memoryUsage := float64(metrics.Usage.Memory().Value()) / float64(node.Status.Capacity.Memory().Value()) * 100.0
		coloredMemory := getColorValue(int(memoryUsage))

		statusEmoji := getEmoji(int(cpuUsage), int(memoryUsage))

		var nodeIPs []string
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP || address.Type == corev1.NodeExternalIP {
				nodeIPs = append(nodeIPs, address.Address)
			}
		}

		var statusMessage string
		if hasNOSchedule {
			statusMessage = "NOSchedule"
		} else {
			statusMessage = "OK"
		}

		if hasLabels {
			labels := strings.Join(getNodeLabels(node), ", ")
			statusMessage += ", Has Labels: " + labels

		}

		row := []string{nodeName, strings.Join(nodeIPs, ", "), coloredCPU, coloredMemory, statusEmoji, statusMessage}
		table.Append(row)

	}
	table.Render()
}

func getNodeLabels(node corev1.Node) []string {
	var labels []string
	for key, value := range node.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, value))
	}
	return labels
}
