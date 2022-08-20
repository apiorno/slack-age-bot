package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("CommandEvent")
		fmt.Println(event.Timestamp.String())
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
	}
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	go printCommandEvents(bot.CommandEvents())

	bot.Command("My yob is <year>", &slacker.CommandDefinition{
		Description: "yob calculator",
		Examples:    []string{"my yob is 2020"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println(err)
			}
			actualYear, _ := strconv.Atoi(os.Getenv("ACTUAL_YEAR"))
			age := actualYear - yob
			r := fmt.Sprintf("age is %d", age)
			response.Reply(r)
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
