package command

import (
	"context"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/helper"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"strconv"
)

var (
	InfoCmdName = "info"
)

var _ Command = new(InfoCommand)

type InfoCommand struct {
	serverSettingsSvc service.IServerSettingsService
}

func NewInfoCommand(serverSettingsSvc service.IServerSettingsService) *InfoCommand {
	return &InfoCommand{
		serverSettingsSvc: serverSettingsSvc,
	}
}

func (c *InfoCommand) Name() string {
	return InfoCmdName
}

func (c *InfoCommand) Help() string {
	return "Check a server info by its ID"
}

func (c *InfoCommand) Usage() string {
	return "/info [id]"
}

func (c *InfoCommand) Execute(ctx context.Context, args []string) CommandResult {
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

	server, err := c.serverSettingsSvc.GetExarotonServerInfo(ctx, uint(serverIdx))
	if err != nil {
		return CommandResult{
			Error: err,
		}
	}

	return CommandResult{
		Text: turnServerInfoIntoText(server, uint(serverIdx)),
	}
}

func turnServerInfoIntoText(server *dto.ExarotonServerInfo, serverIdx uint) string {
	return fmt.Sprintf("ID: %d [%s] \nName: %s\nAddress: %s\nMotd: %s\nStatus: %s\nHost: %s\nPort: %s\nPlayers:%d/%d\nSoftware: %s %s\nShared: %t",
		serverIdx,
		server.ID,
		server.Name,
		server.Address,
		server.Motd,
		server.Status,
		helper.Deref(server.Host),
		helper.If(helper.IsZero(server.Port), "", strconv.Itoa(helper.Deref(server.Port))),
		server.Players.Count,
		server.Players.Max,
		server.Software.Name,
		server.Software.Version,
		server.Shared)
}
