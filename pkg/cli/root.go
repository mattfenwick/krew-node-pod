package cli

import (
	"encoding/json"
	"fmt"
	"github.com/mattfenwick/krew-node-pod/pkg/plugin"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func doOrDie(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func InitAndExecute() {
	rootCmd := setupRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}

type Config struct {
	LogLevel       string
	KubeFlags      *genericclioptions.ConfigFlags
	ShowContainers bool
	Format         string
	ShowStatus     bool
}

func setupRootCmd() *cobra.Command {
	args := &Config{}

	cmd := &cobra.Command{
		Use:   "node-pod",
		Short: "",
		Long:  `.`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, as []string) error {
			// TODO detect whether this is running under kubectl or not, and modify help message accordingly
			//   see https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/#auth-plugins
			//   if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") { }
			runRootCmd(args)
			return nil
		},
	}

	cmd.Flags().StringVar(&args.LogLevel, "v", "info", "log level")

	cmd.Flags().StringVar(&args.Format, "format", "table", "output format: one of json, list, table")

	cmd.Flags().BoolVarP(&args.ShowContainers, "containers", "c", false, "if true, print containers")

	cmd.Flags().BoolVar(&args.ShowStatus, "status", true, "if true, print object status (only works with list and table formats)")

	args.KubeFlags = genericclioptions.NewConfigFlags(false)
	args.KubeFlags.AddFlags(cmd.Flags())

	return cmd
}

func runRootCmd(args *Config) {
	level, err := log.ParseLevel(args.LogLevel)
	doOrDie(err)
	log.SetLevel(level)

	client, err := plugin.NewClientWithDefaultKubeConfigFallback(*args.KubeFlags.KubeConfig)
	doOrDie(err)

	output, err := FetchKubeData(client, *args.KubeFlags.Namespace)
	doOrDie(err)

	if !args.ShowContainers {
		output.RemoveContainers()
	}

	switch args.Format {
	case "list":
		fmt.Println(output.List(args.ShowStatus))
	case "json":
		fmt.Println(output.Json())
	case "table":
		output.Table(args.ShowContainers, args.ShowStatus).Render()
	default:
		doOrDie(errors.Errorf("invalid format '%s'", args.Format))
	}
}

type Output struct {
	Nodes []*Node
}

func NewOutput(nodeMap map[string]*Node) *Output {
	var nodes []*Node
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
	for _, node := range nodes {
		pods := node.Pods
		sort.Slice(pods, func(i, j int) bool {
			if pods[i].Namespace != pods[j].Namespace {
				return pods[i].Namespace < pods[j].Namespace
			}
			return pods[i].Name < pods[j].Name
		})
		node.Pods = pods
		for _, pod := range pods {
			containers := pod.Containers
			sort.Slice(containers, func(i, j int) bool {
				return containers[i].Name < containers[j].Name
			})
			pod.Containers = containers
		}
	}
	return &Output{Nodes: nodes}
}

func (o *Output) RemoveContainers() {
	for _, node := range o.Nodes {
		for _, pod := range node.Pods {
			pod.Containers = nil
		}
	}
}

func (o *Output) Json() string {
	bytes, err := json.MarshalIndent(o, "", "  ")
	doOrDie(err)
	return string(bytes)
}

func (o *Output) Table(showContainers bool, showStatus bool) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Node", "Namespace", "Pod Name"}
	if showContainers {
		headers = append(headers, "Container")
	}
	if showStatus {
		headers = append(headers, "Status")
	}
	table.SetHeader(headers)

	for _, node := range o.Nodes {
		nodeLine := []string{node.Name, "", ""}
		if showContainers {
			nodeLine = append(nodeLine, "")
		}
		if showStatus {
			nodeLine = append(nodeLine, node.Status)
		}
		table.Append(nodeLine)

		for _, pod := range node.Pods {
			podLine := []string{"", pod.Namespace, pod.Name}
			if showContainers {
				podLine = append(podLine, "")
			}
			if showStatus {
				podLine = append(podLine, pod.Status)
			}
			table.Append(podLine)

			if showContainers {
				for _, container := range pod.Containers {
					contLine := []string{"", "", "", container.Name}
					if showStatus {
						contLine = append(contLine, container.Status)
					}
					table.Append(contLine)
				}
			}
		}
	}

	return table
}

func (o *Output) List(showStatus bool) string {
	var lines []string
	for _, node := range o.Nodes {
		if showStatus {
			lines = append(lines, fmt.Sprintf("%s: %s", node.Name, node.Status))
		} else {
			lines = append(lines, fmt.Sprintf("%s", node.Name))
		}
		for _, pod := range node.Pods {
			if showStatus {
				lines = append(lines, fmt.Sprintf(" - %s/%s: %s", pod.Namespace, pod.Name, pod.Status))
			} else {
				lines = append(lines, fmt.Sprintf(" - %s/%s", pod.Namespace, pod.Name))
			}
			for _, container := range pod.Containers {
				if showStatus {
					lines = append(lines, fmt.Sprintf("   - %s: %s", container.Name, container.Status))
				} else {
					lines = append(lines, fmt.Sprintf("   - %s", container.Name))
				}
			}
		}
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

type Node struct {
	Name   string
	Pods   []*Pod
	Status string
}

func (n *Node) AddPod(pod *Pod) {
	n.Pods = append(n.Pods, pod)
}

type Pod struct {
	Name       string
	Namespace  string
	Containers []*Container
	Status     string
}

type Container struct {
	Name   string
	Status string
}

func FetchKubeData(client *plugin.Client, namespace string) (*Output, error) {
	kubePods, err := client.ListPods(namespace)
	if err != nil {
		return nil, err
	}

	kubeNodes, err := client.ListNodes()
	if err != nil {
		return nil, err
	}

	// get the nodes
	nodes := map[string]*Node{}
	for _, kubeNode := range kubeNodes.Items {
		var status = "unknown"
		if len(kubeNode.Status.Conditions) > 0 {
			status = string(kubeNode.Status.Conditions[len(kubeNode.Status.Conditions)-1].Type)
		}
		nodes[kubeNode.Name] = &Node{
			Name:   kubeNode.Name,
			Pods:   nil,
			Status: status,
		}
	}

	// add the pods into the nodes
	for _, kubePod := range kubePods.Items {
		nodeName := kubePod.Spec.NodeName
		node, ok := nodes[nodeName]
		if !ok {
			node = &Node{
				Name:   nodeName,
				Pods:   nil,
				Status: "Unknown",
			}
			nodes[nodeName] = node
			log.Warnf("pod %s/%s assigned to node %s -- but node not found in kube", kubePod.Namespace, kubePod.Name, nodeName)
		}
		node.AddPod(extractPod(&kubePod))
	}

	return NewOutput(nodes), err
}

func extractPod(kubePod *v1.Pod) *Pod {
	var containers []*Container
	for _, cont := range kubePod.Status.ContainerStatuses {
		containers = append(containers, extractContainer(&cont))
	}
	return &Pod{
		Name:       kubePod.Name,
		Namespace:  kubePod.Namespace,
		Containers: containers,
		Status:     string(kubePod.Status.Phase),
	}
}

func extractContainer(kubeContainer *v1.ContainerStatus) *Container {
	state := "Unknown"
	if kubeContainer.State.Running != nil {
		state = "Running"
	} else if kubeContainer.State.Terminated != nil {
		state = "Terminated"
	} else if kubeContainer.State.Waiting != nil {
		state = "Waiting"
	}
	return &Container{
		Name:   kubeContainer.Name,
		Status: state,
	}
}
