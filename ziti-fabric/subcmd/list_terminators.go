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
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/openziti/fabric/pb/mgmt_pb"
	"github.com/openziti/foundation/channel2"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	listTerminatorsClient = NewMgmtClient(listTerminators)
	listCmd.AddCommand(listTerminators)
}

var listTerminators = &cobra.Command{
	Use:   "terminators",
	Short: "Retrieve terminator definitions",
	Run: func(cmd *cobra.Command, args []string) {
		if ch, err := listTerminatorsClient.Connect(); err == nil {
			query := "true limit none"
			if len(args) > 0 {
				query = args[0]
			}
			request := &mgmt_pb.ListTerminatorsRequest{
				Query: query,
			}
			body, err := proto.Marshal(request)
			if err != nil {
				panic(err)
			}
			requestMsg := channel2.NewMessage(int32(mgmt_pb.ContentType_ListTerminatorsRequestType), body)
			responseMsg, err := ch.SendAndWaitWithTimeout(requestMsg, 5*time.Second)
			if err != nil {
				panic(err)
			}
			if responseMsg.ContentType == int32(mgmt_pb.ContentType_ListTerminatorsResponseType) {
				response := &mgmt_pb.ListTerminatorsResponse{}
				if err := proto.Unmarshal(responseMsg.Body, response); err == nil {
					out := fmt.Sprintf("\nTerminators: (%d)\n\n", len(response.Terminators))
					out += fmt.Sprintf("%-10s | %-12s | %-16s | %-12s | %-12s | %-12s | %s\n", "Id", "Service", "Binding", "Static Cost", "Precedence", "Identity", "Destination")
					for _, terminator := range response.Terminators {
						out += fmt.Sprintf("%-10v | %-12v | %-16v | %-12v | %-12v | %-12s | %s\n", terminator.Id, terminator.ServiceId, terminator.Binding,
							terminator.Cost, terminator.Precedence, terminator.Identity, fmt.Sprintf("%-12s -> %s", terminator.RouterId, terminator.Address))
					}
					out += "\n"
					fmt.Print(out)
				} else {
					panic(err)
				}
			} else {
				panic(errors.New("unexpected response"))
			}
		} else {
			panic(err)
		}
	},
}
var listTerminatorsClient *mgmtClient
