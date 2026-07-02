module github.com/k8shell-io/k8shell

go 1.25.0

require (
	github.com/fatih/color v1.19.0
	github.com/k8shell-io/common v0.30.3
	github.com/k8shell-io/k8shell-go v0.1.0
	github.com/spf13/cobra v1.8.0
	golang.org/x/term v0.44.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)

// Points at the local k8shell-go checkout with GetUser/UpdateUser/GetSession added.
// Remove this once those changes are pushed and a new k8shell-go version is tagged.
replace github.com/k8shell-io/k8shell-go => ../k8shell-go
