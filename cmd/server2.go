/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	server "github.com/ccheeliang/websocket-proxy/pkg/server"
	"github.com/spf13/cobra"
)

// client2Cmd represents the client2 command
var client2Cmd = &cobra.Command{
	Use:   "server2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting server instance 2")
		server.StartServer("12339", "http://localhost:12440/ws")
	},
}

func init() {
	rootCmd.AddCommand(client2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// client2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// client2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
