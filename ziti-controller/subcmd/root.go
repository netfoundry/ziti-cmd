/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package subcmd

import (
	"fmt"
	"github.com/michaelquigley/pfxlog"
	edgeSubCmd "github.com/openziti/edge/controller/subcmd"
	"github.com/openziti/ziti/common/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	root.PersistentFlags().BoolVarP(&cliAgentEnabled, "cliagent", "a", false, "Enable CLI Agent (use in dev only)")
	root.PersistentFlags().StringVar(&logFormatter, "log-formatter", "", "Specify log formatter [json|pfxlog|text]")

	edgeSubCmd.AddCommands(root, version.GetCmdBuildInfo())
}

var root = &cobra.Command{
	Use:   "ziti-controller",
	Short: "Ziti Fabric Controller",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}

		switch logFormatter {
		case "pfxlog":
			logrus.SetFormatter(pfxlog.NewFormatterStartingToday())
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		case "text":
			logrus.SetFormatter(&logrus.TextFormatter{})
		default:
			// let logrus do its own thing
		}

	},
}
var verbose bool
var cliAgentEnabled bool
var logFormatter string

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
