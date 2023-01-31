package main

import (
	// stdlib imports

	"bufio"
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
	"github.com/theorlandog/easy_gpg/pkg/utils"
)

var (
	cli   = kingpin.New("easy_gpg", "easy_gpg is a tool for configuring gpg keys.")
	cmd   string
	debug = cli.Flag("debug", "Run in debug mode.").Bool()
	log   = logging.Log
)

func init() {
	// init() is a reserved function in go that executes
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

func collectKeyInfo() (name string, keylength string, email string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your name. []:")
	scanner.Scan()
	name = scanner.Text()
	fmt.Println("Enter your keylength. [4096]:")
	scanner.Scan()
	keylength = scanner.Text()
	fmt.Println("Enter your email. []:")
	scanner.Scan()
	email = scanner.Text()
	return
}

func main() {
	log.Infof("easy_gpg VERSION: %s", config.VERSION)

	operating_system := runtime.GOOS

	switch operating_system {
	case "windows":
		utils.WindowsCheckDependencies()
	case "darwin":
		utils.UnixCheckDependencies()
	case "linux":
		utils.UnixCheckDependencies()
	default:
		fmt.Println("easy_gpg only works on Windows, MacOS, and Linux")
	}

	_, _, _ = collectKeyInfo()
}
