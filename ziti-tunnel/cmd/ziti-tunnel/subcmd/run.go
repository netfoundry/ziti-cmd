// +build linux

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
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/edge/tunnel/intercept"
	"github.com/openziti/edge/tunnel/intercept/tproxy"
	"github.com/openziti/edge/tunnel/intercept/tun"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run <config>",
	Short:   "Auto-select interceptor",
	Long:    "Provided for backwards compatibility with scripts that were coded around older ziti-tunnel versions.",
	Args:    cobra.MaximumNArgs(1),
	Run:     run,
	PostRun: rootPostRun,
}

func init() {
	root.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	log := pfxlog.Logger()
	var err error
	var tProxyInterceptor, tunInterceptor intercept.Interceptor

	if len(args) != 0 {
		_ = cmd.Flag("identity").Value.Set(args[0])
	}

	tProxyInterceptor, err = tproxy.New()
	if err != nil {
		log.Infof("tproxy initialization failed: %v", err)
	} else {
		log.Info("using tproxy interceptor")
		interceptor = tProxyInterceptor
		return
	}

	tunInterceptor, err = tun.New("", tunMtuDefault)
	if err != nil {
		log.Infof("tun initialization failed: %v", err)
	} else {
		log.Info("using tun interceptor")
		interceptor = tunInterceptor
		return
	}

	if interceptor == nil {
		log.Fatal("failed to initialize an interceptor")
	}
}
