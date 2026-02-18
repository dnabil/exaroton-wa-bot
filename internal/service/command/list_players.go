package command

import (
	"context"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"strconv"
)

var (
	ListPlayersCmdName = "players"
)

var _ Command = new(ListPlayersCommand)

type ListPlayersCommand struct {
	serverSettingsSvc service.IServerSettingsService
}

func NewListPlayersCommand(serverSettingsSvc service.IServerSettingsService) *ListPlayersCommand {
	return &ListPlayersCommand{
		serverSettingsSvc: serverSettingsSvc,
	}
}

func (c *ListPlayersCommand) Name() string {
	return ListPlayersCmdName
}

func (c *ListPlayersCommand) Help() string {
	return "List players on a server by its ID"
}

func (c *ListPlayersCommand) Usage() string {
	return "/players [id]"
}

func (c *ListPlayersCommand) Execute(ctx context.Context, args []string) CommandResult {
	if len(args) == 0 {
		return CommandResult{Error: errs.ErrCommandMissingArg}
	}

	var (
		serverIdx int
		err       error
	)
	if serverIdx, err = strconv.Atoi(args[0]); err != nil {
		return CommandResult{
			Error: errs.ErrCommandInvalidArg,
		}
	}

	playerList, err := c.serverSettingsSvc.GetExarotonServerPlayerList(ctx, uint(serverIdx))
	if err != nil {
		return CommandResult{
			Error: err,
		}
	}

	return CommandResult{
		Text: c.formatPlayerListToText(uint(serverIdx), playerList),
	}
}

func (c *ListPlayersCommand) formatPlayerListToText(serverId uint, playerList *dto.ExarotonServerPlayers) string {
	if len(playerList.List) == 0 {
		return "No players online."
	}

	result := fmt.Sprintf("[ServerID: %d] Players online(%d):\n", serverId, len(playerList.List))
	for i, player := range playerList.List {
		result += fmt.Sprintf("%d. %s\n", i+1, player)
	}

	return result
}
