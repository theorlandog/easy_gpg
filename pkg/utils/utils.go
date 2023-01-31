package utils

import (
	"os/exec"

	"github.com/theorlandog/easy_gpg/pkg/logging"
)

var (
	log = logging.Log
)

func UnixCheckDependencies() int {
	cmd := exec.Command("which", "gpg")

	_, err := cmd.Output()
	if err != nil {
		log.Info("gpg binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	log.Debug("Found gpg binary.")

	cmd = exec.Command("which", "gpg-agent")

	_, err = cmd.Output()
	if err != nil {
		log.Info("gpg-agent binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	log.Debug("Found gpg-agent binary.")

	return 0
}

func WindowsCheckDependencies() int {
	cmd := exec.Command("where", "gpg")

	_, err := cmd.Output()
	if err != nil {
		log.Info("gpg binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	log.Debug("Found gpg binary.")

	cmd = exec.Command("where", "gpg-agent")

	_, err = cmd.Output()
	if err != nil {
		log.Info("gpg-agent binary not found in path. Please install from https://gnupg.org/.")
		return 1
	}
	log.Debug("Found gpg-agent binary.")

	return 0
}
