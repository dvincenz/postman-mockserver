package cmd

import (
	"github.com/dvincenz/postman-mockserver/postman"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var token string

// onlineCmd represents the online command
var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "start a postman mock server emulator",
	Long: `This application will start a postman mock server emulation. 
	It will serve all mock endpoints defined in postman.
	For this reason at least an postman api key or an postman collection export is needed to access the mocks.`,
	Run: func(cmd *cobra.Command, args []string) {
		postman.StartServer()
	},
}



func init() {
	rootCmd.AddCommand(onlineCmd)
	onlineCmd.Flags().StringVarP(&token, "token", "t", "", "Add Postman token, check readme for more instructions")
	viper.BindPFlag("postman.token", onlineCmd.Flags().Lookup("token"))

	vi := viper.GetViper()
	log.Trace().Msg(vi.GetString("postman.token"))

}
