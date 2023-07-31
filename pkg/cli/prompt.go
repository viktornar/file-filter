package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Yes string = "yes"
	No  string = "no"
)

func Prompt(text string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("- %s: \n", text)
		if answer, _ := reader.ReadString('\n'); answer != "" && 
			(strings.TrimSpace(strings.ToLower(answer)) == Yes || strings.TrimSpace(strings.ToLower(answer)) == No) {
			return strings.TrimSpace(answer)
		}
	}
}
