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
	table.SetHeader([]string{"Node", "IP", "CPU Usage %", "Memory Usage %", "h"})

	for _, node := range nodes.Items {
		nodeName := node.Name
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

		row := []string{nodeName, strings.Join(nodeIPs, ", "), coloredCPU, coloredMemory, statusEmoji}
		table.Append(row)

	}
	table.Render()
}
