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

package loop3

import (
	"errors"
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/agent"
	"github.com/openziti/foundation/identity/dotziti"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	loop3_pb "github.com/openziti/ziti/ziti-fabric-test/subcmd/loop3/pb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"strings"
)

func init() {
	listenerCmd := newListenerCmd()
	loop3Cmd.AddCommand(listenerCmd.cmd)
}

type listenerCmd struct {
	cmd            *cobra.Command
	identity       string
	bindAddress    string
	edgeConfigFile string
	test           *loop3_pb.Test
}

func newListenerCmd() *listenerCmd {
	result := &listenerCmd{
		cmd: &cobra.Command{
			Use:   "listener",
			Short: "Start loop2 listener",
			Args:  cobra.MaximumNArgs(1),
		},
	}

	result.cmd.Run = result.run

	flags := result.cmd.Flags()
	flags.StringVarP(&result.identity, "identity", "i", "default", ".ziti/identities.yml name")
	flags.StringVarP(&result.bindAddress, "bind", "b", "tcp:127.0.0.1:8171", "Listener bind address")
	flags.StringVarP(&result.edgeConfigFile, "config-file", "c", "", "Edge SDK config file")

	return result
}

func (cmd *listenerCmd) run(_ *cobra.Command, args []string) {
	log := pfxlog.Logger()

	var err error
	if err = agent.Listen(agent.Options{}); err != nil {
		pfxlog.Logger().WithError(err).Error("unable to start CLI agent")
	}

	var scenario *Scenario
	if len(args) == 1 {
		if scenario, err = LoadScenario(args[0]); err != nil {
			panic(err)
		}

		log.Debug(scenario)

		if len(scenario.Workloads) > 1 {
			panic(errors.New("only one workflow may be specified in a listener configuration"))
		} else if len(scenario.Workloads) == 1 {
			_, cmd.test = scenario.Workloads[0].GetTests()
		}
	}

	if scenario != nil && scenario.Metrics != nil {
		closer := make(chan struct{})
		if err := StartMetricsReporter(cmd.edgeConfigFile, scenario.Metrics, closer); err != nil {
			panic(err)
		}
		defer close(closer)
	}

	if strings.HasPrefix(cmd.bindAddress, "edge") {
		cmd.listenEdge()
	} else {
		_, id, err := dotziti.LoadIdentity(cmd.identity)
		if err != nil {
			panic(err)
		}

		bindAddress, err := transport.ParseAddress(cmd.bindAddress)
		if err != nil {
			panic(err)
		}

		cmd.listen(bindAddress, id)
	}
}

func (cmd *listenerCmd) listenEdge() {
	log := pfxlog.ContextLogger(cmd.bindAddress)
	defer log.Error("exited")
	log.Info("started")

	var context ziti.Context
	if cmd.edgeConfigFile != "" {
		zitiCfg, err := config.NewFromFile(cmd.edgeConfigFile)
		if err != nil {
			log.Fatalf("failed to load ziti configuration from %s: %v", cmd.edgeConfigFile, err)
		}
		context = ziti.NewContextWithConfig(zitiCfg)
	} else {
		context = ziti.NewContext()
	}

	service := strings.TrimPrefix(cmd.bindAddress, "edge:")
	listener, err := context.Listen(service)
	if err != nil {
		panic(err)
	}

	for {
		if conn, err := listener.Accept(); err != nil {
			panic(err)
		} else {
			go cmd.handle(conn, cmd.bindAddress)
		}
	}
}

func (cmd *listenerCmd) listen(bind transport.Address, i *identity.TokenId) {
	log := pfxlog.ContextLogger(bind.String())
	defer log.Error("exited")
	log.Info("started")

	incoming := make(chan transport.Connection)
	go func() {
		if _, err := bind.Listen("loop", i, incoming, nil); err != nil {
			panic(err)
		}
	}()
	for {
		select {
		case peer := <-incoming:
			if peer != nil {
				go cmd.handle(peer.Conn(), peer.Detail().String())
			} else {
				return
			}
		}
	}
}

func (cmd *listenerCmd) handle(conn net.Conn, context string) {
	log := pfxlog.ContextLogger(context)
	if proto, err := newProtocol(conn); err == nil {
		var test *loop3_pb.Test
		if cmd.test != nil && cmd.test.IsRxSequential() {
			test = cmd.test
		} else {
			if test, err = proto.rxTest(); err != nil {
				logrus.WithError(err).Error("failure receiving test parameters, closing")
				_ = conn.Close()
				return
			}
		}

		var result *Result
		if err := proto.run(test); err == nil {
			result = &Result{Success: true}
		} else {
			result = &Result{Success: false, Message: err.Error()}
		}
		if err := result.Tx(proto); err != nil {
			log.Errorf("unable to tx result (%s)", err)
		}

	} else {
		log.Errorf("error creating new protocol (%s)", err)
	}
}
