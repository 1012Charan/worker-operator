package cluster

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type Cluster struct {
	Client kubernetes.Interface
	Name   string `json:"clusterName,omitempty"`
}

//NewCluster returns ClusterInterface
func NewCluster(client kubernetes.Interface, clusterName string) ClusterInterface {
	return &Cluster{
		Client: client,
		Name:   clusterName,
	}
}

func (c *Cluster) GetClusterInfo(ctx context.Context) (*ClusterInfo, error) {
	cl := &ClusterInfo{
		Name: c.Name,
	}
	location, err := c.getClusterLocation(ctx)
	if err != nil {
		return nil, err
	}
	cl.ClusterProperty = ClusterProperty{
		GeoLocation: location,
	}
	return cl, nil
}

func (c *Cluster) getClusterLocation(ctx context.Context) (GeoLocation, error) {
	var g GeoLocation

	nodeList, err := c.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return g, fmt.Errorf("can't fetch node list: %+v ", err)
	}

	if len(nodeList.Items) == 0 {
		return g, fmt.Errorf("can't fetch node list , length of node items is zero")
	}

	g.CloudRegion = nodeList.Items[0].ObjectMeta.Labels["topology.kubernetes.io/region"]

	if nodeList.Items[0].Spec.ProviderID != "" {
		g.CloudProvider = strings.Split(nodeList.Items[0].Spec.ProviderID, ":")[0]
	}

	return g, nil
}
