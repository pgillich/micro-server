/*
Copyright © 2024 Peter Gillich <pgillich@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	srv_configs "github.com/pgillich/micro-server/pkg/configs"
	"github.com/pgillich/micro-server/pkg/logger"
	"github.com/pgillich/micro-server/pkg/model"
	"github.com/pgillich/micro-server/pkg/server"
	"github.com/pgillich/micro-server/pkg/utils"
)

var serverViper = viper.New() //nolint:gochecknoglobals // CMD

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "services",
	Short: "Services",
	Long:  `Runs listed microservices`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SetContext(cmd.Parent().Context())
		_, log := logger.FromContext(cmd.Context())
		log.Info("SERVER_TO_RUN", "command", fmt.Sprintf("%+v", cmd.Context().Value(model.CtxKeyCmd)))

		buildInfo, is := cmd.Context().Value(model.CtxKeyBuildInfo).(model.BuildInfo)
		if !is {
			return srv_configs.ErrFatalServerConfig
		}

		serverConfig, is := cmd.Context().Value(model.CtxKeyServerConfig).(srv_configs.ServerConfiger)
		if !is {
			return srv_configs.ErrFatalServerConfig
		}
		InheritViperConfig(serverViper)
		if err := serverViper.Unmarshal(serverConfig); err != nil {
			return err
		}
		serverConfigStr, err := utils.RenderServerConfig(serverConfig)
		if err != nil {
			return err
		}
		log.Debug("SERVER_CONFIG", "config", serverConfigStr)

		testConfig, is := cmd.Context().Value(model.CtxKeyTestConfig).(srv_configs.TestConfiger)
		if !is {
			return srv_configs.ErrFatalServerConfig
		}
		httpServerRunner := testConfig.GetHttpServerRunner()
		if httpServerRunner == nil {
			httpServerRunner = server.RunHttpServer
		}

		err = server.RunServices(cmd.Context(), buildInfo, args, serverConfig, testConfig)
		time.Sleep(time.Second)

		return err
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("listenaddr", "localhost:9090", "Listen address")
	if err := serverViper.BindPFlags(serverCmd.Flags()); err != nil {
		logger.GetLogger(serverCmd.Use, slog.LevelDebug).Error("Unable to bind flags", logger.KeyError, err)
		panic(err)
	}
	//serverViper.AutomaticEnv()
}
