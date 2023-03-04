package cfmt

var defaultPrinter = &colorPrinter{}

func Black() stylerBackgroundSetter   { return defaultPrinter.black() }
func Red() stylerBackgroundSetter     { return defaultPrinter.red() }
func Green() stylerBackgroundSetter   { return defaultPrinter.green() }
func Yellow() stylerBackgroundSetter  { return defaultPrinter.yellow() }
func Blue() stylerBackgroundSetter    { return defaultPrinter.blue() }
func Magenta() stylerBackgroundSetter { return defaultPrinter.magenta() }
func Cyan() stylerBackgroundSetter    { return defaultPrinter.cyan() }
func White() stylerBackgroundSetter   { return defaultPrinter.white() }

func Bold() stylerBackgroundSetter      { return defaultPrinter.Bold() }
func Underline() stylerBackgroundSetter { return defaultPrinter.Underline() }
func Inverse() stylerBackgroundSetter   { return defaultPrinter.Inverse() }
func Reset() stylerBackgroundSetter     { return defaultPrinter.Reset() }

func BlackBG() textStyler   { return defaultPrinter.BlackBG() }
func RedBG() textStyler     { return defaultPrinter.RedBG() }
func GreenBG() textStyler   { return defaultPrinter.GreenBG() }
func YellowBG() textStyler  { return defaultPrinter.YellowBG() }
func BlueBG() textStyler    { return defaultPrinter.BlueBG() }
func MagentaBG() textStyler { return defaultPrinter.MagentaBG() }
func CyanBG() textStyler    { return defaultPrinter.CyanBG() }
func WhiteBG() textStyler   { return defaultPrinter.WhiteBG() }
