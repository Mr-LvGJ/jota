package init

import (
	"github.com/Mr-LvGJ/log"
)

func InitLog(fileName string) {
	c := log.DefaultConfig()
	c.Filename = fileName
	if err := log.NewGlobal(c); err != nil {
		panic(err)
	}
}
