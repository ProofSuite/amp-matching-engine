package server

// var cfgDir string
// var env string

// // rootCmd represents the base command when called without any subcommands
// var rootCmd = &cobra.Command{}

// // Execute adds all child commands to the root command and sets flags appropriately.
// // This is called by main.main(). It only needs to happen once to the rootCmd.
// func Execute() {
// 	if err := rootCmd.Execute(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

// func init() {
// 	// Here you will define your flags and configuration settings.
// 	// Cobra supports persistent flags, which, if defined here,
// 	// will be global for your application.
// 	rootCmd.PersistentFlags().StringVar(&cfgDir, "configDir", "./config", "config directory (default is $PROJECT_PATH/config)")
// 	rootCmd.PersistentFlags().StringVar(&env, "env", "", "Environment to use for deployment(default is '')")

// 	cobra.OnInitialize(initConfig)

// }

// func initConfig() {
// 	if err := app.LoadConfig(cfgDir, env); err != nil {
// 		panic(err)
// 	}

// 	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
// 		panic(err)
// 	}
// }

// serveCmd represents the serve command
// var serveCmd = &cobra.Command{
// 	Use:   "serve",
// 	Short: "Get application up and running",
// 	Long:  `Get application up and running`,
// 	Run:   run,
// }

// func init() {
// 	run()
// }

// cmd *cobra.Command, args []string
