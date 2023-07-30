package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var PodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Show metrics for all pods",
	Run:   showPodMetrics,
}

func showPodMetrics(cmd *cobra.Command, args []string) {
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

	pods, err = clientset.CoreV1().Pods().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error to list pods", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Container", "CPU Usage %", "Memory Usage %"})

	for _, pod := range pods.Items {
		podName := pod.Name
		podNamespace := pod.Namespace
		metrics, err := metricsClientset.MetricsV1beta1().PodMetricses().Get(podNamespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error to get metrics from pod %s: %s\n", podNamespace, podName, err)
			continue
		}

		for _, container := range metrics.Containers {
			cpuUsage := container.Usage.Cpu.MilliValue() / 10
			coloredCPU := getColorValue(cpuUsage)

			memoryUsage := container.Usage.Memory.Value()
			coloredMemory := getColorValue(memoryUsage)

			row := []string{container.Name, fmt.Sprintf("%d", coloredCPU), fmt.Sprintf("%d", coloredMemory)}
			table.Append(row)

		}

	}
	table.Render()
}
