package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"


	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

var RoleMessage discord.Message
var Roles = map[string]string{
    "🧈": "butter",
    "⚙️": "cog",
}

func main() {
	client, err := disgo.New(os.Getenv("DISCORD_TOKEN"),
		// set gateway options
		bot.WithGatewayConfigOpts(
			// set enabled intents
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentDirectMessages,
				gateway.IntentGuildVoiceStates,
				gateway.IntentGuildMessageReactions,
			),
		),
		// add event listeners
		bot.WithEventListenerFunc(func(e *events.GuildVoiceJoin) {
			err := VoiceRoleAdjustment(true, e.GenericGuildVoiceState)
			if err != nil {
				log.Printf("Failed to add member role: %s", err)
			}
		}),
		bot.WithEventListenerFunc(func(e *events.GuildVoiceLeave) {
			err := VoiceRoleAdjustment(false, e.GenericGuildVoiceState)
			if err != nil {
				log.Printf("Failed to remove member role: %s", err)
			}
		}),
	)

	if err != nil {
		panic(err)
	}
	// connect to the gateway
	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func VoiceRoleAdjustment(joined bool, e *events.GenericGuildVoiceState) error {
	userid := e.Member.User.ID
	gid := e.VoiceState.GuildID

	roles, err := e.Client().Rest().GetRoles(gid)

	if err != nil {
		return fmt.Errorf("failed to get the role: %w", err)
	}

	var voicerole discord.Role
	for _, role := range roles {
		if role.Name == "Voice" {
			voicerole = role
		}
	}

	if voicerole.ID == 0 {
		return errors.New("failed to find the role \"Voice\"") 
	}

	if joined {
		err = e.Client().Rest().AddMemberRole(gid, userid, voicerole.ID)
	} else {
		err = e.Client().Rest().RemoveMemberRole(gid, userid, voicerole.ID)
	}
	
	if err != nil {
		return fmt.Errorf("failed to adjust the role: %w", err)
		
	}

	return nil
}

func aquireRoleMessage(client bot.Client) error {
	rolechannelname := "welcome"
	gid, err := snowflake.Parse("635550573470416896")	

	if err != nil {
		return fmt.Errorf("failed to parse snowflake: %w", err)
	}
	
	channels, err := client.Rest().GetGuildChannels(gid)

	if err != nil {
		return fmt.Errorf("failed to get guild channels: %w", err)
	}

	var rolechannel discord.Channel
	for _, channel := range channels {
		if channel.Name() == rolechannelname {
			rolechannel = channel
		}
	}

	var nullflake snowflake.ID

	message, err := client.Rest().GetMessages(rolechannel.ID(), nullflake, nullflake, nullflake, 1)

	if err != nil {
		return fmt.Errorf("failed to get guild channels: %w", err)
	}

	if len(message) == 0 {
		// createdmessage = discord.MessageCreateBuilder
		
		// client.Rest().CreateMessage(rolechannel.ID(), )
	}

}
