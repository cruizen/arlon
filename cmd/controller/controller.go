package controller

import (
	"context"

	"arlon.io/arlon/pkg/argocd"
	"arlon.io/arlon/pkg/controller"
	"github.com/spf13/cobra"
)

func NewCommand(ctx context.Context) *cobra.Command {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	// var argocdConfigPath string
	var argocdUsername string
	var argocdPassword string
	var argocdServerUri string

	argocdServer := argocd.NewArgoCDImpl(argocdUsername, argocdPassword, argocdServerUri)

	command := &cobra.Command{
		Use:               "controller",
		Short:             "Run the Arlon controller",
		Long:              "Run the Arlon controller",
		DisableAutoGenTag: true,
		Run: func(c *cobra.Command, args []string) {
			controller.StartController(ctx, *argocdServer,
				metricsAddr, probeAddr, enableLeaderElection)
		},
	}
	//command.Flags().StringVar(&argocdConfigPath, "argocd-config-path", "", "argocd configuration file path")
	command.Flags().StringVar(&argocdUsername, "argocd-username", "arlon", "username for argocd server")
	command.Flags().StringVar(&argocdPassword, "argocd-password", "", "password for argocd server")
	command.Flags().StringVar(&argocdServerUri, "argocd-server", "argocd-server:80", "uri of argocd server")
	// Ability to mark flags as required or exclusive as a group is new in Cobra https://github.com/spf13/cobra/pull/1654
	// We can use that once a new release is out. It's not available in Cobra 1.4.0
	//command.MarkFlagsRequiredTogether("argocd-username", "argocd-password", "argocd-server")

	command.Flags().StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	command.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	command.Flags().BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	return command
}
