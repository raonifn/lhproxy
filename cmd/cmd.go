package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/murer/lhproxy/pipe"
	"github.com/murer/lhproxy/server"
	"github.com/murer/lhproxy/sockets"

	"github.com/murer/lhproxy/util"
)

var rootCmd *cobra.Command
var clientCmd *cobra.Command
var pipeCmd *cobra.Command

func Config() {
	rootCmd = &cobra.Command{
		Use: "lhproxy", Short: "Last Hope Proxy",
		Version: fmt.Sprintf("%s-%s:%s", runtime.GOOS, runtime.GOARCH, util.Version),
	}
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet")
	cobra.OnInitialize(gconf)

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf(rootCmd.Version)
			return nil
		},
	})

	configServer()

	clientCmd = &cobra.Command{Use: "client"}
	rootCmd.AddCommand(clientCmd)

	pipeCmd = &cobra.Command{Use: "pipe"}
	clientCmd.AddCommand(pipeCmd)

	configPipe()
}

func gconf() {
	quiet, err := rootCmd.PersistentFlags().GetBool("quiet")
	util.Check(err)
	if quiet {
		log.SetOutput(ioutil.Discard)
	}
}

func configServer() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "server <host>:<port>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			server.Start(args[0], util.Secret())
			return nil
		},
	})
}

func configPipe() {
	pipeCmd.AddCommand(&cobra.Command{
		Use:  "native <host>:<port>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := &pipe.Pipe{
				Scks:    sockets.GetNative(),
				Address: args[0],
				Reader:  os.Stdin,
				Writer:  os.Stdout,
			}
			p.Execute()
			return nil
		},
	})

	pipeCmd.AddCommand(&cobra.Command{
		Use:  "http http://<lhproxy:port>/ <host>:<port>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := &pipe.Pipe{
				Scks: &server.HttpSockets{
					URL:    args[0],
					Secret: util.Secret(),
				},
				Address: args[1],
				Reader:  os.Stdin,
				Writer:  os.Stdout,
			}
			p.Execute()
			return nil
		},
	})
}

func Execute() {
	err := rootCmd.Execute()
	util.Check(err)
}
