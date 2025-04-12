package loghandler

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/path"
)

func Run() {
	data, err := path.ReadLogFile()
	if err != nil {
		return
	}

	fmt.Println(data)
}
