package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rimasid/kubesay/internal/k8s"
	"github.com/rimasid/kubesay/internal/output"
)

// This runs automatically when the package loads, and adds our command to the root.
func init() {
	registerCommand(newListCmd)
}

func newListCmd(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Kubernetes resources",
	}
	cmd.AddCommand(newListPodsCmd(v))
	return cmd
}

func newListPodsCmd(v *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:   "pods",
		Short: "List pods in a namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, defaultNS, err := k8s.NewClient("")
			if err != nil {
				return err
			}

			ns := v.GetString(keyNamespace)
			if ns == "" {
				ns = defaultNS
			}

			pods, err := k8s.ListPods(cmd.Context(), client, ns)
			if err != nil {
				return err
			}

			return output.RenderPods(cmd.OutOrStdout(), v.GetString(keyOutputFormat), pods)
		},
	}
}
