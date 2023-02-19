package main

import (
	"flag"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dsychin/habitica-todoist-task-redeemer/config"
	"github.com/dsychin/habitica-todoist-task-redeemer/controller"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	flag.String("habitica_key", "", "Your Habitica API Token")
	flag.String("habitica_user_id", "", "Your Habitica User ID")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("redeem")
	viper.BindEnv("habitica_key")
	viper.BindEnv("habitica_user_id")

	config.APIToken = viper.GetString("habitica_key")
	config.UserID = viper.GetString("habitica_user_id")

	if config.APIToken == "" {
		log.Fatal("API Token is empty. Please set the REDEEM_HABITICA_KEY environment variable.")
	}
	if config.UserID == "" {
		log.Fatal("API Token is empty. Please set the REDEEM_HABITICA_USER_ID environment variable.")
	}

	lambda.Start(controller.HandleRequest)
}
