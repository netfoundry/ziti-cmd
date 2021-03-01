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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/Jeffail/gabs"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/spf13/cobra"
)

type createConfigOptions struct {
	edgeOptions
	jsonFile string
}

// newCreateConfigCmd creates the 'edge controller create service-policy' command
func newCreateConfigCmd(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &createConfigOptions{
		edgeOptions: edgeOptions{
			CommonOptions: common.CommonOptions{Factory: f, Out: out, Err: errOut},
		},
	}

	cmd := &cobra.Command{
		Use:   "config <name> <type> [JSON configuration data]",
		Short: "creates a config managed by the Ziti Edge Controller",
		Long:  "creates a config managed by the Ziti Edge Controller",
		Args:  cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := runCreateConfig(options)
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringVarP(&options.jsonFile, "json-file", "f", "", "Read config JSON from a file instead of the command line")
	options.AddCommonFlags(cmd)

	return cmd
}

// runCreateConfig create a new config on the Ziti Edge Controller
func runCreateConfig(o *createConfigOptions) error {
	var jsonBytes []byte

	if len(o.Args) == 3 {
		jsonBytes = []byte(o.Args[2])
	}

	if o.Cmd.Flags().Changed("json-file") {
		if len(o.Args) == 3 {
			return errors.New("config json specified both in file and on command line. please pick one")
		}
		var err error
		if jsonBytes, err = ioutil.ReadFile(o.jsonFile); err != nil {
			return fmt.Errorf("failed to read config json file %v: %w", o.jsonFile, err)
		}
	}

	if len(jsonBytes) == 0 {
		return errors.New("no config json specified")
	}

	dataMap := map[string]interface{}{}
	if err := json.Unmarshal(jsonBytes, &dataMap); err != nil {
		fmt.Printf("Attempted to parse: %v\n", string(jsonBytes))
		fmt.Printf("Failing parsing JSON: %+v\n", err)
		return errors.Errorf("unable to parse data as json: %v", err)
	}

	configTypeId, err := mapNameToID("config-types", o.Args[1], o.edgeOptions)
	if err != nil {
		return err
	}

	entityData := gabs.New()
	setJSONValue(entityData, o.Args[0], "name")
	setJSONValue(entityData, configTypeId, "configTypeId")
	setJSONValue(entityData, dataMap, "data")
	result, err := createEntityOfType("configs", entityData.String(), &o.edgeOptions)

	if err != nil {
		panic(err)
	}

	configId := result.S("data", "id").Data()

	if _, err := fmt.Fprintf(o.Out, "%v\n", configId); err != nil {
		panic(err)
	}

	return err
}
