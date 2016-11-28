package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/infrakit/pkg/cli"
	instance_plugin "github.com/docker/infrakit/pkg/rpc/instance"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	// PluginName is the name of the plugin in the Docker Hub / registry
	PluginName = "ProfitBricksInstance"

	// PluginType is the type / interface it supports
	PluginType = "infrakit.InstancePlugin/1.0"

	// Version is the build release identifier.
	Version = "1.0.0"

	// Revision is the build source control revision.
	Revision = "Unspecified"
)

func main() {
	builder := &Builder{}
	var logLevel int
	var name string

	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "ProfitBricks instance plugin",
		Run: func(c *cobra.Command, args []string) {
			instancePlugin, err := builder.BuildInstancePlugin()
			if err != nil {
				panic(err)
			}
			cli.SetLogLevel(logLevel)
			cli.RunPlugin(name, instance_plugin.PluginServer(instancePlugin))
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print build version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			buff, err := json.MarshalIndent(map[string]interface{}{
				"name":     PluginName,
				"type":     PluginType,
				"version":  Version,
				"revision": Revision,
			}, "  ", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(buff))
			return nil
		},
	})

	cmd.Flags().IntVar(&logLevel, "log", cli.DefaultLogLevel, "Logging level. 0 is least verbose. Max is 5")
	cmd.Flags().StringVar(&name, "name", "infrakit-instance-profitbricks", "Plugin name to advertise for discovery")

	cmd.Flags().AddFlagSet(builder.Flags())

	cmd.AddCommand(cli.VersionCommand())

	err := cmd.Execute()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
