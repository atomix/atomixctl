// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import "github.com/spf13/cobra"

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {databases, primitives}",
		Short: "List resources in the cluster",
	}
	cmd.AddCommand(newGetDatabasesCommand())
	cmd.AddCommand(newGetPrimitivesCommand())
	cmd.AddCommand(newGetCountersCommand())
	cmd.AddCommand(newGetElectionsCommand())
	cmd.AddCommand(newGetIndexedMapsCommand())
	cmd.AddCommand(newGetLeaderLatchesCommand())
	cmd.AddCommand(newGetListsCommand())
	cmd.AddCommand(newGetLocksCommand())
	cmd.AddCommand(newGetLogsCommand())
	cmd.AddCommand(newGetMapsCommand())
	cmd.AddCommand(newGetSetsCommand())
	cmd.AddCommand(newGetValuesCommand())
	addClientFlags(cmd)
	return cmd
}
