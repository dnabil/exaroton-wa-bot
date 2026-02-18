package command

import (
	"context"
	"exaroton-wa-bot/internal/service"
)

type (
	Command interface {
		Name() string
		Help() string
		Usage() string
		Execute(c context.Context, args []string) CommandResult
	}

	Registry struct {
		commands map[string]Command
	}

	CommandResult struct {
		Text  string
		Error error
	}
)

func NewRegistry(WhatsappService service.IWhatsappService, serverSettingsSvc service.IServerSettingsService) *Registry {
	r := &Registry{
		commands: make(map[string]Command),
	}

	// register commands here...
	r.Register(NewHelpCommand(r))
	r.Register(NewListServerCommand(serverSettingsSvc))
	r.Register(NewStartServerCommand(serverSettingsSvc))
	r.Register(NewInfoCommand(serverSettingsSvc))
	r.Register(NewStopServerCommand(serverSettingsSvc))

	return r
}

// Register a new command to the registry.
// The command is identified by its Name method.
func (r *Registry) Register(cmd Command) {
	_, ok := r.commands[cmd.Name()]
	if ok {
		panic("command name must be unique")
	}

	r.commands[cmd.Name()] = cmd
}

// List all registered commands.
func (r *Registry) List() []Command {
	out := make([]Command, 0, len(r.commands))
	for _, c := range r.commands {
		out = append(out, c)
	}
	return out
}

// Get a command by its name.
// Returns the command and whether it was found.
func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}
