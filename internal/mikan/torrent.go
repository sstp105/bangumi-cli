package mikan

type TorrentsFilter struct {
	RSS     RSS
	Filters Filters
}

func (tf *TorrentsFilter) Apply() RSS {
	return Filter(tf.RSS, tf.Filters)
}

func (tf *TorrentsFilter) URLs() []string {
	filtered := tf.Apply()
	var urls []string
	for _, item := range filtered.Channel.Items {
		urls = append(urls, item.Enclosure.URL)
	}
	return urls
}
