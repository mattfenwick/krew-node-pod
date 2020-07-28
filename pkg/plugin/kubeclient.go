package plugin

import (
	v1 "k8s.io/api/core/v1"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	client *kubernetes.Clientset
}

func PathToKubeConfig() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "unable to get home dir")
	}
	return path.Join(home, ".kube", "config"), nil
}

func NewDefaultClient() (*Client, error) {
	kubeConfigPath, err := PathToKubeConfig()
	if err != nil {
		return nil, err
	}
	return NewClient(kubeConfigPath)
}

func NewClientWithDefaultKubeConfigFallback(kubeConfigPath string) (*Client, error) {
	if kubeConfigPath == "" {
		return NewDefaultClient()
	}
	return NewClient(kubeConfigPath)
}

func NewClient(kubeConfigPath string) (*Client, error) {
	log.Debugf("instantiating k8s client from config path: '%s'", kubeConfigPath)
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	// kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to build config from flags")
	}
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to instantiate client")
	}
	return &Client{
		client: client,
	}, nil
}

func (kc *Client) ListPods(namespace string) (*v1.PodList, error) {
	pods, err := kc.client.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	return pods, errors.Wrapf(err, "unable to list pods in ns '%s'", namespace)
}

func (kc *Client) ListNodes() (*v1.NodeList, error) {
	nodes, err := kc.client.CoreV1().Nodes().List(metav1.ListOptions{})
	return nodes, errors.Wrapf(err, "unable to list nodes")
}
