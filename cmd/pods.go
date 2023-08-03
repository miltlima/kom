package cmd

import (
	"context"
	"fmt"
	"kom/kube"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var pods *corev1.PodList

var PodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Show metrics for all pods",
	Run:   showPodMetrics,
}

func showPodMetrics(cmd *cobra.Command, args []string) {
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

	pods, err = clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error to list pods", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Node", "Namespace", "Pod", "Pod IP", "CPU Usage %", "Memory Usage %"})

	for _, pod := range pods.Items {
		podName := pod.Name
		podNamespace := pod.Namespace
		podIP := pod.Status.PodIP
		nodeName := pod.Spec.NodeName

		podMetrics, err := metricsClientset.MetricsV1beta1().PodMetricses(podNamespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error to get metrics from pod %s%s: %s\n", podNamespace, podName, err)
			continue
		}

		node, err := clientset.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error to get node %s: %s\n", nodeName, err)
			continue
		}

		cpuUsage := podMetrics.Containers[0].Usage.Cpu().MilliValue() / 10
		coloredCPU := getColorValue(int(cpuUsage))

		nodeMemoryCapacity := float64(node.Status.Capacity.Memory().Value())
		podMemoryUsage := float64(podMetrics.Containers[0].Usage.Memory().Value())

		memoryUsage := (podMemoryUsage / nodeMemoryCapacity) * 100.0
		coloredMemory := getColorValue(int(memoryUsage))

		row := []string{nodeName, podNamespace, podName, podIP, coloredCPU, coloredMemory}
		table.Append(row)

	}

	table.Render()

}
