package main

import (
	"flag"
	"github.com/Aoi-hosizora/github-telebot/internal/bot"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/database"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/Aoi-hosizora/github-telebot/internal/task"
	"log"
)

var (
	fHelp   = flag.Bool("h", false, "show help")
	fConfig = flag.String("config", "./config.yaml", "change the config path")
)

func main() {
	flag.Parse()
	if *fHelp {
		flag.Usage()
	} else {
		run()
	}
}

func run() {
	err := config.Load(*fConfig)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	err = logger.Setup()
	if err != nil {
		log.Fatalln("Failed to setup logger:", err)
	}
	err = database.SetupGorm()
	if err != nil {
		log.Fatalln("Failed to connect mysql:", err)
	}
	err = database.SetupRedis()
	err = database.SetupRedis()
	if err != nil {
		log.Fatalln("Failed to connect redis:", err)
	}

	err = bot.Setup()
	if err != nil {
		log.Fatalln("Failed to load telebot:", err)
	}
	err = task.Setup()
	if err != nil {
		log.Fatalln("Failed to setup cron:", err)
	}

	defer task.Cron().Stop()
	task.Cron().Start()

	defer server.Bot().Stop()
	server.Bot().Start()
}