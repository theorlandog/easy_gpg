package main

import (
	// stdlib imports
	"fmt"
	"os"
	"runtime"
	"strings"

	// External imports
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	// Local imports
	"github.com/theorlandog/easy_gpg/pkg/config"
	"github.com/theorlandog/easy_gpg/pkg/logging"
	"github.com/theorlandog/easy_gpg/pkg/steps"
)

var (
	cli   = kingpin.New("easy_gpg", "easy_gpg is a tool for configuring gpg keys.")
	cmd   string
	debug = cli.Flag("debug", "Run in debug mode.").Bool()
	log   = logging.Log
)

func init() {
	// init() is a reserved function in golang that executes
	// when the package is first imported.
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			split := strings.SplitN(arg, "=", 2)
			split[0] = strings.ReplaceAll(split[0], "_", "-")
			os.Args[i] = strings.Join(split, "=")
		}
	}

	cli.Version("easy_gpg " + config.VERSION)
	cmd = kingpin.MustParse(cli.Parse(os.Args[1:]))

	switch {
	case *debug:
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	log.Infof("easy_gpg VERSION: %s", config.VERSION)

	operating_system := runtime.GOOS

	switch operating_system {
	case "windows":
		steps.WindowsCheckDependencies()
	case "darwin":
		steps.UnixCheckDependencies()
	case "linux":
		steps.UnixCheckDependencies()
	default:
		fmt.Println("easy_gpg only works on Windows, MacOS, and Linux")
	}

	name, keyType, keyLength, email, password, expireDays := steps.CollectKeyInfo()
	configString := steps.GenerateKeyGenConfigString(name, keyType, keyLength, email, password, expireDays)
	print(configString)
}
