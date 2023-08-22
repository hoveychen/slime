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
	"net/http"

	"github.com/hoveychen/slime/pkg/hub"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run --secret <secret> --appPassword <appPassword>",
	Short: "Run a hub server.",
	Long:  `A hub server accepts http requests, and forwards the requests to the agents.`,
	Run: func(cmd *cobra.Command, args []string) {
		secret := viper.GetString("secret")
		if secret == "" {
			logrus.Fatal("secret is required")
		}

		host := viper.GetString("host")
		port := viper.GetInt("port")

		var opts []hub.HubServerOption
		if concurrent := viper.GetInt("concurrent"); concurrent > 0 {
			opts = append(opts, hub.WithConcurrent(concurrent))
		}
		if appPassword := viper.GetString("appPassword"); appPassword != "" {
			opts = append(opts, hub.WithAppPassword(appPassword))
		}

		hub := hub.NewHubServer(secret, opts...)

		addr := fmt.Sprintf("%s:%d", host, port)
		logrus.WithField("addr", addr).Info("Starting hub server")
		err := http.ListenAndServe(addr, hub)
		if err != nil {
			logrus.WithError(err).Error("Hub server terminated")
		}
	},
}

func init() {
	HubCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	runCmd.PersistentFlags().String("appPassword", "", "The password for the application to connect to the hub")
	runCmd.PersistentFlags().Int("port", 8080, "Port to listen on")
	runCmd.PersistentFlags().String("host", "0.0.0.0", "Host to listen on")
	runCmd.PersistentFlags().Int("concurrent", 0, "The number of concurrent requests from the applications")
	viper.BindPFlags(runCmd.PersistentFlags())
}
