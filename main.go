package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/urfave/cli"
)

const (
	// ExitCodeOK exit code case of OK
	ExitCodeOK int = iota //0
	// ExitCodeError exit code case of Error
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
		cli.StringFlag{
			Name:  "separator, s",
			Usage: "separator charcter",
		},
		cli.StringFlag{
			Name:  "limits, l",
			Usage: "separated string limit lenght",
		},
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

	limitsInput := c.String("limits")
	if limitsInput != "" {
		tmp, err := parceLimits(limitsInput)
		if err != nil {
			return err
		}
		fmt.Println("limits ", tmp)
		paddedSheet = trancateLimitedLength(paddedSheet, tmp)
		fmt.Println("limited sheet ", paddedSheet)
	}

	separator := c.String("separator")
	if separator != "" {
		draw(paddedSheet, separator)
	} else {
		draw(paddedSheet, "\t")
	}

	return nil
}

func toSheetString(val string) ([][]string, int) {
	lines := strings.Split(val, "\n")

	var columnSize int
	sheet := [][]string{}
	for _, v := range lines {
		if v == "" {
			continue
		}

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
			runeLen := runewidth.StringWidth(word)
			if counts[i] < runeLen {
				counts[i] = runeLen
			}
		}
	}
	return counts
}

func paddingSheet(sheet [][]string, counts []int) [][]string {
	for i, words := range sheet {
		for j, word := range words {
			sheet[i][j] = padRight(word, counts[j])
		}
	}
	return sheet
}

func padRight(str string, length int) string {
	ws := fmt.Sprintf("%-*s", length-runewidth.StringWidth(str), "")
	return fmt.Sprint(str, ws)
}

func parceLimits(str string) ([]int, error) {
	str = strings.Replace(str, " ", "", -1)
	sp := strings.Split(str, ",")

	limits := make([]int, len(sp))
	for i, limitstr := range sp {
		limit, err := strconv.Atoi(limitstr)
		if err != nil {
			return nil, fmt.Errorf("invalid limit option")
		}
		limits[i] = limit
	}
	return limits, nil
}

func trancateLimitedLength(sheet [][]string, limits []int) [][]string {
	for _, words := range sheet {
		for i, limit := range limits {
			words[i] = cutLeft(words[i], limit)
		}
	}
	return sheet
}

func cutLeft(str string, length int) string {
	return runewidth.Truncate(str, length, "")
}

func draw(sheet [][]string, separator string) {
	for _, words := range sheet {
		fmt.Println(strings.Join(words, separator))
	}
}
