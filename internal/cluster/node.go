package cluster

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var nodeInfo NodeInfo

const (
	NodeExternalIP corev1.NodeAddressType = "ExternalIP"
)

// Node info structure.
// Protected by a mutex and contains information about the kubeslice gateway nodes in the cluster.
type NodeInfo struct {
	Client     client.Client
	ExternalIP []string
	sync.Mutex
}

//GetNodeExternalIpList gets the list of External Node IPs of avesha-gateway nodes

func GetNodeExternalIpList() ([]string, error) {
	// If node IP is set as an env variable, we use that as the only
	// node IP available to us. Early exit from here, and there is no need
	// spawn the node watcher thread.
	staticNodeIp := os.Getenv("NODE_IP")
	if staticNodeIp != "" {
		nodeInfo.ExternalIP = append(nodeInfo.ExternalIP, staticNodeIp)
		return nodeInfo.ExternalIP, nil
	}
	// Dynamic node IP deduction if there is no static node IP provided
	nodeInfo.Lock()
	defer nodeInfo.Unlock()

	if len(nodeInfo.ExternalIP) == 0 {
		err := nodeInfo.populateNodeIpList()
		if err != nil {
			return nil, err
		}
	}
	return nodeInfo.ExternalIP, nil
}

func (n *NodeInfo) populateNodeIpList() error {
	nodes := corev1.NodeList{}
	listOpts := []client.ListOption{
		client.MatchingLabels(map[string]string{
			"avesha/node-type": "gateway",
		}),
	}
	err := n.Client.List(context.Background(), &nodes, listOpts...)
	if err != nil {
		return fmt.Errorf("can't fetch node list: %+v ", err)
	}
	//TODO(rahulsawra):check if we can optimize this
	nodeIpArr := []corev1.NodeAddress{}
	for i := 0; i < len(nodes.Items); i++ {
		nodeIpArr = append(nodeIpArr, nodes.Items[i].Status.Addresses...)
	}

	for i := 0; i < len(nodeIpArr); i++ {
		if nodeIpArr[i].Type == NodeExternalIP {
			n.ExternalIP = append(n.ExternalIP, nodeIpArr[i].Address)
		}
	}

	return err
}
