package main

import (
	"log"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/shi0rik0/docker-ebpf-plugin/driver"
)

const PLUGIN_NAME = "ebpf"

func main() {
	driver := driver.NewDriver()
	handler := network.NewHandler(driver)
	err := handler.ServeUnix(PLUGIN_NAME, 0)
	if err != nil {
		log.Fatal(err)
	}
}
