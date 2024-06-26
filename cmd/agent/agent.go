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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// agentCmd represents the agent command
var AgentCmd = &cobra.Command{
	Use:   "agent [command] [options]",
	Short: "Operations on the agent",
	Long:  `An agent server is a server that can be used to proxy the traffic to the upstream server.`,
}

func init() {

	// Here you will define your flags and configuration settings.

	AgentCmd.PersistentFlags().String("token", "", "The agent token for the agent to communicate with the hub")
	AgentCmd.PersistentFlags().String("hub", "", "The hub address")
	AgentCmd.PersistentFlags().Int("agentID", 0, "Override the agent ID")
	viper.BindPFlags(AgentCmd.PersistentFlags())
}
