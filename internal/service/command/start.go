package command

import (
	"context"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/constants/messages"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"strconv"
	"time"
)

var (
	StartServerCmdName = "start"
)

var _ Command = new(StartServerCommand)

type StartServerCommand struct {
	serverSettingsSvc service.IServerSettingsService
}

func NewStartServerCommand(serverSettingsSvc service.IServerSettingsService) *StartServerCommand {
	return &StartServerCommand{
		serverSettingsSvc: serverSettingsSvc,
	}
}

func (c *StartServerCommand) Name() string {
	return StartServerCmdName
}

func (c *StartServerCommand) Help() string {
	return "Start a server by its ID"
}

func (c *StartServerCommand) Usage() string {
	return "/start [id]"
}

func (c *StartServerCommand) Execute(ctx context.Context, args []string) CommandResult {
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

	startStatus := c.serverSettingsSvc.StartExarotonServer(ctx, uint(serverIdx), service.WithPolling(
		50*time.Second,
		10*time.Second,
	))
	if startStatus.Err != nil {
		return CommandResult{
			Error: err,
		}
	}

	lastStatus := dto.ServerStatusStarting
	for v := range startStatus.Status {
		lastStatus = v
	}

	return CommandResult{
		Text: fmt.Sprintf(messages.ServerStartFinish, serverIdx, lastStatus.String()),
	}
}
