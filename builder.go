package main

import (
	"fmt"
	"github.com/docker/infrakit/pkg/spi/instance"
	"github.com/spf13/pflag"
	"os"
)

type options struct {
	username string
	password string
	dir      string
}

//Builder test
type Builder struct {
	options options
}

//Flags test
func (b *Builder) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("oneandone", pflag.PanicOnError)
	flags.StringVar(&b.options.username, "username", "", "ProfitBricks username")
	flags.StringVar(&b.options.password, "password", "", "ProfitBricks username")
	flags.StringVar(&b.options.dir, "dir", os.TempDir(), "Existing directory for storing the plugin files")
	return flags
}

// BuildInstancePlugin creates an instance Provisioner configured with the Flags.
func (b *Builder) BuildInstancePlugin() (instance.Plugin, error) {
	var ok bool

	if len(b.options.username) == 0 {
		if b.options.username, ok = os.LookupEnv("PROFITBRICKS_USERNAME"); !ok {
			fmt.Errorf("'PROFITBRICKS_USERNAME' is not set")
			os.Exit(1)
		}
	}

	if len(b.options.password) == 0 {
		if b.options.password, ok = os.LookupEnv("PROFITBRICKS_PASSWORD"); !ok {
			fmt.Errorf("'PROFITBRICKS_PASSWORD' is not set")
			os.Exit(1)
		}
	}

	if len(b.options.dir) == 0 {
		b.options.dir = "./"
	}

	return InstancePlugin(b.options.username, b.options.password, b.options.dir), nil
}
