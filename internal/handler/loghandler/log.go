package loghandler

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/path"
	"time"
)

func Run() {
	date := time.Now().Format("2006-01-02")

	fn := fmt.Sprintf("%s.log", date)

	data, err := path.ReadLogFile(fn)
	if err != nil {
		return
	}

	fmt.Println(data)
}
