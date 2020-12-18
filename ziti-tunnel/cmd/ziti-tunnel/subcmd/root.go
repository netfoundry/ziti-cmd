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
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/edge/tunnel/dns"
	"github.com/openziti/edge/tunnel/entities"
	"github.com/openziti/edge/tunnel/intercept"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/openziti/ziti/common/enrollment"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	svcPollRateFlag   = "svcPollRate"
	resolverCfgFlag   = "resolver"
	dnsSvcIpRangeFlag = "dnsSvcIpRange"
)

func init() {
	root.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose mode")
	root.PersistentFlags().StringP("identity", "i", "", "Path to JSON file that contains an enrolled identity")
	root.PersistentFlags().StringP("identityDir", "", "", "Path to directory of JSON files that each contain an enrolled identity")
	root.PersistentFlags().Uint(svcPollRateFlag, 15, "Set poll rate for service updates (seconds)")
	root.PersistentFlags().StringP(resolverCfgFlag, "r", "udp://127.0.0.1:53", "Resolver configuration")
	root.PersistentFlags().StringVar(&logFormatter, "log-formatter", "", "Specify log formatter [json|pfxlog|text]")
	root.PersistentFlags().StringP(dnsSvcIpRangeFlag, "d", "100.64.0.1/10", "cidr to use when assigning IPs to unresolvable intercept hostnames")

	root.AddCommand(enrollment.NewEnrollCommand())
}

var root = &cobra.Command{
	Use:              filepath.Base(os.Args[0]),
	Short:            "Ziti Tunnel",
	PersistentPreRun: rootPreRun,
}

var interceptor intercept.Interceptor
var resolver dns.Resolver
var logFormatter string

func Execute() {
	if err := root.Execute(); err != nil {
		pfxlog.Logger().Errorf("error: %s", err)
		os.Exit(1)
	}
}

func rootPreRun(cmd *cobra.Command, _ []string) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		println("err")
	}
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
}

func getIdFiles(cmd *cobra.Command) []string {
	log := pfxlog.Logger()

	var idFiles []string
	idDir := cmd.Flag("identityDir").Value.String()
	if idDir != "" {
		files, err := ioutil.ReadDir(idDir)
		if err != nil {
			log.Fatalf("failed to scan directory %s: %v", idDir, err)
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".json" {
				fn, err := filepath.Abs(filepath.Join(idDir, file.Name()))
				if err != nil {
					log.Fatalf("failed to listing file %s: %v", file.Name(), err)
				}
				idFiles = append(idFiles, fn)
			}
		}
	}

	identityJson := cmd.Flag("identity").Value.String()
	if identityJson != "" {
		idFiles = append(idFiles, identityJson)
	}
	return idFiles
}
func rootPostRun(cmd *cobra.Command, _ []string) {
	log := pfxlog.Logger()

	idFiles := getIdFiles(cmd)
	if len(idFiles) == 0 {
		log.Fatalf("no identityJson files found")
	}

	svcPollRate, _ := cmd.Flags().GetUint(svcPollRateFlag)
	resolverConfig := cmd.Flag("resolver").Value.String()
	resolver = dns.NewResolver(resolverConfig)
	dnsIpRange, _ := cmd.Flags().GetString(dnsSvcIpRangeFlag)
	err := intercept.SetDnsInterceptIpRange(dnsIpRange)
	if err != nil {
		log.Fatalf("invalid dns service IP range %s: %v", dnsIpRange, err)
	}

	var wg sync.WaitGroup
	for _, f := range idFiles {
		cfg, err := config.NewFromFile(f)
		if err != nil {
			log.Fatalf("failed to load ziti configuration from %s: %v", f, err)
		}

		cfg.ConfigTypes = []string{
			entities.ClientConfigV1,
			entities.ServerConfigV1,
		}
		ztx := ziti.NewContextWithConfig(cfg)

		wg.Add(1)
		go intercept.ServicePoller(ztx, interceptor, resolver, time.Duration(svcPollRate)*time.Second)
	}
	wg.Wait()
}
