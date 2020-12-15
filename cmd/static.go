
package cmd

import (
	"github.com/dvincenz/postman-mockserver/postman"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var collectionFilePath string

var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "run PMS from postman collection export file",
	Long: `You can run PMS from a exported collection file.
	The program will read all the mocks in the definition and serve this mocks as a webservice`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("start mock from static file")
		postman.StartServerFromStaticFile()
	},
}


func init() {
	rootCmd.AddCommand(staticCmd)
	staticCmd.Flags().StringVarP(&collectionFilePath, "path", "p", "./collection.json", "path to postman collection json file")
	viper.BindPFlag("static.path", staticCmd.Flags().Lookup("path"))
}
