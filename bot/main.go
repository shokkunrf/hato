package main

import (
	"hato/config"
	"hato/discord"
	"hato/mqtt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	brokerConfig, err := config.GetBrokerConfig()
	if err != nil {
		log.Fatalln(err)
	}
	publisher := mqtt.MakePublisher(brokerConfig)
	err = publisher.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		publisher.Disconnect()
	}()

	discordConfig, err := config.GetDiscordConfig()
	if err != nil {
		log.Fatalln(err)
	}
	bot, err := discord.MakeBot(*discordConfig, publisher)
	if err != nil {
		log.Fatalln(err)
	}
	err = bot.Start()
	if err != nil {
		log.Fatalln(err)
	}
	defer bot.Stop()

	log.Println("Start")

	// 終了を待機
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	select {
	case <-signalChan:
		return
	}
}
