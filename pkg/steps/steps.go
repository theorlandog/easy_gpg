package steps

import (
	// stdlib imports
	"os/exec"
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
