package init

import (
	"github.com/Mr-LvGJ/jota/log"
)

func InitLog(fileName string) {
	c := log.DefaultConfig()
	c.Filename = fileName
	if err := log.NewGlobal(c); err != nil {
		panic(err)
	}
}
