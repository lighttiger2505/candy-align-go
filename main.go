package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

const (
	ExitCodeOK    int = iota //0
	ExitCodeError int = iota //1
)

func main() {
	err := newApp().Run(os.Args)
	var exitCode = ExitCodeOK
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		exitCode = ExitCodeError
	}
	os.Exit(exitCode)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "candy-align"
	app.HelpName = "canal"
	app.Usage = "Sweet text alignment command."
	app.UsageText = "canal [options]"
	app.Version = "0.0.1"
	app.Author = "lighttiger2505"
	app.Email = "lighttiger2505@gmail.com"
	app.Flags = []cli.Flag{
		// 		cli.StringFlag{
		// 			Name:  "suffix, x",
		// 			Usage: "Diary file suffix",
		// 		},
	}
	app.Action = run
	return app
}

func run(c *cli.Context) error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed get stdin")
	}
	fmt.Println(string(b))
	return nil
}
