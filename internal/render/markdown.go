package render

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const ANSIReset = "\033[0m"

type ANSIFormat int

const (
	Regular ANSIFormat = iota
	Bold
	Dim
	Italic
	Underline
	SlowBlink
	FastBlink
	Invert
	Hide
	Strike
)

type ANSIColor int

const (
	Black ANSIColor = 30 + iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Custom
	Default
)

func createAnsiCode(format ANSIFormat, color ANSIColor, isHighIntensity, isBg bool) string {
	if isHighIntensity {
		color += 60
	}
	if isBg {
		color += 10
	}
	ansiCode := fmt.Sprintf("\033[%d;%dm", format, color)
	return ansiCode
}

func formatStrAnsi(str string, format ANSIFormat, color ANSIColor, isHighIntensity, isBg bool) string {
	ansiCode := createAnsiCode(format, color, isHighIntensity, isBg)
	return fmt.Sprintf("%s%s%s", ansiCode, str, ANSIReset)
}

func formatMdToAnsi(md string) (string, error) {
	titleRe := regexp.MustCompile(`(?m)^(#+)\s+(.*)$`)
	md = titleRe.ReplaceAllString(md, formatStrAnsi("$2", Bold, Blue, false, false))

	boldRe := regexp.MustCompile(`\*\*(.*?)\*\*`)
	md = boldRe.ReplaceAllString(md, formatStrAnsi("$1", Bold, Default, false, false))

	ulRe := regexp.MustCompile(`(?m)^(\s*)([\*\-])(\s+)`)
	md = ulRe.ReplaceAllString(md, "$1"+formatStrAnsi("$2", Bold, Blue, false, false)+"$3")
	return md, nil
}

func checkExecExist(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func Render(content string, pager string) error {
	pagerArgs := strings.Fields(pager)
	if len(pagerArgs) < 1 {
		return fmt.Errorf("render failed: pager not provided")
	}

	if !checkExecExist(pagerArgs[0]) {
		return fmt.Errorf("cannot find pager '%s' in PATH", pagerArgs[0])
	}
	cmd := exec.Command(pagerArgs[0], pagerArgs[1:]...)
	fmt.Fprintf(os.Stderr, "Launching '%s' to view content...\n", cmd.String())
	formattedContent, err := formatMdToAnsi(content)
	if err != nil {
		return err
	}
	cmd.Stdin = strings.NewReader(formattedContent)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
