package discord

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"pkd-bot/calc"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var BotCommandsChannelID = ""

var ct2blrk = map[string]string{
	"Early 3-1":   "Early 3+1",
	"Glass Neo":   "Rng Skip",
	"Overhead 4B": "Overhead 4b",
}

type BoostRoomsResponse struct {
	Name     string  `json:"name"`
	Pacelock float64 `json:"pacelock"`
	Index    int     `json:"index"`
}

var seedCache = NewSeedCache(1 * time.Hour)

// ------------------------
// /allsplits handler
// ------------------------
func AllSplitsHandle(i *discordgo.InteractionCreate) (calc.CalcSeedResult, []BoostRoomsResponse, error) {
	if s == nil {
		return calc.CalcSeedResult{}, nil, fmt.Errorf("discord session is not initialized")
	}

	// Get filter option
	filter := "all"
	if len(i.ApplicationCommandData().Options) > 0 {
		opt := i.ApplicationCommandData().Options[0]
		filter = strings.ToLower(opt.StringValue())
	}

	// Get rooms from options
	rooms := []string{}
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "rooms" && len(opt.Options) > 0 {
			for _, r := range opt.Options {
				rooms = append(rooms, strings.ToLower(r.StringValue()))
			}
		}
	}

	// Apply room name mapping
	for i, r := range rooms {
		if mapped, exists := ct2blrk[r]; exists {
			rooms[i] = mapped
		}
	}

	// Add finish room for calculation
	roomsWithFinish := append(rooms, "finish room")

	if BotCommandsChannelID == "" {
		BotCommandsChannelID = GetChannelIDByName("bot-commands")
		if BotCommandsChannelID == "" {
			return calc.CalcSeedResult{}, nil, fmt.Errorf("could not find #bot-commands channel")
		}
	}

	if err := checkBotPermissions(BotCommandsChannelID); err != nil {
		return calc.CalcSeedResult{}, nil, fmt.Errorf("permission error: %w", err)
	}

	// Calculate all splits
	results, err := calc.CalcSeed(roomsWithFinish)
	if err != nil {
		return calc.CalcSeedResult{}, nil, fmt.Errorf("error calculating seed: %w", err)
	}
	if len(results) == 0 {
		return calc.CalcSeedResult{}, nil, fmt.Errorf("no results found")
	}

	bestResult := results[0]

	// Filter rooms by option
	filteredRooms := []string{}
	for _, r := range rooms {
		isEasy := calc.RoomMap[r].Difficulty == calc.Easy
		isHard := calc.RoomMap[r].Difficulty == calc.Hard

		switch filter {
		case "easy":
			if isEasy {
				filteredRooms = append(filteredRooms, r)
			}
		case "hard":
			if isHard {
				filteredRooms = append(filteredRooms, r)
			}
		default:
			filteredRooms = append(filteredRooms, r)
		}
	}


	// Build BoostRoomsResponse
	boostRooms := []BoostRoomsResponse{}
	for _, room := range bestResult.BoostRooms {
		name := calc.RoomMap[roomsWithFinish[room.Ind]].Name
		strat := calc.RoomMap[roomsWithFinish[room.Ind]].BoostStrats[room.StratInd].Name
		boostRooms = append(boostRooms, BoostRoomsResponse{
			Name:     fmt.Sprintf("%s (%s)", name, strat),
			Pacelock: room.Pacelock,
			Index:    room.Ind,
		})
	}

	// Send result to Discord
	if bestResult.BoostTime < 130 && !seedCache.HasSeen(strings.Join(roomsWithFinish, "|")) {
		seedCache.MarkSeen(strings.Join(roomsWithFinish, "|"))

		img, err := drawCalcResults(roomsWithFinish, []calc.CalcSeedResult{bestResult})
		if err != nil {
			return calc.CalcSeedResult{}, nil, fmt.Errorf("error drawing seed results: %w", err)
		}

		content := fmt.Sprintf("Found a %s seed!", FormatTime(bestResult.BoostTime))
		calcCommand := createCalcCommand(filteredRooms)

		components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: ButtonShowCalc,
						Label:    "How did you get this?",
						Style:    discordgo.SuccessButton,
					},
					discordgo.Button{
						CustomID: ButtonCopyCalcCommand,
						Label:    "Copy Calc Command",
						Style:    discordgo.PrimaryButton,
						Emoji: &discordgo.ComponentEmoji{
							Name: "📋",
						},
					},
				},
			},
		}

		message, err := s.ChannelMessageSendComplex(BotCommandsChannelID, &discordgo.MessageSend{
			Content:    content,
			Components: components,
			Files: []*discordgo.File{
				{
					Name:   "seed.png",
					Reader: bytes.NewReader(img.Bytes()),
				},
			},
		})
		if err != nil {
			return calc.CalcSeedResult{}, nil, fmt.Errorf("error sending message: %w", err)
		}

		messageStates[message.ID] = &ResultState{
			Rooms:       filteredRooms,
			Results:     []calc.CalcSeedResult{bestResult},
			Index:       0,
			Filter:      ButtonAnyBoost,
			CalcCommand: calcCommand,
		}

		cleanupTimers[message.ID] = cleanupMessageState(message.ID, s, BotCommandsChannelID, true)
	}

	return bestResult, boostRooms, nil
}

// ------------------------
// Helpers (from your existing code)
// ------------------------
func GetChannelIDByName(channelName string) string {
	if s == nil {
		log.Error("Discord session is not initialized")
		return ""
	}

	if GuildID == "" {
		guilds, err := s.UserGuilds(100, "", "", false)
		if err != nil {
			log.Errorf("Error getting user guilds: %v", err)
			return ""
		}

		for _, guild := range guilds {
			channels, err := s.GuildChannels(guild.ID)
			if err != nil {
				log.Errorf("Error getting channels for guild %s: %v", guild.ID, err)
				continue
			}

			for _, channel := range channels {
				if channel.Type == discordgo.ChannelTypeGuildText && channel.Name == channelName {
					return channel.ID
				}
			}
		}
	} else {
		channels, err := s.GuildChannels(GuildID)
		if err != nil {
			log.Errorf("Error getting channels for guild %s: %v", GuildID, err)
			return ""
		}

		for _, channel := range channels {
			if channel.Type == discordgo.ChannelTypeGuildText && channel.Name == channelName {
				return channel.ID
			}
		}
	}

	log.Errorf("Channel '%s' not found", channelName)
	return ""
}

func checkBotPermissions(channelID string) error {
	log.Info("checking bot permissions")

	if s == nil {
		err := fmt.Errorf("discord session is not initialized")
		log.Error(err)
		return err
	}

	_, err := s.Channel(channelID)
	if err != nil {
		err := fmt.Errorf("error getting channel info: %w", err)
		log.Error(err)
		return err
	}

	permissions, err := s.State.UserChannelPermissions(s.State.User.ID, channelID)
	if err != nil {
		err := fmt.Errorf("error getting permissions: %w", err)
		log.Error(err)
		return err
	}

	requiredPerms := discordgo.PermissionViewChannel |
		discordgo.PermissionSendMessages |
		discordgo.PermissionAttachFiles

	if permissions&int64(requiredPerms) != int64(requiredPerms) {
		return fmt.Errorf("bot lacks necessary permissions for channel %s. Has: %d, Needs: %d",
			channelID, permissions, requiredPerms)
	}

	return nil
}

func createCalcCommand(rooms []string) string {
	var commandParts []string
	commandParts = append(commandParts, "/calc")

	for i, room := range rooms {
		roomName := room
		if len(roomName) > 0 {
			roomName = strings.ToUpper(roomName[:1]) + roomName[1:]
		}
		commandParts = append(commandParts, fmt.Sprintf("room_%d:%s", i+1, roomName))
	}

	return strings.Join(commandParts, " ")
}
