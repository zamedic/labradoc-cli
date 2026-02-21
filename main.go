package main

import (
	"github.com/zamedic/labrador-cli/cmd"
	"github.com/zamedic/labrador-cli/internal/config"
)

func main() {
	config.InitConfig()
	cmd.Execute()
}
