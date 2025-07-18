package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/prompt"
)

func filterRSS(rss mikan.RSS) ([]model.Torrent, *model.Filters) {
	log.Infof("The current subscribed RSS contains the following results:\n%s", rss.String())

	filters := promptFilters()
	log.Infof("Filtering results based on keywords %v...", filters)

	rss = rss.Filter(filters)
	log.Infof("The filtered results are as follows:\n%s", rss.String())

	return rss.Torrents(), &filters
}

func promptFilters() model.Filters {
	var input string
	input = prompt.ReadUserInput("Please enter keywords to include, separated by commas (,). Press Enter when done.")
	include := libs.SplitToSlice(input, ",")

	return model.Filters{
		Include: include,
	}
}
