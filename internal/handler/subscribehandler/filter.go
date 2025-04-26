package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

func filterRSS(rss mikan.RSS) ([]string, *model.Filters) {
	log.Infof("当前订阅的RSS包含以下结果: \n%s", rss.String())

	filters := promptFilters()
	log.Infof("根据关键词 %v 筛选结果...", filters)

	rss = rss.Filter(filters)
	log.Infof("筛选后的结果如下:\n%s", rss.String())

	return rss.TorrentURLs(), &filters
}

func promptFilters() model.Filters {
	var input string
	input = prompt.ReadUserInput("请输入想要包含的关键词，多个请以英文逗号(,)隔开, 完成后请按回车")
	include := libs.SplitToSlice(input, ",")

	return model.Filters{
		Include: include,
	}
}
