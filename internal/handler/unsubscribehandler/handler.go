package unsubscribehandler

import (
	"fmt"

	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
)

type Handler struct {
	id           int
	subscription []model.BangumiBase
}

func NewHandler(id int) (*Handler, error) {
	subscription, err := path.ReadSubscriptionConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to init handler: %w", err)
	}

	return &Handler{id, subscription}, nil
}

func (h *Handler) Run() {
	switch h.id {
	case -1:
		UnsubscribeAll(h.subscription)
	default:
		UnsubscribeByID(h.subscription, h.id)
	}
	log.Info("unsubscribe 任务结束")
}

func UnsubscribeAll(subscriptions []model.BangumiBase) {
	var failures []model.BangumiBase
	for _, sub := range subscriptions {
		if err := path.DeleteJSONConfigFile(sub.ConfigFileName()); err != nil {
			log.Warnf("取消订阅 %s 错误: %s", sub.Name, err)
			failures = append(failures, sub)
			continue
		}
		log.Debugf("已取消订阅: %s", sub.Name)
	}

	if err := path.DeleteJSONConfigFile(path.SubscriptionConfigFile); err != nil {
		log.Errorf("删除订阅列表失败: %s", err)
	}

	if len(failures) != 0 {
		log.Warnf("%d 番剧订阅未成功移除, 请重试或使用 --id 指定删除:", len(failures))
		for _, b := range failures {
			log.Debug(b.Name)
		}
		return
	}

	log.Successf("成功取消 %d 部番剧的订阅", len(subscriptions))
}

func UnsubscribeByID(subscriptions []model.BangumiBase, id int) {
	updated := filterOutByID(subscriptions, id)
	if len(updated) == len(subscriptions) {
		log.Errorf("未找到 ID 为 %d 的订阅", id)
		return
	}

	if err := path.SaveJSONConfigFile(path.SubscriptionConfigFile, updated); err != nil {
		log.Errorf("保存订阅配置失败: %s", err)
		return
	}

	fn := fmt.Sprintf("%d.json", id)
	if err := path.DeleteJSONConfigFile(fn); err != nil {
		log.Errorf("删除配置文件 %s 失败: %s", fn, err)
		return
	}

	log.Successf("成功取消订阅: %d", id)
}

func filterOutByID(subscriptions []model.BangumiBase, id int) []model.BangumiBase {
	targetID := fmt.Sprintf("%d", id)

	var result []model.BangumiBase
	for _, s := range subscriptions {
		if s.ID != targetID {
			result = append(result, s)
		}
	}
	return result
}
