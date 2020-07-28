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
			// TODO detect whether this is running under kubectl or not, and modify help message accordingly
			//   see https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/#auth-plugins
			//   if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") { }
			runRootCmd(args)
			return nil
		},
	}

	cmd.Flags().StringVar(&args.LogLevel, "v", "info", "log level")

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
	level, err := log.ParseLevel(args.LogLevel)
	doOrDie(err)
	log.SetLevel(level)

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
		sort.Slice(pods, func(i, j int) bool {
			if pods[i].Namespace != pods[j].Namespace {
				return pods[i].Namespace < pods[j].Namespace
			}
			return pods[i].Name < pods[j].Name
		})
		fmt.Printf("%s:\n", node)
		for _, pod := range pods {
			fmt.Printf(" - %s/%s\n", pod.Namespace, pod.Name)
		}
		fmt.Println()
	}
}
