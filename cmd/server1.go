/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	server "github.com/ccheeliang/websocket-proxy/pkg/server"
	"github.com/spf13/cobra"
)

// client1Cmd represents the client1 command
var client1Cmd = &cobra.Command{
	Use:   "server1",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting server instance 1")
		server.StartServer("12338", "http://localhost:12440/ws")
	},
}

func init() {
	rootCmd.AddCommand(client1Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// client1Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// client1Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
