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
			Name:  "input-delimiter, s",
			Usage: "Delimiter of input string. Use for splitting",
		},
		cli.StringFlag{
			Name:  "output-delimiter, d",
			Usage: "Delimiter of output string",
		},
		cli.StringFlag{
			Name:  "width, w",
			Usage: "Output column width, separate columns by commas",
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

	inputDelimiter := c.String("input-delimiter")
	sheet, columnSize := toSheetString(string(b), inputDelimiter)

	counts := countColumn(sheet, columnSize)
	paddedSheet := paddingSheet(sheet, counts)

	widthInput := c.String("width")
	if widthInput != "" {
		width, err := parceLimits(widthInput)
		if err != nil {
			return err
		}
		paddedSheet = trancateLimitedLength(paddedSheet, width)
	}

	lines := createDrawLines(sheet, c.String("output-delimiter"))

	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
}

func toSheetString(str, delimiter string) ([][]string, int) {
	lines := strings.Split(str, "\n")

	var columnSize int
	sheet := [][]string{}
	for _, v := range lines {
		if v == "" {
			continue
		}

		var columns []string
		if delimiter != "" {
			columns = strings.Split(v, delimiter)
			for i, column := range columns {
				columns[i] = strings.TrimSpace(column)
			}
		} else {
			columns = strings.Fields(v)
		}

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

func createDrawLines(sheet [][]string, delimiter string) []string {
	if delimiter == "" {
		delimiter = "\t"
	}

	res := make([]string, len(sheet))
	for i, fields := range sheet {
		res[i] = strings.Join(fields, delimiter)
	}
	return res
}
