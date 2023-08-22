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
	"github.com/hoveychen/slime/pkg/agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run --token <token> --hub <hub_address> --upstream <upstream_address>",
	Short: "Run a agent server.",
	Long:  `An agent server is a server that can be used to proxy the traffic to the upstream server.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := cmd.Flag("token").Value.String()
		if token == "" {
			logrus.Fatal("No token is provided")
		}
		hub := cmd.Flag("hub").Value.String()
		if hub == "" {
			logrus.Fatal("No hub address is provided")
		}
		upstreams, _ := cmd.PersistentFlags().GetStringSlice("upstream")
		if len(upstreams) == 0 {
			logrus.Fatal("No upstream address is provided")
		}
		reportHardware, _ := cmd.PersistentFlags().GetBool("reportHardware")

		var opts []agent.AgentServerOption
		if numWorker, err := cmd.Flags().GetInt("numWorker"); err == nil && numWorker > 1 {
			opts = append(opts, agent.WithNumWorker(numWorker))
		}
		if !reportHardware {
			opts = append(opts, agent.WithReportHardware(false))
		}

		grp, ctx := errgroup.WithContext(cmd.Context())
		for _, upstream := range upstreams {
			logrus.WithField("upstream", upstream).Info("Starting agent for upstream")
			agent, err := agent.NewAgentServer(hub, upstream, token, opts...)
			if err != nil {
				panic(err)
			}

			grp.Go(func() error {
				if err := agent.Run(ctx); err != nil {
					logrus.WithField("upstream", upstream).WithError(err).Error("Agent server terminated")
				}
				return err
			})
		}
		grp.Wait()
	},
}

func init() {
	AgentCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	runCmd.PersistentFlags().StringSlice("upstream", nil, "The upstream address")
	runCmd.PersistentFlags().Int("numWorker", 1, "The number of workers to handle the requests")
	runCmd.PersistentFlags().Bool("reportHardware", true, "Report the hardware information to the hub")
}
