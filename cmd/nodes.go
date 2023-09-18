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
	"k8s.io/client-go/kubernetes"
)

var nodes *corev1.NodeList

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Show metrics for all nodes",
	Run:   showNodesMetrics,
}

func getNodeLabels(node corev1.Node) []string {
	var labels []string
	for key, value := range node.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, value))
	}
	return labels
}

func getKubernetesVersion(clientset *kubernetes.Clientset) (string, error) {
	versionInfo, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return versionInfo.GitVersion, nil
}

func getPodCount(node corev1.Node, clientset *kubernetes.Clientset) int {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + node.Name,
	})
	if err != nil {
		fmt.Printf("Error getting pods on node %s: %s\n", node.Name, err)
		return 0
	}
	return len(pods.Items)
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

	kubernetesVersion, err := getKubernetesVersion(clientset)
	if err != nil {
		fmt.Printf("Error getting Kubernetes version: %s\n", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Node", "K8s Version", "IP", "Pod Count", "CPU Usage %", "Memory Usage %", "h", "Status", "Labels", "Taints"})

	for _, node := range nodes.Items {
		nodeName := node.Name
		hasLabels := len(node.Labels) > 0
		hasTaints := len(node.Spec.Taints) > 0
		podCount := getPodCount(node, clientset)

		var hasNOSchedule bool
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.TaintNodeUnschedulable && condition.Status == corev1.ConditionTrue {
				hasNOSchedule = true
				break
			}
		}

		var taintValue string
		if hasTaints {
			taintValue = "Yes"
		} else {
			taintValue = "No"
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

		var statusLabels string

		if hasLabels {
			labels := strings.Join(getNodeLabels(node), ", ")
			statusLabels += labels

		}

		row := []string{
			nodeName,
			kubernetesVersion,
			strings.Join(nodeIPs,
				", "),
			fmt.Sprintf("%d", podCount),
			coloredCPU,
			coloredMemory,
			statusEmoji,
			statusMessage,
			statusLabels,
			fmt.Sprintf("%v", taintValue)}
		table.Append(row)

	}
	table.Render()
}
