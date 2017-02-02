package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sbruhns/ya/yubi"
	"github.com/spf13/cobra"
)

type Key struct {
	Number int
	Name   string
	Token  string
}

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy key to current clipboard",
	Long:  "Copy key to current clipboard",

	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		out, _ := exec.Command(yubi.YubioathCmdName).Output()
		rows := strings.Split(string(out), "\n")
		keys := []Key{}
		number := 0
		for _, value := range rows {
			number++
			parsedKey := strings.Split(value, " ")
			l := len(parsedKey)
			if l > 1 {
				key := Key{number, parsedKey[0], parsedKey[l-1]}
				keys = append(keys, key)
				fmt.Println(strconv.Itoa(key.Number) + ". " + key.Name + " " + key.Token)
			}
		}
		if len(keys) < 1 {
			fmt.Println("No key was found!")
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the number to copy: ")
		text, _ := reader.ReadString('\n')
		selected, err := strconv.Atoi(strings.Trim(text, "\n"))
		if err != nil || selected > len(keys) {
			fmt.Println("Invalid input")
			os.Exit(1)
		}
		token := keys[selected-1].Token
		cpCmd := exec.Command("xsel", "-b", "-i")

		cpCmd.Stdin = strings.NewReader(token)
		err = cpCmd.Run()
		if err != nil {
			fmt.Println("Could not add Token to clipboard")
			os.Exit(1)
		}

		fmt.Println(keys[selected-1].Name + " token was copied!")
	},
}

func init() {
	RootCmd.AddCommand(copyCmd)
}
