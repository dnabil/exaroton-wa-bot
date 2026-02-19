package command

import (
	"context"
	"exaroton-wa-bot/internal/constants/messages"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"strconv"
)

var (
	ListServerCmdName = "servers"
)

var _ Command = new(ListServerCommand)

type ListServerCommand struct {
	serverSettingsSvc service.IServerSettingsService
}

func NewListServerCommand(serverSettingsSvc service.IServerSettingsService) *ListServerCommand {
	return &ListServerCommand{
		serverSettingsSvc: serverSettingsSvc,
	}
}

func (c *ListServerCommand) Name() string {
	return ListServerCmdName
}

func (c *ListServerCommand) Help() string {
	return "Show available servers"
}

func (c *ListServerCommand) Usage() string {
	return "/servers [page]"
}

func (c *ListServerCommand) Execute(ctx context.Context, args []string) CommandResult {
	servers, err := c.serverSettingsSvc.ListExarotonServer(ctx)
	if err != nil {
		return CommandResult{Error: err}
	}

	// pagination
	var (
		page       = 1
		limit      = 4
		totalItems = len(servers)
	)
	if len(args) > 0 {
		page, err = strconv.Atoi(args[0])
		if err != nil {
			return CommandResult{Error: err}
		}
	}

	pag := dto.NewPagination(page, limit, totalItems)

	text := fmt.Sprintf(messages.CmdShowingPage, c.Name(), pag.CurrentPage, pag.TotalPage) + "\n\n"
	text += c.formatServersIntoText(servers[pag.Start():pag.End()])

	return CommandResult{Text: text}
}

func (c *ListServerCommand) formatServersIntoText(servers []*dto.ExarotonServerInfo) string {
	var text string
	for i, srv := range servers {
		text += formatServerIntoText(i, srv) + "\n"
	}

	if text == "" {
		text = "No servers found"
	}

	return text
}

func formatServerIntoText(i int, srv *dto.ExarotonServerInfo) string {
	return fmt.Sprintf("ID: %d [%s]\nName: %s\nAddress: %s\nStatus: %s\nSoftware: %s %s\n",
		i,
		srv.ID,
		srv.Name,
		srv.Address,
		srv.Status,
		srv.Software.Name,
		srv.Software.Version)
}
