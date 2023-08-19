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
package hub

import (
	"github.com/spf13/cobra"
)

// hubCmd represents the hub command
var HubCmd = &cobra.Command{
	Use:   "hub [command] [options]",
	Short: "Operations on the hub",
	Long:  `A hub server accepts http requests, and forwards the requests to the agents.`,
}

func init() {
	// Here you will define your flags and configuration settings.
	HubCmd.PersistentFlags().String("secret", "", "The secret key for the hub communicate with the agent")
	HubCmd.MarkPersistentFlagRequired("secret")
}
