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

package edge_controller

import (
	"fmt"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/openziti/ziti/ziti/cmd/ziti/util"
	"github.com/spf13/cobra"
	"io"
)

// versionOptions are the flags for version commands
type versionOptions struct {
	commonOptions
}

// newVersionCmd creates the command
func newVersionCmd(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &versionOptions{
		commonOptions: commonOptions{
			CommonOptions: common.CommonOptions{Factory: f, Out: out, Err: errOut},
		},
	}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Returns the version information reported by the edge controller ",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	options.AddCommonFlags(cmd)

	return cmd
}

// Run implements this command
func (o *versionOptions) Run() error {
	jsonParsed, err := util.EdgeControllerList("version", nil, o.OutputJSONResponse, o.Out, o.commonOptions.Timeout, o.commonOptions.Verbose)
	if err != nil {
		return err
	}

	if !o.OutputJSONResponse {
		fmt.Printf("Version     : %v\n", jsonParsed.S("data", "version").Data())
		fmt.Printf("GIT revision: %v\n", jsonParsed.S("data", "revision").Data())
		fmt.Printf("Build Date  : %v\n", jsonParsed.S("data", "buildDate").Data())
		fmt.Printf("Runtime     : %v\n", jsonParsed.S("data", "runtimeVersion").Data())
	}

	return nil
}
