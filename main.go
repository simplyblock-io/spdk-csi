package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	e2elog "k8s.io/kubernetes/test/e2e/framework/log"
)

var wg *sync.WaitGroup

const (
	MGMT_IP        = ""
	CLUSTER_ID     = ""
	CLUSTER_SECRET = ""
)

func executeKubectlCommand(command string) (string, error) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return "", fmt.Errorf("error executing kubectl command: %v, output: %s", err, string(out))
	}
	return string(out), nil
}

func checkCachingNodes(timeout time.Duration) error {
	err := wait.PollImmediate(3*time.Second, timeout, func() (bool, error) {
		fmt.Println("-- caching nodes --")
		out, err := executeKubectlCommand("kubectl get nodes -l type=cache")
		if err != nil {
			e2elog.Logf("failed %s", err)
			return false, err
		}
		fmt.Println(out)

		//deployCachenode()
		// _, err = executeKubectlCommand("apply -f caching-node.yaml")
		// if err != nil {
		// 	e2elog.Logf("failed %s", err)
		// 	return false, err
		// }

		out, err = executeKubectlCommand("kubectl wait --timeout=3m --for=condition=ready pod -l app=caching-node")
		if err != nil {
			e2elog.Logf("failed %s", err)
			return false, err
		}
		fmt.Println(out)

		out, err = executeKubectlCommand("kubectl get pods -l app=caching-node -o wide | awk 'NR>1 {print $(NF-3)}'")
		if err != nil {
			fmt.Println("this is the cause of the error")
			//fmt.Println(out)
			e2elog.Logf("failed %s", err)
			return false, err
		}
		fmt.Println(out)

		nodeIPs := strings.Split(strings.TrimSpace(out), "\n")
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Println(nodeIPs)
		for _, node := range nodeIPs {
			fmt.Printf("Adding caching node: %s\n", node)

			// cmd := exec.Command("curl", "--location", fmt.Sprintf("http://%s/cachingnode/", MGMT_IP), "--header", "Content-Type: application/json", "--header", fmt.Sprintf("Authorization: %s %s", CLUSTER_ID, CLUSTER_SECRET), "--data", fmt.Sprintf("{\"cluster_id\": \"%s\", \"node_ip\": \"%s:5000\", \"iface_name\": \"eth0\", \"spdk_mem\": \"2g\"}", CLUSTER_ID, node))
			// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
			// _, err = cmd.CombinedOutput()
			// if err != nil {
			// 	e2elog.Logf("failed %s", err)
			// 	return false, err
			// }
			////////////////////////////////////////
			type output struct {
				out []byte
				err error
			}

			ch := make(chan output)

			go func() {
				//cmd := exec.Command("curl", "--location", fmt.Sprintf("http://%s/cachingnode/", MGMT_IP), "--header", "Content-Type: application/json", "--header", fmt.Sprintf("Authorization: %s %s", CLUSTER_ID, CLUSTER_SECRET), "--data", fmt.Sprintf("{\"cluster_id\": \"%s\", \"node_ip\": \"%s:5000\", \"iface_name\": \"eth0\", \"spdk_mem\": \"2g\"}", CLUSTER_ID, node))
				//cmd := exec.Command(fmt.Sprintf("curl --location \"http://%s/cachingnode/\" --header \"Content-Type: application/json\" --header \"Authorization: %s %s\" --data '{\"cluster_id\": \"%s\", \"node_ip\": \"%s:5000\", \"iface_name\": \"eth0\", \"spdk_mem\": 8589934592}'", MGMT_IP, CLUSTER_ID, CLUSTER_SECRET, CLUSTER_ID, node))
				out, err := exec.Command("bash", "-c", "curl --location \"http://%s/cachingnode/\" --header \"Content-Type: application/json\" --header \"Authorization: %s %s\" --data '{\"cluster_id\": \"%s\", \"node_ip\": \"%s:5000\", \"iface_name\": \"eth0\", \"spdk_mem\": 8589934592}'", MGMT_IP, CLUSTER_ID, CLUSTER_SECRET, CLUSTER_ID, node).Output()
				//out, _ = cmd.CombinedOutput()
				ch <- output{out, err}
			}()

			select {
			case <-time.After(2 * time.Second):
				fmt.Println("timed out")
			case x := <-ch:
				fmt.Printf("program done; out: %q\n", string(x.out))
				if x.err != nil {
					fmt.Printf("program errored: %s\n", x.err)
				}
			}
			///////////////////////////////////////
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("failed to create cache node %w", err)
	}
	return nil
}

func main() {

	if err := checkCachingNodes(5 * time.Minute); err != nil {
		fmt.Println("Error:", err)
	}
}
