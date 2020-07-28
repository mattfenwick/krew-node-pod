package cli

import (
	"fmt"
	"github.com/mattfenwick/krew-node-pod/pkg/plugin"
	v1 "k8s.io/api/core/v1"
	//"os"
	"sort"

	//"strings"
	//"time"

	log "github.com/sirupsen/logrus"
	//"github.com/mattfenwick/krew-node-pod/pkg/logger"
	//"github.com/mattfenwick/krew-node-pod/pkg/plugin"
	//"github.com/pkg/errors"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
	//"github.com/tj/go-spin"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

func doOrDie(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

type Config struct {
	LogLevel  string
	KubeFlags *genericclioptions.ConfigFlags
}

func setupRootCmd() *cobra.Command {
	args := &Config{}

	cmd := &cobra.Command{
		Use:           "node-pod",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, as []string) error {
			runRootCmd(args)
			return nil
		},
	}

	args.KubeFlags = genericclioptions.NewConfigFlags(false)
	args.KubeFlags.AddFlags(cmd.Flags())

	return cmd
}

func InitAndExecute() {
	rootCmd := setupRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func runRootCmd(args *Config) {
	// TODO use args
	ListPods()
}

func ListPods() {
	client, err := plugin.NewDefaultClient()
	doOrDie(err)

	pods, err := client.ListPods(v1.NamespaceAll)
	doOrDie(err)

	nodeToPods := map[string][]v1.Pod{}
	for _, pod := range pods.Items {
		node := pod.Spec.NodeName
		nodeToPods[node] = append(nodeToPods[node], pod)
	}

	nodes := []string{}
	for node := range nodeToPods {
		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i] < nodes[j]
	})

	for _, node := range nodes {
		pods := nodeToPods[node]
		fmt.Printf("%s:\n", node)
		for _, pod := range pods {
			fmt.Printf(" - %s/%s\n", pod.Namespace, pod.Name)
		}
		fmt.Println()
	}
}
