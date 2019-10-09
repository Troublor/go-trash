package service

var DefaultEventNames = []string{
	"onTestStart",
	"onTestEnd",
}

func init() {
	for _, name := range DefaultEventNames {
		if err := CreateEvent(name); err != nil {
			panic(err)
		}
	}
}
