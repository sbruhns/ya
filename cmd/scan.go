package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan qr-code",
	Long:  `scan screen for qr-codes and add it to yubikey`,

	Run: func(cmd *cobra.Command, args []string) {
		code, err := scanCode()
		if err != nil || code == "" {
			fmt.Println("Error while scanning barcode: %s", err)
			os.Exit(-1)
		}

		yubioathCmdName := "yubioath-cli"
		if exec.Command("which", yubioathCmdName).Run() != nil {
			yubioathCmdName = "yubioath"
		}

		if exec.Command(yubioathCmdName, "put", code).Run() != nil {
			fmt.Println("error while put code to yubikey")
			os.Exit(-1)
		}
	},
}

func scanCode() (string, error) {
	// TODO Add some error handling
	imgSelCmd := exec.Command("import", ":-")
	scanCmd := exec.Command("zbarimg", ":-", "2> /dev/null")

	r, w := io.Pipe()
	imgSelCmd.Stdout = w
	scanCmd.Stdin = r

	var b2 bytes.Buffer
	scanCmd.Stdout = &b2

	imgSelCmd.Start()
	scanCmd.Start()
	imgSelCmd.Wait()
	w.Close()
	scanCmd.Wait()

	splitedOutput := strings.SplitAfter(b2.String(), "QR-Code:")
	if len(splitedOutput) < 2 {
		return "", fmt.Errorf("No barcode found in scaned image")
	}

	return splitedOutput[1], nil
}

func init() {
	RootCmd.AddCommand(scanCmd)
}
