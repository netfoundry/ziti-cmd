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
	"fmt"
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/fabric/router/xgress_transport"
	"github.com/openziti/foundation/identity/dotziti"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	loop3_pb "github.com/openziti/ziti/ziti-fabric-test/subcmd/loop3/pb"
	"github.com/spf13/cobra"
	"net"
	"strings"
	"time"
)

func init() {
	dialerCmd := newDialerCmd()
	loop3Cmd.AddCommand(dialerCmd.cmd)
}

type dialerCmd struct {
	cmd            *cobra.Command
	identity       string
	endpoint       string
	direct         bool
	service        string
	edgeConfigFile string
}

func newDialerCmd() *dialerCmd {
	result := &dialerCmd{
		cmd: &cobra.Command{
			Use:   "dialer <scenarioFile>",
			Short: "Start loop2 dialer",
			Args:  cobra.ExactArgs(1),
		},
	}

	result.cmd.Run = result.run

	flags := result.cmd.Flags()
	flags.StringVarP(&result.identity, "identity", "i", "default", ".ziti/identities.yml name")
	flags.StringVarP(&result.endpoint, "endpoint", "e", "tls:127.0.0.1:7001", "Endpoint address")
	flags.BoolVarP(&result.direct, "direct", "d", false, "Transmit direct (no ingress)")
	flags.StringVarP(&result.service, "service", "s", "loop", "Service name for ingress")
	flags.StringVarP(&result.edgeConfigFile, "config-file", "c", "", "Edge SDK config file")

	return result
}

func (cmd *dialerCmd) run(_ *cobra.Command, args []string) {
	log := pfxlog.Logger()

	scenario, err := LoadScenario(args[0])
	if err != nil {
		panic(err)
	}
	log.Debug(scenario)

	if scenario.Metrics != nil {
		closer := make(chan struct{})
		if err := StartMetricsReporter(cmd.edgeConfigFile, scenario.Metrics, closer); err != nil {
			panic(err)
		}
		defer close(closer)
	}

	resultChs := make(map[string]chan *Result)
	for _, workload := range scenario.Workloads {
		log.Infof("executing workload [%s] with concurrency [%d]", workload.Name, workload.Concurrency)

		var conns []net.Conn
		for i := 0; i < int(workload.Concurrency); i++ {
			conns = append(conns, cmd.connect())
		}

		for i, conn := range conns {
			name := fmt.Sprintf("%s:%d", workload.Name, i)
			resultCh := make(chan *Result, 1)
			resultChs[name] = resultCh

			go func() {
				workload := scenario.Workloads[0]
				local := &loop3_pb.Test{
					Name:             name,
					TxRequests:       workload.Dialer.TxRequests,
					TxPacing:         workload.Dialer.TxPacing,
					TxMaxJitter:      workload.Dialer.TxMaxJitter,
					RxRequests:       workload.Listener.TxRequests,
					RxTimeout:        workload.Dialer.RxTimeout,
					PayloadMinBytes:  workload.Dialer.PayloadMinBytes,
					PayloadMaxBytes:  workload.Dialer.PayloadMaxBytes,
					LatencyFrequency: workload.Dialer.LatencyFrequency,
				}
				remote := &loop3_pb.Test{
					Name:             name,
					TxRequests:       workload.Listener.TxRequests,
					TxPacing:         workload.Listener.TxPacing,
					TxMaxJitter:      workload.Listener.TxMaxJitter,
					RxRequests:       workload.Dialer.TxRequests,
					RxTimeout:        workload.Listener.RxTimeout,
					PayloadMinBytes:  workload.Listener.PayloadMinBytes,
					PayloadMaxBytes:  workload.Listener.PayloadMaxBytes,
					LatencyFrequency: workload.Listener.LatencyFrequency,
				}

				if proto, err := newProtocol(conn); err == nil {
					if err := proto.txTest(remote); err == nil {
						if err := proto.run(local); err == nil {
							if result, err := proto.rxResult(); err == nil {
								resultCh <- result
							} else {
								panic(err)
							}
						} else {
							panic(err)
						}
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}
			}()

			time.Sleep(time.Duration(scenario.ConnectionDelay) * time.Millisecond)
		}
	}

	failed := false
	for name, resultCh := range resultChs {
		result := <-resultCh
		if !result.Success {
			failed = true
			log.Errorf("[%s] -> %s", name, result.Message)
		} else {
			log.Infof("[%s] -> success", name)
		}
	}
	if failed {
		panic("failures detected")
	} else {
		log.Info("success")
	}
}

func (cmd *dialerCmd) connect() net.Conn {
	log := pfxlog.Logger()

	start := time.Now()

	var conn net.Conn
	var err error
	if strings.HasPrefix(cmd.endpoint, "edge:") {
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

		service := strings.TrimPrefix(cmd.endpoint, "edge:")
		conn, err = context.DialWithOptions(service, &ziti.DialOptions{
			ConnectTimeout: time.Second * 30,
		})
		if err != nil {
			panic(err)
		}
	} else {
		endpoint, err := transport.ParseAddress(cmd.endpoint)
		if err != nil {
			panic(err)
		}

		_, id, err := dotziti.LoadIdentity(cmd.identity)
		if err != nil {
			panic(err)
		}

		if cmd.direct {
			if conn, err = dialDirect(endpoint, id); err != nil {
				panic(err)
			}
		} else {
			serviceId := &identity.TokenId{Token: cmd.service}
			if conn, err = dialIngress(endpoint, id, serviceId); err != nil {
				panic(err)
			}
		}
	}

	ConnectionTime.Update(time.Now().Sub(start))

	return conn
}

func dialDirect(endpoint transport.Address, id *identity.TokenId) (net.Conn, error) {
	peer, err := endpoint.Dial("loop", id, nil)
	if err != nil {
		return nil, err
	}

	return peer.Conn(), nil
}

func dialIngress(endpoint transport.Address, id, serviceId *identity.TokenId) (net.Conn, error) {
	peer, err := xgress_transport.ClientDial(endpoint, id, serviceId, nil)
	if err != nil {
		return nil, err
	}

	return peer.Conn(), nil
}
