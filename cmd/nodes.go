package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Show metrics for all nodes",
	Run:   showNodesMetrics,
}

func showNodesMetrics(cmd *cobra.Command, args []string) {
	config, err := getConfig()
	if err != nil {
		fmt.Println("Error to get config", err)
		os.Exit(1)
	}

	clientset, err := NewClientSet(config)
	if err != nil {
		fmt.Println("Error to create kubernetes client", err)
		os.Exit(1)
	}

	metricsClientset, err := NewMetricsClientSet(config)
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
	table.SetHeader([]string{"Node", "CPU Usage %", "Memory Usage %"})

	for _, node := range nodes.Items {
		nodeName := node.Name
		metrics, err := metricsClientset.MetricsV1beta1().NodesMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error to get metrics from node %s: %s\n", nodeName, err)
			continue
		}

		cpuUsage := metrics.Usage.CPU.MilliValue() / 10
		coloredCPU := getColorValue(cpuUsage)

		memoryUsage := metrics.Usage.Memory.Value()
		coloredMemory := getColorValue(memoryUsage)

		row := []string{nodeName, fmt.Sprintf("%d", coloredCPU), fmt.Sprintf("%d", coloredMemory)}
		table.Append(row)

	}
	table.Render()
}
