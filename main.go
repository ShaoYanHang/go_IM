package main

import (
	"IM/router"
	"IM/utils"
	"fmt"
	"bytes"
	"os"
	"os/exec"
)

func main() {

	runCommand()

	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()

	r := router.Router()
	r.Run()
}

func runCommand() {
	cmd := exec.Command("swag", "init")
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out.String())
}
