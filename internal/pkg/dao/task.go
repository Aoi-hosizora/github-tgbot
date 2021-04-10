package dao

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-db/xredis"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
	"strings"
	"time"
)

const MagicToken = "$$"

func getActivityPattern(chatId, id, t, repo string) string {
	repo = strings.ReplaceAll(repo, "-", MagicToken)
	return fmt.Sprintf("gh-activity-ev-%s-%s-%s-%s", chatId, id, t, repo)
}

func parseActivityPattern(key string) (chatId int64, id, t, repo string) {
	sp := strings.Split(key, "-")
	chatId, _ = xnumber.ParseInt64(sp[3], 10)
	id = sp[4]
	t = sp[5]
	repo = strings.ReplaceAll(sp[6], MagicToken, "-")
	return
}

func getIssuePattern(chatId, id, event, repo, ct, num string) string {
	event = strings.ReplaceAll(event, "-", MagicToken)
	repo = strings.ReplaceAll(repo, "-", MagicToken)
	return fmt.Sprintf("gh-issue-ev-%s-%s-%s-%s-%s-%s", chatId, id, event, repo, ct, num)
}

func parseIssuePattern(key string) (chatId, id int64, event, repo string, ct time.Time, num int32) {
	sp := strings.Split(key, "-")
	chatId, _ = xnumber.ParseInt64(sp[3], 10)
	id, _ = xnumber.ParseInt64(sp[4], 10)
	event = strings.ReplaceAll(sp[5], MagicToken, "-")
	repo = strings.ReplaceAll(sp[6], MagicToken, "-")
	ctn, _ := xnumber.ParseInt64(sp[7], 10)
	ct = time.Unix(ctn, 0)
	num, _ = xnumber.ParseInt32(sp[8], 10)
	return
}

func GetOldActivities(chatId int64) ([]*model.ActivityEvent, bool) {
	pattern := getActivityPattern(xnumber.I64toa(chatId), "*", "*", "*")
	keys, err := database.Redis().Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, false
	}

	evs := make([]*model.ActivityEvent, len(keys))
	for idx := range evs {
		_, id, t, repo := parseActivityPattern(keys[idx])
		evs[idx] = &model.ActivityEvent{
			Id: id, Type: t, Repo: &struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			}{Name: repo},
		}
	}
	return evs, true
}

func SetOldActivities(chatId int64, evs []*model.ActivityEvent) bool {
	chatIdStr := xnumber.I64toa(chatId)
	pattern := getActivityPattern(chatIdStr, "*", "*", "*")
	_, err := xredis.DelAll(database.Redis(), context.Background(), pattern)
	if err != nil {
		return false
	}

	keys := make([]string, 0)
	values := make([]string, 0)
	for _, ev := range evs {
		pattern := getActivityPattern(chatIdStr, ev.Id, ev.Type, ev.Repo.Name)
		keys = append(keys, pattern)
		values = append(values, chatIdStr)
	}

	_, err = xredis.SetAll(database.Redis(), context.Background(), keys, values)
	return err == nil
}

func GetOldIssues(chatId int64) ([]*model.IssueEvent, bool) {
	pattern := getIssuePattern(xnumber.I64toa(chatId), "*", "*", "*", "*", "*")
	keys, err := database.Redis().Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, false
	}

	evs := make([]*model.IssueEvent, 0)
	for _, key := range keys {
		_, id, event, repo, ct, num := parseIssuePattern(key)
		evs = append(evs, &model.IssueEvent{Id: id, Event: event, Repo: repo, CreatedAt: ct, Number: num})
	}
	return evs, true
}

func SetOldIssues(chatId int64, evs []*model.IssueEvent) bool {
	chatIdStr := xnumber.I64toa(chatId)
	pattern := getIssuePattern(chatIdStr, "*", "*", "*", "*", "*")
	_, err := xredis.DelAll(database.Redis(), context.Background(), pattern)
	if err != nil {
		return false
	}

	// set to redis, and check if duplicate in last history
	keys := make([]string, 0)
	values := make([]string, 0)
	for _, ev := range evs {
		pattern := getIssuePattern(chatIdStr, xnumber.I64toa(ev.Id), ev.Event, ev.Repo, xnumber.I64toa(ev.CreatedAt.Unix()), xnumber.I32toa(ev.Number))
		keys = append(keys, pattern)
		values = append(values, chatIdStr)
	}

	_, err = xredis.SetAll(database.Redis(), context.Background(), keys, values)
	return err == nil
}