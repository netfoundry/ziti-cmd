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
	"github.com/openziti/foundation/agent"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// PsStackOptions the options for the create spring command
type PsStackOptions struct {
	PsOptions
	CtrlListener string
}

// NewCmdPsStack creates a command object for the "create" command
func NewCmdPsStack(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &PsStackOptions{
		PsOptions: PsOptions{
			CommonOptions: CommonOptions{
				Factory: f,
				Out:     out,
				Err:     errOut,
			},
		},
	}

	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		Use:  "stack <optional-target>",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdhelper.CheckErr(err)
		},
	}

	options.addCommonFlags(cmd)

	return cmd
}

// Run implements the command
func (o *PsStackOptions) Run() error {
	addr, err := agent.ParseGopsAddress(o.Args)
	if err != nil {
		return err
	}
	return agent.MakeRequest(addr, agent.StackTrace, nil, os.Stdout)
}
