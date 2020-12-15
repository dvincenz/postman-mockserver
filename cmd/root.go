/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"github.com/dvincenz/postman-mockserver/postman"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var cfgFile string
var port int

var rootCmd = &cobra.Command{
	Use:   "postman-mockserver",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("mode") == "online" {
			postman.StartServer()
		}
		postman.StartServerFromStaticFile()
	},
}



func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.pms/config.yaml)")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P", 8080, "starting port")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

}



func initConfig() {
	cfgFile = viper.GetString("config")
	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			log.Error().Msg("config file " + cfgFile + " does not exists")
		}else{
			log.Info().Msg("Config file: " + viper.GetString("config"))

		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName("/.pms/config")


	}
	viper.SetEnvPrefix("pms")
	viper.AutomaticEnv() // read in environment variables that match
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)


	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msg("error reading config " + err.Error())
	}
	//todo log need to be setup before logging config-logs but configuration is needed for setup logger :-/
	initLogger()
}

func initLogger() {
	if !viper.GetBool("logging.jsonLogging"){
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	loglevel, err := zerolog.ParseLevel(viper.GetString("logging.level"))
	if err != nil || loglevel == zerolog.NoLevel{
		log.Warn().Msg("loglevel not set, default level is set to trace")
		loglevel = zerolog.TraceLevel
	}
	zerolog.SetGlobalLevel(loglevel)
}
