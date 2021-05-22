package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sreesanthv/go-api-base/api"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := api.NewServer()
		if err != nil {
			log.Fatal(err)
		}
		s.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// flags configuration
	// TODO FLAGS - read from env
	serveCmd.Flags().StringP("port", "p", "localhost:3000", "Port to run Application server on")
	serveCmd.Flags().StringP("log_level", "l", "debug", "Log level")
	serveCmd.Flags().BoolP("allow_cors", "c", false, "Log level")

	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("log_level", serveCmd.Flags().Lookup("log_level"))
	viper.BindPFlag("allow_cors", serveCmd.Flags().Lookup("allow_cors"))
}
