package cmd

func init() {
	rootCmd.AddCommand(NewCheckCmd().Command)
	rootCmd.AddCommand(NewInCmd().Command)
}
