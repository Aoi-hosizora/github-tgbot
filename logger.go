package main

import (
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func botLog(fromMsg *telebot.Message, sendMsg *telebot.Message, err error) {
	if fromMsg != nil { // bot
		if sendMsg == nil { // receive
			log.Printf("[bot] receive \t %d \t \"%s\"\n", fromMsg.ID, fromMsg.Text)
		} else { // send
			if err == nil {
				timeSpan := float64(sendMsg.Time().Sub(fromMsg.Time()).Nanoseconds()) / 1e6
				log.Printf("[bot] reply \t %d \t %.0fms (from %d \"%s\")\n", sendMsg.ID, timeSpan, fromMsg.ID, fromMsg.Text)
			} else {
				log.Printf("[bot] failed to reply bot of %d\n", fromMsg.ID)
			}
		}
	} else { // channel
		if err == nil {
			log.Printf("[channel] send \t %d \t \"%s\"", sendMsg.ID, sendMsg.Chat.Title)
		} else {
			log.Printf("[channel] failed to send message to channel \"%s\"\n", sendMsg.Chat.Title)
		}
	}
}