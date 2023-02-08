package config

const VERSION = "0.1"

const GENPARAMSTEMPLATE = `%echo Generating a OpenPGP key

{{if eq .Password ""}}
%no-protection
{{else}}
Passphrase: {{.Password}}
{{end}}

Key-Type: {{.KeyType}}
Key-Length: {{.KeyLength}}
# Key generated is a master key ("certificate")
Key-Usage: cert

# Parameters to generate subkeys
Subkey-Type: ELG-E
Subkey-Length: {{.KeyLength}}
Subkey-Usage: encrypt

Subkey-Type: ELG-S
Subkey-Length: {{.KeyLength}}
Subkey-Usage: sign

Subkey-Type: ELG-A
Subkey-Length: {{.KeyLength}}
Subkey-Usage: auth

# select a name and email address - neither has to be valid nor existing
Name-Real: {{.Name}}
Name-Email: {{.Email}}

# Set the key to expiration (0 is never)
Expire-Date: {{.ExpireDays}}

# Do a commit here, so that we can later print "done" :-)
%commit

%echo done
`
