package task

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/bot"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"github.com/Aoi-hosizora/ah-tgbot/util"
	"time"
)

var (
	oldActivities = make(map[int64][]*model.ActivityEvent, 0)
	oldIssues     = make(map[int64][]*model.IssueEvent, 0)
)

func sliceActivityDiff(s1 []*model.ActivityEvent, s2 []*model.ActivityEvent) []*model.ActivityEvent {
	result := make([]*model.ActivityEvent, 0)
	for _, item1 := range s1 {
		exist := false
		for _, item2 := range s2 {
			if model.ActivityEventEqual(item1, item2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, item1)
		}
	}
	return result
}

func sliceIssueDiff(s1 []*model.IssueEvent, s2 []*model.IssueEvent) []*model.IssueEvent {
	result := make([]*model.IssueEvent, 0)
	for _, item1 := range s1 {
		exist := false
		for _, item2 := range s2 {
			if model.IssueEventEqual(item1, item2) {
				exist = true
				break
			}
		}
		if !exist {
			result = append(result, item1)
		}
	}
	return result
}

func ActivityTask() {
	for {
		users := model.GetUsers()
		for _, user := range users {
			// get event and unmarshal
			resp, err := util.GetGithubActivityEvents(user.Username, user.Private, user.Token, 1)
			if err != nil {
				continue
			}
			events, err := model.UnmarshalActivityEvents(resp)
			if err != nil {
				continue
			}

			// check map and diff
			if _, ok := oldActivities[user.ChatID]; !ok {
				oldActivities[user.ChatID] = []*model.ActivityEvent{}
			}
			diff := sliceActivityDiff(events, oldActivities[user.ChatID])
			if len(diff) != 0 {
				// render and send
				render := util.RenderGithubActivityString(diff)
				flag := fmt.Sprintf("%s\n---\nFrom [%s](https://github.com/%s) updated.", render, user.Username, user.Username)
				bot.SendToChat(user.ChatID, flag)
			}

			// update old map
			oldActivities[user.ChatID] = events
		}

		// wait to send next time
		time.Sleep(time.Duration(config.Configs.TaskConfig.PollingActivityDuration) * time.Second)
	}
}

func IssueTask() {
	for {
		users := model.GetUsers()
		for _, user := range users {
			if !user.Private {
				continue
			}

			// get event and unmarshal
			resp, err := util.GetGithubIssueEvents(user.Username, user.Private, user.Token, 1)
			if err != nil {
				continue
			}
			events, err := model.UnmarshalIssueEvents(resp)
			if err != nil {
				continue
			}

			// check map and diff
			if _, ok := oldIssues[user.ChatID]; !ok {
				oldIssues[user.ChatID] = []*model.IssueEvent{}
			}
			diff := sliceIssueDiff(events, oldIssues[user.ChatID])
			if len(diff) != 0 {
				// render and send
				render := util.RenderGithubIssueString(diff)
				flag := fmt.Sprintf("%s\n---\nFrom [%s](https://github.com/%s) updated.", render, user.Username, user.Username)
				bot.SendToChat(user.ChatID, flag)
			}

			// update old map
			oldIssues[user.ChatID] = events
		}

		// wait to send next time
		time.Sleep(time.Duration(config.Configs.TaskConfig.PollingIssueDuration) * time.Second)
	}
}
