package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
	sheet, columnSize := toSheetString(string(b))
	counts := countColumn(sheet, columnSize)
	paddedSheet := paddingSheet(sheet, counts)
	draw(paddedSheet)

	return nil
}

func toSheetString(val string) ([][]string, int) {
	lines := strings.Split(val, "\n")

	var columnSize int
	sheet := [][]string{}
	for _, v := range lines {
		columns := strings.Fields(v)

		tmpSize := len(columns)
		if columnSize < tmpSize {
			columnSize = tmpSize
		}

		sheet = append(sheet, columns)
	}
	return sheet, columnSize
}

func countColumn(sheet [][]string, columnSize int) []int {
	counts := make([]int, columnSize)
	for _, words := range sheet {
		for i, word := range words {
			wordlen := len(word)
			if counts[i] < wordlen {
				counts[i] = wordlen
			}
		}
	}
	return counts
}

func paddingSheet(sheet [][]string, counts []int) [][]string {
	for i, words := range sheet {
		for j, word := range words {
			sheet[i][j] = padRight(word, counts[j], " ")
		}
	}
	return sheet
}

func padRight(str string, length int, padChar string) string {
	return str + times(padChar, length-len(str))
}

func times(str string, n int) (out string) {
	for i := 0; i < n; i++ {
		out += str
	}
	return
}

func draw(sheet [][]string) {
	for _, words := range sheet {
		fmt.Println(strings.Join(words, " "))
	}
}
