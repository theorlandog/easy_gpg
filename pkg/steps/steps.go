package steps

import (
	// stdlib imports
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	// Local imports
	"github.com/theorlandog/easy_gpg/pkg/logging"
)

var (
	log = logging.Log
)

func UnixCheckDependencies() int {
	cmd := exec.Command("which", "gpg")

	output, err := cmd.Output()
	if err != nil {
		log.Info("gpg binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	// Don't forget to remove the newline suffix of output.
	log.Debugf("Found gpg binary at %s.", strings.Split(string(output), "\n")[0])

	cmd = exec.Command("which", "gpg-agent")

	output, err = cmd.Output()
	if err != nil {
		log.Info("gpg-agent binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	log.Debugf("Found gpg-agent binary at %s.", strings.Split(string(output), "\n")[0])

	return 0
}

func WindowsCheckDependencies() int {
	cmd := exec.Command("where", "gpg")

	output, err := cmd.Output()
	if err != nil {
		log.Info("gpg binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	// Don't forget to remove the newline suffix of output.
	log.Debugf("Found gpg binary at %s.", strings.Split(string(output), "\n")[0])

	cmd = exec.Command("where", "gpg-agent")

	output, err = cmd.Output()
	if err != nil {
		log.Info("gpg-agent binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	log.Debugf("Found gpg-agent binary at %s.", strings.Split(string(output), "\n")[0])

	return 0
}

func CollectKeyInfo() (name string, keylength string, email string) {
	scanner := bufio.NewScanner(os.Stdin)
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Enter your name. [%s]:\n", currentUser.Name)
	scanner.Scan()
	name = strings.TrimSpace(scanner.Text())
	if name == "" {
		// Set the default to the OS username
		name = currentUser.Name
	}
	fmt.Println("Enter your keylength. [4096]:")
	scanner.Scan()
	keylength = strings.TrimSpace(scanner.Text())
	if keylength == "" {
		// Set default keylength
		keylength = "4096"
	}
	fmt.Println("Enter your email. []:")
	scanner.Scan()
	email = strings.TrimSpace(scanner.Text())
	return
}
