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

package cmd

import (
	"io"

	"github.com/blang/semver"
	"github.com/openziti/ziti/common/version"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/templates"
	c "github.com/openziti/ziti/ziti/cmd/ziti/constants"
	"github.com/openziti/ziti/ziti/cmd/ziti/internal/log"
	"github.com/spf13/cobra"
)

var (
	installZitiEdgeTunnelLong = templates.LongDesc(`
		Installs the Ziti Edge Tunnel app if it has not been installed already
`)

	installZitiEdgeTunnelExample = templates.Examples(`
		# Install the Ziti Edge Tunnel app 
		ziti install ziti-edge-tunnel
	`)
)

// InstallZitiEdgeTunnelOptions the options for the upgrade ziti-edge-tunnel command
type InstallZitiEdgeTunnelOptions struct {
	InstallOptions

	Version string
}

// NewCmdInstallZitiEdgeTunnel defines the command
func NewCmdInstallZitiEdgeTunnel(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &InstallZitiEdgeTunnelOptions{
		InstallOptions: InstallOptions{
			CommonOptions: CommonOptions{
				Factory: f,
				Out:     out,
				Err:     errOut,
			},
		},
	}

	cmd := &cobra.Command{
		Use:     "ziti-edge-tunnel",
		Short:   "Installs the Ziti Edge Tunnel app - if it has not been installed already",
		Aliases: []string{"edge-tunnel"},
		Long:    installZitiEdgeTunnelLong,
		Example: installZitiEdgeTunnelExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdhelper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&options.Version, "version", "v", "", "The specific version to install")
	options.addCommonFlags(cmd)
	return cmd
}

// Run implements the command
func (o *InstallZitiEdgeTunnelOptions) Run() error {
	newVersion, err := o.getLatestGitHubReleaseVersion(version.GetBranch(), c.ZITI_EDGE_TUNNEL_GITHUB)
	if err != nil {
		return err
	}

	if o.Version != "" {
		newVersion, err = semver.Make(o.Version)
	}

	log.Infoln("Attempting to install '" + c.ZITI_EDGE_TUNNEL + "' version: " + newVersion.String())

	return o.installGitHubRelease(version.GetBranch(), c.ZITI_EDGE_TUNNEL, c.ZITI_EDGE_TUNNEL_GITHUB, false, newVersion.String())

}
