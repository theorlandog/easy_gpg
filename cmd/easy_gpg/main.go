package main

import (
	// stdlib imports
	"fmt"
	"os"
	"strings"

	// External imports
	"gopkg.in/alecthomas/kingpin.v2"
	// Local imports
	"github.com/theorlandog/easy_gpg/pkg/config"
	"github.com/theorlandog/easy_gpg/pkg/utils"
)

var (
	cli   = kingpin.New("easy_gpg", "easy_gpg is a tool for configuring gpg keys.")
	cmd   string
	debug = cli.Flag("debug", "Run in debug mode.").Bool()
)

func init() {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			split := strings.SplitN(arg, "=", 2)
			split[0] = strings.ReplaceAll(split[0], "_", "-")
			os.Args[i] = strings.Join(split, "=")
		}
	}

	cli.Version("easy_gpg " + config.VERSION)
	cmd = kingpin.MustParse(cli.Parse(os.Args[1:]))
}

func main() {
	fmt.Println("Hello World")
	// use an imported constant
	fmt.Println("VERSION:", config.VERSION)
	// use an imported function
	fmt.Println("Answer is", utils.AddNumbers(2, 3))
}
