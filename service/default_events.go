package service

var DefaultEventNames = []string{
	"onTestStart",
	"onTestEnd",
	"onCmdStart",
	"onCmdEnd",
	"onCmdExitWithErr",
}

func init() {
	for _, name := range DefaultEventNames {
		if err := CreateEvent(name); err != nil {
			panic(err)
		}
	}
}
