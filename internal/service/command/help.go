package command

import (
	"context"
	"exaroton-wa-bot/internal/constants/errs"
	"exaroton-wa-bot/internal/constants/messages"
	"exaroton-wa-bot/internal/dto"
	"fmt"
	"sort"
	"strconv"
)

var (
	HelpCmdName = "help"
)

var _ Command = new(HelpCommand)

type HelpCommand struct {
	registry *Registry
}

func NewHelpCommand(r *Registry) *HelpCommand {
	return &HelpCommand{registry: r}
}

func (c *HelpCommand) Name() string {
	return HelpCmdName
}

func (c *HelpCommand) Help() string {
	return "Show available commands"
}

func (c *HelpCommand) Usage() string {
	return "/help [page|command]"
}

func (c *HelpCommand) Execute(ctx context.Context, args []string) CommandResult {
	cmds := c.registry.List()

	if len(args) == 0 {
		return c.showPage(ctx, cmds, 1)
	}

	// pagination
	if page, err := strconv.Atoi(args[0]); err == nil {
		return c.showPage(ctx, cmds, page)
	}

	// command detail
	return c.showCommandDetail(ctx, args[0])
}

func (c *HelpCommand) showPage(ctx context.Context, cmds []Command, page int) CommandResult {
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Name() < cmds[j].Name()
	})

	pag := dto.NewPagination(page, 7, len(cmds))

	msg := fmt.Sprintf(messages.CmdShowingPage, c.Name(), pag.CurrentPage, pag.TotalPage) + "\n\n"

	for _, cmd := range cmds[pag.Start():pag.End()] {
		msg += fmt.Sprintf(
			"/%-10s %s\n",
			cmd.Name(),
			cmd.Help(),
		)
	}

	return CommandResult{Text: msg}
}

func (c *HelpCommand) showCommandDetail(ctx context.Context, name string) CommandResult {
	cmd, ok := c.registry.Get(name)
	if !ok {
		return CommandResult{Error: errs.ErrCommandNotFound}
	}

	msg := fmt.Sprintf(
		"%s\n\n%s\n\nUsage:\n%s",
		cmd.Name(),
		cmd.Help(),
		cmd.Usage(),
	)

	return CommandResult{Text: msg}
}
