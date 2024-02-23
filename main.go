/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/ccheeliang/websocket-proxy/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cmd.Execute()
}
