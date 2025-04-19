package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/console"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

func filter(rss mikan.RSS) ([]string, *model.Filters) {
	r := rss

	console.Info("当前订阅的RSS包含以下结果:")
	filteredRSS(r)

	filters := promptFilters()

	console.Infof("根据关键词 %v 筛选结果...", filters)
	r = applyFilters(r, filters)

	console.Info("筛选后的结果如下:")
	filteredRSS(r)

	torrents := r.TorrentURLs()

	return torrents, &filters
}

func applyFilters(rss mikan.RSS, filters model.Filters) mikan.RSS {
	return rss.Filter(filters)
}

func promptFilters() model.Filters {
	var input string
	input = prompt.ReadUserInput("请输入想要包含的关键词，多个请以英文逗号(,)隔开, 完成后请按回车")
	include := libs.SplitToSlice(input, ",")

	return model.Filters{
		Include: include,
	}
}

func filteredRSS(rss mikan.RSS) {
	res := rss.String()
	console.Plain(res)
}
