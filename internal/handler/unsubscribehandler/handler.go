package unsubscribehandler

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
)

func Run(id int) {
	if id == -1 {
		removeAll()
		return
	}

	remove(id)
}

func removeAll() {
	var subscription []model.BangumiBase
	if err := path.ReadJSONConfigFile(path.SubscriptionConfigFile, &subscription); err != nil {
		log.Errorf("读取本地订阅配置文件错误:%s", err)
		return
	}

	var failures []model.BangumiBase
	for _, s := range subscription {
		fn := s.ConfigFileName()
		if err := path.DeleteJSONConfigFile(fn); err != nil {
			log.Warnf("取消订阅%s错误:%s", s.Name, err)
			failures = append(failures, s)
		}
		log.Debugf("已取消订阅:%s", s.Name)
	}

	if err := path.DeleteJSONConfigFile(path.SubscriptionConfigFile); err != nil {
		log.Errorf("删除本地订阅配置文件错误:%s", err)
	}

	if len(failures) == 0 {
		log.Successf("已成功取消订阅 %d 部番剧", len(subscription))
		return
	}

	log.Warnf("以下%d 部番剧取消订阅失败, 请重试 bangumi unsubscribe 或指定 --id", len(failures))
	for _, b := range failures {
		log.Debug(b.Name)
	}
}

func remove(id int) {
	var subscription []model.BangumiBase
	if err := path.ReadJSONConfigFile(path.SubscriptionConfigFile, &subscription); err != nil {
		log.Errorf("读取本地订阅配置文件错误:%s", err)
		return
	}

	subscription = removeByID(subscription, id)
	if err := path.SaveJSONConfigFile(path.SubscriptionConfigFile, subscription); err != nil {
		log.Errorf("保存本地订阅配置文件错误:%s", err)
		return
	}

	fn := fmt.Sprintf("%d.json", id)
	if err := path.DeleteJSONConfigFile(fn); err != nil {
		log.Errorf("取消订阅失败:%s", err)
		return
	}

	log.Successf("取消订阅 %d 成功", id)
}

func removeByID(subscription []model.BangumiBase, id int) []model.BangumiBase {
	ids := fmt.Sprintf("%d", id)

	var res []model.BangumiBase
	for _, s := range subscription {
		if s.ID != ids {
			res = append(res, s)
		}
	}
	return res
}
