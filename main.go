/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"arlon.io/arlon/cmd/bundle"
	"arlon.io/arlon/cmd/cluster"
	"arlon.io/arlon/cmd/clusterspec"
	"arlon.io/arlon/cmd/controller"
	"arlon.io/arlon/cmd/list_clusters"
	"arlon.io/arlon/cmd/profile"
	"github.com/spf13/cobra"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

func main() {

	ctx, cancel := Context()
	defer cancel()

	command := &cobra.Command{
		Use:               "arlon",
		Short:             "Run the Arlon program",
		Long:              "Run the Arlon program",
		DisableAutoGenTag: true,
		Run: func(c *cobra.Command, args []string) {
			c.Println(c.UsageString())
		},
	}
	// don't display usage upon error
	command.SilenceUsage = true
	command.AddCommand(controller.NewCommand(ctx))
	command.AddCommand(list_clusters.NewCommand())
	command.AddCommand(bundle.NewCommand())
	command.AddCommand(profile.NewCommand())
	command.AddCommand(clusterspec.NewCommand())
	command.AddCommand(cluster.NewCommand())

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	// override default log level, which is initially set to 'debug'
	flag.Set("zap-log-level", "info")
	flag.Parse()
	logger := zap.New(zap.UseFlagOptions(&opts))
	ctrl.SetLogger(logger)
	args := flag.Args()
	command.SetArgs(args)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func Context() (ctx context.Context, cancel func()) {
	log.Println("Setting up cancellable context")

	// trap Ctrl+C and call cancel on the context
	ctx, origCancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	cancel = func() {
		signal.Stop(signalChan)
		origCancel()
	}

	// Cancel the context on a received signal.
	go func() {
		select {
		case sig := <-signalChan:
			log.Printf("cancellation signal received: %s\n", sig.String())
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	return ctx, cancel
}
