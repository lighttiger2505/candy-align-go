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
			Usage: "Output column width, separate fields by commas",
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

	var width []int
	widthInput := c.String("width")
	if widthInput != "" {
		width, err = parceWidthFlag(widthInput)
		if err != nil {
			return err
		}
	}

	lines := Align(
		string(b),
		&AlignOption{
			inputDelimiter:  c.String("input-delimiter"),
			outputDelimiter: c.String("output-delimiter"),
			width:           width,
		},
	)

	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
}

type AlignOption struct {
	inputDelimiter  string
	outputDelimiter string
	width           []int
}

func Align(text string, opt *AlignOption) []string {
	table, count := Separate(text, opt.inputDelimiter)
	return Format(table, count, opt.width, opt.inputDelimiter)
}

func Separate(str, delimiter string) ([][]string, []int) {
	table, fieldNum := splitToTable(str, delimiter)
	counts := countFields(table, fieldNum)
	return table, counts
}

func Format(table [][]string, counts, width []int, delimiter string) []string {
	table = padFields(table, counts)
	if len(width) > 0 {
		table = trancateProtrudeString(table, width)
	}
	return createDrawLines(table, delimiter)
}

func splitToTable(str, delimiter string) ([][]string, int) {
	lines := strings.Split(str, "\n")

	var fieldNum int
	table := [][]string{}
	for _, v := range lines {
		if v == "" {
			continue
		}

		var fields []string
		if delimiter != "" {
			fields = strings.Split(v, delimiter)
			for i, column := range fields {
				fields[i] = strings.TrimSpace(column)
			}
		} else {
			fields = strings.Fields(v)
		}

		tmpSize := len(fields)
		if fieldNum < tmpSize {
			fieldNum = tmpSize
		}
		table = append(table, fields)
	}
	return table, fieldNum
}

func countFields(table [][]string, fieldNum int) []int {
	counts := make([]int, fieldNum)
	for _, fields := range table {
		for i, field := range fields {
			runeLen := runewidth.StringWidth(field)
			if counts[i] < runeLen {
				counts[i] = runeLen
			}
		}
	}
	return counts
}

func padFields(table [][]string, counts []int) [][]string {
	for i, fields := range table {
		for j, field := range fields {
			table[i][j] = padRight(field, counts[j])
		}
	}
	return table
}

func padRight(str string, length int) string {
	ws := fmt.Sprintf("%-*s", length-runewidth.StringWidth(str), "")
	return fmt.Sprint(str, ws)
}

func parceWidthFlag(str string) ([]int, error) {
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

func trancateProtrudeString(table [][]string, limits []int) [][]string {
	for _, fields := range table {
		for i, limit := range limits {
			fields[i] = runewidth.Truncate(fields[i], limit, "")
		}
	}
	return table
}

func createDrawLines(table [][]string, delimiter string) []string {
	if delimiter == "" {
		delimiter = "\t"
	}

	res := make([]string, len(table))
	for i, fields := range table {
		res[i] = strings.Join(fields, delimiter)
	}
	return res
}
