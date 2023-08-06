package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"kom/kube"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var LogsCmd = &cobra.Command{
	Use:   "logs <pod-name>",
	Short: "Collect logs from a Kubernetes pod and save logs to a file",
	RunE:  collectLogsFromPod,
}

var saveToFile bool
var filePath string
var namespace string
var containerName string

func init() {
	LogsCmd.Flags().BoolVarP(&saveToFile, "save", "s", false, "Save the output to a log file")
	LogsCmd.Flags().StringVarP(&filePath, "output", "o", "output.log", "File path to save the output")
	LogsCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Specify the namespace of the pod")
	LogsCmd.Flags().StringVarP(&containerName, "container", "c", "", "Specify the name of the container in the pod (optional)")
}

func collectLogsFromPod(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide the pod name")
	}

	podName := args[0]

	config, err := kube.GetConfig()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	clientset, err := kube.NewClientSet(config)
	if err != nil {
		return fmt.Errorf("error creating Kubernetes client: %w", err)
	}

	var containerNamePtr string
	if containerName != "" {
		containerNamePtr = containerName
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: containerNamePtr,
	})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return fmt.Errorf("error opening stream for pod logs: %w", err)
	}
	defer podLogs.Close()

	currentTime := time.Now().Format(time.RFC3339)

	folderName := "komlogs"
	err = os.MkdirAll(folderName, 0755)
	if err != nil {
		return fmt.Errorf("error creating folder: %w", err)
	}

	fileName := fmt.Sprintf("%s_%s.log", podName, currentTime)

	filePath := filepath.Join(folderName, fileName)

	header := fmt.Sprintf("POD: %s, CONTAINER: %s, NAMESPACE: %s, DATE: %s\n",
		podName, containerNamePtr, namespace, currentTime)

	var buffer bytes.Buffer
	buffer.WriteString(header)

	if _, err := io.Copy(&buffer, podLogs); err != nil {
		return fmt.Errorf("error reading logs: %w", err)
	}

	if saveToFile {
		if err := saveLogsToFile(filePath, &buffer); err != nil {
			return fmt.Errorf("error saving logs to file: %w", err)
		}
	} else {
		displayLogs(&buffer)
	}
	return nil
}

func displayLogs(podsLogs io.Reader) {
	scanner := bufio.NewScanner(podsLogs)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}

func saveLogsToFile(filePath string, buffer *bytes.Buffer) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, buffer); err != nil {
		return fmt.Errorf("error writing logs to file: %w", err)
	}

	log.Printf("Logs saved to %s", filePath)
	return nil
}
