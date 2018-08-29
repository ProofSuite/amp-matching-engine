package cmd

import (
	"fmt"
	"os"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/spf13/cobra"
)

var cfgDir string
var env string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "amp-matching-engine",
	Short: "Decentralised cryptocurrency matching engine",
	Long:  `Decentralised cryptocurrency matching engine`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgDir, "configDir", "./config", "config directory (default is $PROJECT_PATH/config)")
	rootCmd.PersistentFlags().StringVar(&env, "env", "", "Environment to use for deployment(default is '')")

	cobra.OnInitialize(initConfig)

}

func initConfig() {
	if err := app.LoadConfig(cfgDir, env); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}
}
