package main

import (
	"os"

	"github.com/bluele/slack"
)

type slackNotifier struct {
	slackInfo *slack.WebHook
	channel   string
	iconEmoji string
	username  string
}

func newSlackNotifier() *slackNotifier {

	webhook, ok := os.LookupEnv("SLACK_WEBHOOK")
	if !ok {
		panic("NO SLACK WEBHOOK")
	}

	channel, ok := os.LookupEnv("SLACK_CHANNEL")
	if !ok {
		panic("NO SLACK CHANNEL")
	}

	iconEmoji, ok := os.LookupEnv("SLACK_ICON_EMOJI")
	if !ok {
		iconEmoji = ":cookie:"
	}

	username, ok := os.LookupEnv("SLACK_USERNAME")
	if !ok {
		username = "Webhook"
	}

	return &slackNotifier{
		slackInfo: slack.NewWebHook(webhook),
		channel:   channel,
		iconEmoji: iconEmoji,
		username:  username,
	}

}

func (slacker *slackNotifier) notify(text string) {
	slacker.slackInfo.PostMessage(&slack.WebHookPostPayload{
		Text:      text,
		Channel:   slacker.channel,
		IconEmoji: slacker.iconEmoji,
		Username:  slacker.username,
	})
}

func (slacker *slackNotifier) notifyWithAttachment(text string, attachment *slack.Attachment) {
	slacker.slackInfo.PostMessage(&slack.WebHookPostPayload{
		Text:      text,
		Channel:   slacker.channel,
		IconEmoji: slacker.iconEmoji,
		Username:  slacker.username,
		Attachments: []*slack.Attachment{
			attachment,
		},
	})
}
