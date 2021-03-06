package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/button"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/dao"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	ISSUE_ONLY_FOR_TOKEN           = "Send issue can only be allowed for users that bind with token."
	ALLOW_ISSUE_Q                  = "Would you need to filter the message generated by yourself?"
	ALLOW_ISSUE_FAILED             = "Failed to allow bot to send issue events periodically."
	ALLOW_ISSUE_FILTER_SUCCESS     = "Success to allow bot to send issue events periodically (filter message generated by myself)."
	ALLOW_ISSUE_NOT_FILTER_SUCCESS = "Success to allow bot to send issue events periodically (do not filter message generated by myself)."

	DISALLOW_ISSUE_FAILED  = "Failed to disallow bot to send issue events periodically."
	DISALLOW_ISSUE_SUCCESS = "Success to disallow bot to send issue events periodically."

	GITHUB_PAGE_Q     = "Please send the page number you want to get. Send /cancel to cancel."
	UNEXPECTED_NUMBER = "Unexpected page number. Please send an integer value. Send /cancel to cancel."
	EMPTY_EVENT       = "You have empty event."
)

// /allowissue
func AllowIssueCtrl(m *telebot.Message) {
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}
	if user.Token == "" {
		_ = server.Bot().Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	_ = server.Bot().Reply(m, ALLOW_ISSUE_Q, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*button.InlineBtnFilter, *button.InlineBtnNotFilter}, {*button.InlineBtnCancel},
		},
	})
}

// button.InlineBtnFilter
func InlineBtnFilterCtrl(c *telebot.Callback) {
	m := c.Message
	_, _ = server.Bot().Edit(m, fmt.Sprintf("%s (filter)", m.Text))

	flag := ""
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		flag = BIND_NOT_YET
	} else if user.Token == "" {
		flag = ISSUE_ONLY_FOR_TOKEN
	} else {
		status := dao.UpdateUserAllowIssue(user.ChatID, true, true)
		if status == xstatus.DbNotFound {
			flag = BIND_NOT_YET
		} else if status == xstatus.DbFailed {
			flag = ALLOW_ISSUE_FAILED
		} else {
			flag = ALLOW_ISSUE_FILTER_SUCCESS
		}
	}
	_ = server.Bot().Reply(m, flag)
}

// button.InlineBtnNotFilter
func InlineBtnNotFilterCtrl(c *telebot.Callback) {
	m := c.Message
	_, _ = server.Bot().Edit(m, fmt.Sprintf("%s (not filter)", m.Text))

	flag := ""
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		flag = BIND_NOT_YET
	} else if user.Token == "" {
		flag = ISSUE_ONLY_FOR_TOKEN
	} else {
		status := dao.UpdateUserAllowIssue(user.ChatID, true, false)
		if status == xstatus.DbNotFound {
			flag = BIND_NOT_YET
		} else if status == xstatus.DbFailed {
			flag = ALLOW_ISSUE_FAILED
		} else {
			flag = ALLOW_ISSUE_NOT_FILTER_SUCCESS
		}
	}
	_ = server.Bot().Reply(m, flag)
}

// /disallowissue
func DisallowIssueCtrl(m *telebot.Message) {
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	} else if user.Token == "" {
		_ = server.Bot().Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	flag := ""
	status := dao.UpdateUserAllowIssue(user.ChatID, false, false)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = DISALLOW_ISSUE_FAILED
	} else {
		flag = DISALLOW_ISSUE_SUCCESS
	}
	_ = server.Bot().Reply(m, flag)
}

// /activity
func ActivityCtrl(m *telebot.Message) {
	m.Text = "1"
	FromActivityPageCtrl(m)
}

// /activitypage
func ActivityPageCtrl(m *telebot.Message) {
	server.Bot().SetStatus(m.Chat.ID, fsm.ActivityPage)
	_ = server.Bot().Reply(m, GITHUB_PAGE_Q)
}

// fsm.ActivityPage
func FromActivityPageCtrl(m *telebot.Message) {
	page, err := xnumber.Atoi(m.Text)
	if err != nil {
		_ = server.Bot().Reply(m, UNEXPECTED_NUMBER)
		return
	}
	if page <= 0 {
		page = 1
	}
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		server.Bot().SetStatus(m.Chat.ID, fsm.None)
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}

	flag := ""
	v2md := false
	if resp, err := service.GetActivityEvents(user.Username, user.Token, page); err != nil {
		flag = GITHUB_FAILED
	} else if events, err := model.UnmarshalActivityEvents(resp); err != nil {
		flag = GITHUB_FAILED
	} else if render := service.RenderActivityEvents(events); render == "" {
		flag = EMPTY_EVENT
	} else {
		flag = service.ConcatListAndUsername(render, user.Username) + fmt.Sprintf(" \\(page %d\\)", page) // <<<
		v2md = true
	}

	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	if !v2md {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdown)
	} else {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdownV2)
	}
}

// /issue
func IssueCtrl(m *telebot.Message) {
	m.Text = "1"
	FromIssuePageCtrl(m)
}

// /issuepage
func IssuePageCtrl(m *telebot.Message) {
	server.Bot().SetStatus(m.Chat.ID, fsm.IssuePage)
	_ = server.Bot().Reply(m, GITHUB_PAGE_Q)
}

// fsm.IssuePage
func FromIssuePageCtrl(m *telebot.Message) {
	page, err := xnumber.Atoi(m.Text)
	if err != nil {
		_ = server.Bot().Reply(m, UNEXPECTED_NUMBER)
		return
	}
	if page <= 0 {
		page = 1
	}
	user := dao.QueryUser(m.Chat.ID)
	if user == nil {
		server.Bot().SetStatus(m.Chat.ID, fsm.None)
		_ = server.Bot().Reply(m, BIND_NOT_YET)
		return
	}
	if user.Token == "" {
		_ = server.Bot().Reply(m, ISSUE_ONLY_FOR_TOKEN)
		return
	}

	flag := ""
	v2md := false
	if resp, err := service.GetIssueEvents(user.Username, user.Token, page); err != nil {
		flag = GITHUB_FAILED
	} else if events, err := model.UnmarshalIssueEvents(resp); err != nil {
		flag = GITHUB_FAILED
	} else if render := service.RenderIssueEvents(events); render == "" {
		flag = EMPTY_EVENT
	} else {
		flag = service.ConcatListAndUsername(render, user.Username) + fmt.Sprintf(" \\(page %d\\)", page) // <<<
		v2md = true
	}

	server.Bot().SetStatus(m.Chat.ID, fsm.None)
	if !v2md {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdown)
	} else {
		_ = server.Bot().Reply(m, flag, telebot.ModeMarkdownV2)
	}
}
