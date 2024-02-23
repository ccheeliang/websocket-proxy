/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	websocketserver "github.com/ccheeliang/websocket-proxy/pkg/websocket_server"
	"github.com/spf13/cobra"
)

// websocketCmd represents the websocket command
var websocketCmd = &cobra.Command{
	Use:   "websocket",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("websocket called")
		websocketserver.ListenAndRunWebsocket("12440")
	},
}

func init() {
	rootCmd.AddCommand(websocketCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// websocketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// websocketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
