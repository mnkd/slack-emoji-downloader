package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
	ExitCodeFileError
)

var (
	Version  string
	Revision string
)

func usage() {
	str := `
USAGE
  slack-emoji-downloader --token your_api_token

SEE ALSO
   https://api.slack.com/methods/emoji.list

`
	fmt.Fprintln(os.Stderr, str)
}

var app *App

func init() {
	var token string
	var ver bool

	flag.StringVar(&token, "token", "", "Authentication token.")
	flag.BoolVar(&ver, "v", false, "Print version.")
	flag.Parse()

	if ver {
		fmt.Fprintln(os.Stdout, "Version:", Version)
		fmt.Fprintln(os.Stdout, "Revision:", Revision)
		os.Exit(ExitCodeOK)
	}

	if len(token) == 0 {
		usage()
		os.Exit(ExitCodeOK)
	}

	app = NewApp(token)
}

func main() {
	os.Exit(app.Run())
}
