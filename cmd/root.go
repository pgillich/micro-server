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
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pgillich/micro-server/configs"
	"github.com/pgillich/micro-server/internal/buildinfo"
	"github.com/pgillich/micro-server/pkg/logger"
	"github.com/pgillich/micro-server/pkg/model"
	"github.com/pgillich/micro-server/pkg/utils"
)

var cfgFile string //nolint:gochecknoglobals // cobra

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{ //nolint:gochecknoglobals // cobra
	Use:   buildinfo.BuildInfo.AppName(),
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var rootViper = viper.New()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context, args []string, testConfig configs.TestConfig) {
	ctx = context.WithValue(ctx, model.CtxKeyCmd, strings.Join(append([]string{rootCmd.Use}, args...), " "))
	ctx = context.WithValue(ctx, model.CtxKeyTestConfig, &testConfig)
	rootCmd.SetArgs(args)
	rootCmd.SetContext(ctx)
	if err := rootCmd.Execute(); err != nil {
		logger.GetLogger(rootCmd.Use, slog.LevelDebug).Error("EXECUTE_FAILED", logger.KeyError, err, "args", args)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.server_runner.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
// TODO use local viper instance
// TODO do not use global cfgFile (?)
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		rootViper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".server_runner" (without extension).
		rootViper.AddConfigPath(home)
		rootViper.SetConfigType("yaml")
		rootViper.SetConfigName(".server_runner")
	}

	// rootViper.AutomaticEnv() // read in environment variables that match

	log := logger.GetLogger(rootCmd.Use, slog.LevelDebug)
	// If a config file is found, read it in.
	if err := rootViper.ReadInConfig(); err != nil {
		log.Error("EXECUTE_FAILED", logger.KeyError, err, "path", rootViper.ConfigFileUsed())
		os.Exit(1)
	}
	logger.GetLogger(rootCmd.Use, slog.LevelDebug).Info("CONFIG_FILE", "path", rootViper.ConfigFileUsed())
}

func InheritViperConfig(child *viper.Viper) {
	for _, key := range rootViper.AllKeys() {
		if !rootViper.IsSet(key) {
			continue
		}
		value := rootViper.Get(key)
		if !utils.IsEmpty(value) {
			child.Set(key, value)
		}
	}
}
