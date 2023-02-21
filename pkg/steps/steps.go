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

func GenerateKeyGenConfigString(keyGenConfigParams KeyGenConfigParams) (configstring string) {
	var buf bytes.Buffer

	configTemplate := template.New("configTemplate")
	configTemplate = template.Must(configTemplate.Parse(config.GENPARAMSTEMPLATE))

	// template.Execute writes to an io.writer object
	configTemplate.Execute(&buf, keyGenConfigParams)

	configstring = buf.String()

	return
}

func ExtractFingerprintFromKeyList(keyList []string) (fingerprint string) {
	fingerprint = ""

	// using for loop, extract the finger print.
	// example string: `fpr:::::::::577963869F20A39E58065014A24551D20C9D6802:`
	for i := 0; i < len(keyList); i++ {
		if strings.HasPrefix(keyList[i], "fpr") {
			fingerprint = strings.Split(keyList[i], ":")[9]
			break
		}
	}

	// log if no fingerprint detected
	if strings.TrimSpace(fingerprint) == "" {
		log.Debug("No fingerprint detected.")
	}

	return
}

func GenerateKeys(keyGenConfigParams KeyGenConfigParams) int {
	expire := strconv.FormatInt(int64(keyGenConfigParams.ExpireDays), 10)
	if expire == "0" {
		expire = "none"
	}

	// gpg --quick-gen-key --passphrase "" --batch "jack mehoff" RSA4096 cert,sign none
	cmd := exec.Command("gpg", "--quick-gen-key", "--passphrase"+"="+keyGenConfigParams.Password, "--batch", keyGenConfigParams.Name, keyGenConfigParams.KeyType+keyGenConfigParams.KeyLength, "cert,sign", expire)
	_, err := cmd.Output()
	if err != nil {
		log.Info("gpg failure.")
		fmt.Println(cmd)
		// return 1
	}

	// Prepend an '=' to the key users name to only list exact matches to the name
	cmd = exec.Command("gpg", "--list-keys", "--with-colons", "--batch", "="+keyGenConfigParams.Name)
	output, err := cmd.Output()
	if err != nil {
		log.Info("gpg failure.")
		fmt.Println(cmd)
		// return 1
	}

	// Turn the output into an array
	keyList := strings.Split(string(output), "\n")
	fingerprint := ExtractFingerprintFromKeyList(keyList)

	// Following ideas from https://serverfault.com/questions/818289/add-second-sub-key-to-unattended-gpg-key#962553
	cmd = exec.Command("gpg", "--quick-add-key", "--no-tty", "--batch", "--passphrase"+"="+keyGenConfigParams.Password, fingerprint, keyGenConfigParams.KeyType+keyGenConfigParams.KeyLength, "auth", expire)
	_, err = cmd.Output()
	if err != nil {
		print("auth\n")
		fmt.Println(cmd)
		log.Info(err.Error())
		// return 1
	}

	cmd = exec.Command("gpg", "--quick-add-key", "--no-tty", "--batch", "--passphrase"+"="+keyGenConfigParams.Password, fingerprint, keyGenConfigParams.KeyType+keyGenConfigParams.KeyLength, "sign", expire)
	_, err = cmd.Output()
	if err != nil {
		print("sign\n")
		fmt.Println(cmd)
		log.Info(err.Error())
		// return 1
	}

	cmd = exec.Command("gpg", "--quick-add-key", "--no-tty", "--batch", "--passphrase"+"="+keyGenConfigParams.Password, fingerprint, keyGenConfigParams.KeyType+keyGenConfigParams.KeyLength, "encrypt", expire)
	_, err = cmd.Output()
	if err != nil {
		print("encrypt\n")
		fmt.Println(cmd)
		log.Info(err.Error())
		// return 1
	}

	return 0
}
