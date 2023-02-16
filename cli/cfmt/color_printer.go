package cfmt

import (
	"fmt"
	"io"
	"os"
)

type printer interface {
	Printf(format string, args ...any) (n int, err error)
	Print(args ...any) (n int, err error)
	Println(args ...any) (n int, err error)
	Fprintf(w io.Writer, format string, args ...any) (n int, err error)
	Fprint(w io.Writer, args ...any) (n int, err error)
	Fprintln(w io.Writer, args ...any) (n int, err error)
}

type stringPrinter interface {
	Sprintf(format string, args ...any) string
	Sprint(args ...any) string
	Sprintln(args ...any) string
}

type colorChooser interface {
	printer
	stringPrinter
	black() stylerBackgroundSetter
	red() stylerBackgroundSetter
	green() stylerBackgroundSetter
	yellow() stylerBackgroundSetter
	blue() stylerBackgroundSetter
	magenta() stylerBackgroundSetter
	cyan() stylerBackgroundSetter
	white() stylerBackgroundSetter
}

type textBackground interface {
	printer
	stringPrinter
	BlackBG() textStyler
	RedBG() textStyler
	GreenBG() textStyler
	YellowBG() textStyler
	BlueBG() textStyler
	MagentaBG() textStyler
	CyanBG() textStyler
	WhiteBG() textStyler
}

type textStyler interface {
	printer
	stringPrinter
	Underline() stylerBackgroundSetter
	Inverse() stylerBackgroundSetter
	Reset() stylerBackgroundSetter
	Bold() stylerBackgroundSetter
}

type stylerBackgroundSetter interface {
	textStyler
	textBackground
}

type colorPrinter struct {
	option     string
	color      string
	background string
}

func (c *colorPrinter) black() stylerBackgroundSetter   { c.color = "30;"; return c }
func (c *colorPrinter) red() stylerBackgroundSetter     { c.color = "31;"; return c }
func (c *colorPrinter) green() stylerBackgroundSetter   { c.color = "32;"; return c }
func (c *colorPrinter) yellow() stylerBackgroundSetter  { c.color = "33;"; return c }
func (c *colorPrinter) blue() stylerBackgroundSetter    { c.color = "34;"; return c }
func (c *colorPrinter) magenta() stylerBackgroundSetter { c.color = "35;"; return c }
func (c *colorPrinter) cyan() stylerBackgroundSetter    { c.color = "36;"; return c }
func (c *colorPrinter) white() stylerBackgroundSetter   { c.color = "97;"; return c }

func (c *colorPrinter) Bold() stylerBackgroundSetter      { c.option += "1;"; return c }
func (c *colorPrinter) Underline() stylerBackgroundSetter { c.option += "4;"; return c }
func (c *colorPrinter) Inverse() stylerBackgroundSetter   { c.option += "7;"; return c }
func (c *colorPrinter) Reset() stylerBackgroundSetter {
	c.option = "0"
	c.color = ""
	c.background = ""
	return c
}

func (c *colorPrinter) BlackBG() textStyler   { c.background = "40"; return c }
func (c *colorPrinter) RedBG() textStyler     { c.background = "41"; return c }
func (c *colorPrinter) GreenBG() textStyler   { c.background = "42"; return c }
func (c *colorPrinter) YellowBG() textStyler  { c.background = "43"; return c }
func (c *colorPrinter) BlueBG() textStyler    { c.background = "44"; return c }
func (c *colorPrinter) MagentaBG() textStyler { c.background = "45"; return c }
func (c *colorPrinter) CyanBG() textStyler    { c.background = "46"; return c }
func (c *colorPrinter) WhiteBG() textStyler   { c.background = "47"; return c }

func (c *colorPrinter) Print(args ...any) (n int, err error) {
	return c.Fprint(os.Stdin, args...)
}

func (c *colorPrinter) Println(args ...any) (n int, err error) {
	return c.Fprintln(os.Stdin, args...)
}

func (c *colorPrinter) Printf(format string, args ...any) (n int, err error) {
	return c.Fprintf(os.Stdin, format, args...)
}

func (c *colorPrinter) Fprintf(w io.Writer, format string, args ...any) (n int, err error) {
	return fmt.Fprintf(w, c.prepareSettings(format), args...)
}

func (c *colorPrinter) Fprint(w io.Writer, args ...any) (n int, err error) {
	return fmt.Fprint(w, c.prepareSettings(fmt.Sprint(args...)))
}

func (c *colorPrinter) Fprintln(w io.Writer, args ...any) (n int, err error) {
	return fmt.Fprintln(w, c.prepareSettings(fmt.Sprint(args...)))
}

func (c *colorPrinter) Sprintf(format string, args ...any) string {
	return fmt.Sprintf(c.prepareSettings(format), args...)
}

func (c *colorPrinter) Sprint(args ...any) string {
	return fmt.Sprint(c.prepareSettings(fmt.Sprint(args...)))
}

func (c *colorPrinter) Sprintln(args ...any) string {
	return fmt.Sprintln(c.prepareSettings(fmt.Sprint(args...)))
}

func (c *colorPrinter) prepareSettings(text string) string {
	if len(c.background) == 0 && len(c.color) > 0 {
		c.color = c.color[:len(c.color)-1]
	}
	if len(c.color) == 0 && len(c.background) == 0 && len(c.option) > 0 {
		c.option = c.option[:len(c.option)-1]
	}

	return "\033[" + c.option + c.color + c.background + "m" + text + "\033[0m"
}

func (c *colorPrinter) resetSettings() {
	c.color = ""
	c.option = ""
	c.background = ""
}
