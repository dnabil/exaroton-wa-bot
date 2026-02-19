package command

import (
	"context"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/constants/messages"
	"exaroton-wa-bot/internal/service"
	"strconv"
)

var (
	StopServerCmdName = "stop"
)

var _ Command = new(StopServerCommand)

type StopServerCommand struct {
	serverSettingsSvc service.IServerSettingsService
}

func NewStopServerCommand(serverSettingsSvc service.IServerSettingsService) *StopServerCommand {
	return &StopServerCommand{
		serverSettingsSvc: serverSettingsSvc,
	}
}

func (c *StopServerCommand) Name() string {
	return StopServerCmdName
}

func (c *StopServerCommand) Help() string {
	return "Stop a server by its ID"
}

func (c *StopServerCommand) Usage() string {
	return "/stop [id]"
}

func (c *StopServerCommand) Execute(ctx context.Context, args []string) CommandResult {
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

	if err = c.serverSettingsSvc.StopExarotonServer(ctx, uint(serverIdx)); err != nil {
		return CommandResult{
			Error: err,
		}
	}

	return CommandResult{
		Text: messages.ServerIsStopping,
	}
}
