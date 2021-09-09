package k8sexec

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/argoproj/argo-workflows/v3/errors"
	"github.com/argoproj/argo-workflows/v3/workflow/common"
)

func KillContainer(restConfig *rest.Config, pod *corev1.Pod, containerName string, sig syscall.Signal) error {
	command := []string{"/bin/sh", "-c", fmt.Sprintf("kill -%d 1", sig)}
	exec, err := common.ExecPodContainer(restConfig, pod.Namespace, pod.Name, containerName, true, true, command...)
	if err != nil {
		return err
	}
	_, _, err = GetExecutorOutput(exec)
	return err
}

// GetExecutorOutput returns the output of an remotecommand.Executor
func GetExecutorOutput(exec remotecommand.Executor) (*bytes.Buffer, *bytes.Buffer, error) {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	err := exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdOut,
		Stderr: &stdErr,
		Tty:    false,
	})
	if err != nil {
		return nil, nil, errors.InternalWrapError(err)
	}
	return &stdOut, &stdErr, nil
}
func TestKillContainer(t *testing.T) {
	cliConfig, err := GetConfig("", "")
	if err != nil {
		log.Fatalf("new kube client error: %v", err)
	}

	kubecli, err := kubernetes.NewForConfig(cliConfig)
	if err != nil {
		log.Fatalf("NewForConfig error: %v", err)
	}

	podList, err := kubecli.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("list pod error:%s", err.Error())
	}

	var wg sync.WaitGroup
	count := 0
	for i, pod := range podList.Items {
		if !strings.HasPrefix(pod.Name, "exec-test-") {
			continue
		}
		if pod.Status.Conditions == nil {
			continue
		}

		isRunning := false
		for _, c := range pod.Status.Conditions {
			if c.Type == "Ready" && c.Status == "True" {
				isRunning = true
				break
			}
		}
		if !isRunning {
			continue
		}
		wg.Add(1)
		go func(pod *corev1.Pod) {
			defer wg.Done()
			stime := time.Now()
			if err := KillContainer(cliConfig, pod, pod.Spec.Containers[0].Name, 15); err != nil {
				log.Printf("KillContainer pod:%s/%s error:%s latency:%fs\n", pod.Namespace, pod.Name, err.Error(), float64(time.Now().UnixNano()-stime.UnixNano())/1e9)
			} else {
				t.Logf("pod: %s/%s kill latency:%fs ", pod.Namespace, pod.Name, float64(time.Now().UnixNano()-stime.UnixNano())/1e9)
			}

		}(&podList.Items[i])

		count++
		if count >= 50 {
			break
		}
	}
	wg.Wait()
	t.Logf("all pod killed~~")
}

// GetConfig returns a rest.Config to be used for kubernetes client creation.
// It does so in the following order:
//   1. Use the passed kubeconfig/masterURL.
//   2. Fallback to the KUBECONFIG environment variable.
//   3. Fallback to in-cluster config.
//   4. Fallback to the ~/.kube/config.
func GetConfig(masterURL, kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	// If we have an explicit indication of where the kubernetes config lives, read that.
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	}
	// If not, try the in-cluster config.
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}
	// If no in-cluster config, try the default location in the user's home directory.
	if usr, err := user.Current(); err == nil {
		if c, err := clientcmd.BuildConfigFromFlags("", filepath.Join(usr.HomeDir, ".kube", "config")); err == nil {
			return c, nil
		}
	}

	return nil, fmt.Errorf("could not create a valid kubeconfig")
}
