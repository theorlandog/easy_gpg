package steps

import (
	// stdlib imports
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"text/template"

	// Local imports
	"github.com/theorlandog/easy_gpg/pkg/config"
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

func CollectKeyInfo() (name string, keyType string, keylength string, email string, password string, expireDays int) {
	scanner := bufio.NewScanner(os.Stdin)

	// Get current user for default name
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

	fmt.Println("Enter your keytype. [RSA]:")
	scanner.Scan()
	keyType = strings.ToUpper(strings.TrimSpace(scanner.Text()))
	if keyType == "" {
		// Set default keylength
		keyType = "RSA"
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

	fmt.Println("Enter a password for you key. []:")
	scanner.Scan()
	password = strings.TrimSpace(scanner.Text())

	fmt.Println("How many days until the master key should expire? (0 is never). [0]:")
	scanner.Scan()
	expireDaysText := strings.TrimSpace(scanner.Text())
	if expireDaysText == "" {
		expireDays = 0
	} else {
		expireDays, err = strconv.Atoi(expireDaysText)
		if err != nil {
			// ... handle error
			log.Fatalf(err.Error())
		}
	}

	return
}

type KeyGenConfigParams struct {
	Name       string
	Email      string
	KeyType    string
	KeyLength  string
	Password   string
	ExpireDays int
}

func GenerateKeyGenConfigString(name string, keyType string, keyLength string, email string, password string, expireDays int) (configstring string) {
	var buf bytes.Buffer
	keyGenConfigParams := KeyGenConfigParams{
		Name:       name,
		Email:      email,
		KeyType:    keyType,
		KeyLength:  keyLength,
		Password:   password,
		ExpireDays: expireDays,
	}
	configTemplate := template.New("configTemplate")
	configTemplate = template.Must(configTemplate.Parse(config.GENPARAMSTEMPLATE))

	// template.Execute writes to an io.writer object
	configTemplate.Execute(&buf, keyGenConfigParams)

	configstring = buf.String()

	return
}
