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
	"github.com/Jeffail/gabs"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
)

type createCaOptions struct {
	edgeOptions
	tags             map[string]string
	name             string
	caPath           string
	caPemBytes       []byte
	autoCaEnrollment bool
	ottCaEnrollment  bool
	authEnabled      bool
	identityRoles    []string
}

// newCreateCaCmd creates the 'edge controller create ca local' command for the given entity type
func newCreateCaCmd(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &createCaOptions{
		edgeOptions: edgeOptions{
			CommonOptions: common.CommonOptions{
				Factory: f,
				Out:     out,
				Err:     errOut,
			},
		},
		tags: make(map[string]string),
	}

	cmd := &cobra.Command{
		Use:   "ca <name> <pemCertFile> [--autoca, --ottca, --auth]",
		Short: "creates a ca managed by the Ziti Edge Controller",
		Long:  "creates a ca managed by the Ziti Edge Controller",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("requires at least %d arg(s), only received %d", 2, len(args))
			}

			options.name = args[0]
			options.caPath = args[1]

			var err error
			options.caPemBytes, err = ioutil.ReadFile(options.caPath)

			if err != nil {
				return fmt.Errorf("could not read CA certificate file: %s", err)
			}

			return nil

		},
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := runCreateCa(options)
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringToStringVarP(&options.tags, "tags", "t", nil, "Add tags to service definition")
	cmd.Flags().BoolVarP(&options.authEnabled, "auth", "e", false, "Whether the CA can be used for authentication or not")
	cmd.Flags().BoolVarP(&options.ottCaEnrollment, "ottca", "o", false, "Whether the CA can be used for one-time-token CA enrollment")
	cmd.Flags().BoolVarP(&options.autoCaEnrollment, "autoca", "u", false, "Whether the CA can be used for auto CA enrollment")
	cmd.Flags().StringSliceVarP(&options.identityRoles, "role-attributes", "a", []string{}, "A csv string of role attributes enrolling identities receive")
	options.AddCommonFlags(cmd)

	return cmd
}

func runCreateCa(options *createCaOptions) (err error) {
	data := gabs.New()
	setJSONValue(data, options.name, "name")
	setJSONValue(data, options.autoCaEnrollment, "isAutoCaEnrollmentEnabled")
	setJSONValue(data, options.ottCaEnrollment, "isOttCaEnrollmentEnabled")
	setJSONValue(data, options.authEnabled, "isAuthEnabled")
	setJSONValue(data, string(options.caPemBytes), "certPem")
	setJSONValue(data, options.identityRoles, "identityRoles")

	result, err := createEntityOfType("cas", data.String(), &options.edgeOptions)

	if err != nil {
		panic(err)
	}

	id := result.S("data", "id").Data()

	if _, err = fmt.Fprintf(options.Out, "%v\n", id); err != nil {
		panic(err)
	}

	return err
}
