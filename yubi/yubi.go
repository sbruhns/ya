package yubi

import "os/exec"

var YubioathCmdName string

func init() {
	YubioathCmdName = "yubioath-cli"
	if exec.Command("which", YubioathCmdName).Run() != nil {
		YubioathCmdName = "yubioath"
	}
}
