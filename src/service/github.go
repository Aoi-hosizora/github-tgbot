package service

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"net/http"
)

func CheckUser(username string, private bool, token string) (bool, error) {
	url := fmt.Sprintf(UserApi, username)
	header := &http.Header{}
	if private {
		header.Add("Authorization", fmt.Sprintf("Token %s", token))
	}

	code, _, err := HttpGet(url, header)
	return code == 200, err
}

func GetActivityEvents(username string, private bool, token string, page int) (response string, err error) {
	url := fmt.Sprintf(ActivityEventApi, username)
	url = fmt.Sprintf("%s?page=%d", url, page)
	header := &http.Header{}
	if private {
		header.Add("Authorization", fmt.Sprintf("Token %s", token))
	}

	_, bs, err := HttpGet(url, header)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func GetIssueEvents(username string, private bool, token string, page int) (response string, err error) {
	url := fmt.Sprintf(IssueEventApi, username)
	url = fmt.Sprintf("%s?page=%d", url, page)
	header := &http.Header{}
	if private {
		header.Add("Authorization", fmt.Sprintf("Token %s", token))
	}

	_, bs, err := HttpGet(url, header)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func RenderActivities(objs []*model.ActivityEvent) string {
	if len(objs) == 1 {
		return RenderActivity(objs[0])
	}

	result := ""
	for idx, obj := range objs {
		result += fmt.Sprintf("%d. %s\n", idx+1, RenderActivity(obj))
	}
	return result[:len(result)-1]
}

func RenderIssues(objs []*model.IssueEvent) string {
	if len(objs) == 1 {
		return RenderIssue(objs[0])
	}

	result := ""
	for idx, obj := range objs {
		result += fmt.Sprintf("%d. %s\n", idx+1, RenderIssue(obj))
	}
	return result[:len(result)-1]
}
