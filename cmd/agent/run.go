/*
Copyright © 2023 Harry C <hoveychen@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package agent

import (
	"fmt"

	"github.com/hoveychen/slime/pkg/agent"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run --token <token> --hub <hub_address> --upstream <upstream_address>",
	Short: "Run a agent server.",
	Long:  `An agent server is a server that can be used to proxy the traffic to the upstream server.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := cmd.Flag("token").Value.String()
		hub := cmd.Flag("hub").Value.String()
		upstream := cmd.Flag("upstream").Value.String()

		var opts []agent.AgentServerOption
		if numWorker, err := cmd.Flags().GetInt("numWorker"); err == nil && numWorker > 1 {
			opts = append(opts, agent.WithNumWorker(numWorker))
		}

		fmt.Printf("Starting agent")
		agent, err := agent.NewAgentServer(hub, upstream, token, opts...)
		if err != nil {
			panic(err)
		}
		if err := agent.Run(cmd.Context()); err != nil {
			panic(err)
		}
	},
}

func init() {
	AgentCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	runCmd.PersistentFlags().String("upstream", "", "The upstream address")
	runCmd.MarkPersistentFlagRequired("upstream")
	runCmd.PersistentFlags().Int("numWorker", 1, "The number of workers to handle the requests")
}
