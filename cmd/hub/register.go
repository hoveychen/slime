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
	"fmt"
	"math/rand"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hoveychen/slime/pkg/token"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register --secret <secret>",
	Short: "Register an agent, and get an agent token for the agent to communicate with the hub",
	Long:  `The agent token is the only way to authenticate an agent to the hub.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flag("name").Value.String()
		age, _ := cmd.Flags().GetDuration("age")
		scopePaths, _ := cmd.Flags().GetStringSlice("scopePaths")
		secret := cmd.Flag("secret").Value.String()
		if secret == "" {
			logrus.Fatal("The secret is required")
		}

		if name == "" {
			// generate a random name
			name = petname.Generate(2, "-")
		}

		agentToken := token.AgentToken{
			Id:         rand.Int63(),
			Name:       name,
			ScopePaths: scopePaths,
		}
		if age > 0 {
			agentToken.ExpireAt = time.Now().Add(age).Unix()
		}

		tokenMgr := token.NewTokenManager([]byte(secret))
		data, err := tokenMgr.Encrypt(&agentToken)
		if err != nil {
			logrus.WithError(err).Error("Failed to encrypt the agent token")
			return
		}

		fmt.Println(string(data))
	},
}

func init() {
	HubCmd.AddCommand(registerCmd)

	registerCmd.PersistentFlags().String("name", "", "The agent name")
	registerCmd.PersistentFlags().Duration("age", 0, "When specified, the token will be expired after the specified age. format like '1h2m3s'")
	registerCmd.PersistentFlags().StringSlice("scopePaths", []string{}, "When specified, the agent accepts only the scoped paths")
}
