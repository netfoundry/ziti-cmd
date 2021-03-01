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
	"io"

	"github.com/pkg/errors"

	"github.com/Jeffail/gabs"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/spf13/cobra"
)

type updateIdentityOptions struct {
	edgeOptions
	name           string
	roleAttributes []string
}

func newUpdateIdentityCmd(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &updateIdentityOptions{
		edgeOptions: edgeOptions{
			CommonOptions: common.CommonOptions{Factory: f, Out: out, Err: errOut},
		},
	}

	cmd := &cobra.Command{
		Use:   "identity <idOrName>",
		Short: "updates a identity managed by the Ziti Edge Controller",
		Long:  "updates a identity managed by the Ziti Edge Controller",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := runUpdateIdentity(options)
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	options.AddCommonFlags(cmd)
	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringVarP(&options.name, "name", "n", "", "Set the name of the identity")
	cmd.Flags().StringSliceVarP(&options.roleAttributes, "role-attributes", "a", nil,
		"Set role attributes of the identity. Use --role-attributes '' to set an empty list")

	return cmd
}

// runUpdateIdentity update a new identity on the Ziti Edge Controller
func runUpdateIdentity(o *updateIdentityOptions) error {
	id, err := mapNameToID("identities", o.Args[0], o.edgeOptions)
	if err != nil {
		return err
	}
	entityData := gabs.New()
	change := false

	if o.Cmd.Flags().Changed("name") {
		setJSONValue(entityData, o.name, "name")
		change = true
	}

	if o.Cmd.Flags().Changed("role-attributes") {
		setJSONValue(entityData, o.roleAttributes, "roleAttributes")
		change = true
	}

	if !change {
		return errors.New("no change specified. must specify at least one attribute to change")
	}

	_, err = patchEntityOfType(fmt.Sprintf("identities/%v", id), entityData.String(), &o.edgeOptions)
	return err
}
