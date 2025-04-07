package config

import (
	"exaroton-wa-bot/internal/errs"
	"fmt"
	"log/slog"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Cfg struct {
	*koanf.Koanf
	isYmlSet bool

	Args *Args
}

// config keys
const (
	// app's port
	KeyPort = "port"

	// html pages dir
	KeyPagesDir  = "pages_dir"  // string
	KeyPublicDir = "public_dir" // string

	keyLogLevel  = "log_level"
	keyIsJsonLog = "is_json_log" // bool
)

// log keys
const (
	KeyLogErr      = "error"
	KeyLogErrStack = "error_stack"
	// only used when panic happens.
	keyLogPanicID = "panic_id"

	// status request (int)
	KeyLogStatus = "status"

	// uri request (string)
	KeyLogURI = "uri"
)

// app environment
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type Args struct {
	CfgPath string
	Env     string
}

func InitConfig(args *Args) (*Cfg, error) {
	cfg := &Cfg{
		Koanf: koanf.New("."),
		Args:  args,
	}

	if err := cfg.LoadYmlConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Cfg) LoadYmlConfig() error {
	if cfg.Args.CfgPath == "" {
		return errs.ErrCfgYmlPathNotSet
	}

	cfgProvider := file.Provider(cfg.Args.CfgPath)
	cfgParser := yaml.Parser()

	err := cfg.Load(cfgProvider, cfgParser)
	if err != nil {
		return err
	}

	err = cfg.watchCfgFile(cfgProvider, ".yml", cfgParser)
	if err != nil {
		return err
	}

	cfg.isYmlSet = true
	return nil
}

// watchCfgFile watches a configuration file for changes and reloads the configuration if necessary.
//
// Parameters:
// - f: a pointer to a file.File object representing the configuration file.
// - cfgName: a string representing the name of the configuration file.
// - pa: a koanf.Parser object used to parse the configuration file.
// - opts: optional koanf.Option objects for additional configuration options.
//
// Returns:
// - error: an error object if there was an error watching the configuration file or reloading the configuration.
func (cfg *Cfg) watchCfgFile(f *file.File, cfgName string, pa koanf.Parser, opts ...koanf.Option) error {
	if !cfg.isYmlSet {
		if err := f.Watch(func(event interface{}, err error) {
			if err != nil {
				slog.Error(fmt.Sprintf("init watch config file (%s) error", cfgName), KeyLogErr, err)
				return
			}

			err = cfg.Load(f, pa, opts...)
			if err != nil {
				slog.Error(fmt.Sprintf("reloading config file (%s) error", cfgName), KeyLogErr, err)
				return
			}

			slog.Info("reloaded cfg: " + cfgName)
		}); err != nil {
			return err
		}
	}

	return nil
}
