package config

import (
	"errors"
	"exaroton-wa-bot/internal/constants/errs"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type Cfg struct {
	*koanf.Koanf
	isYmlSet bool

	Args *Args
}

// config keys
const (
	// app's port
	KeyPort = "app.port"

	// html pages dir
	KeyPagesDir  = "app.pages_dir"  // string
	KeyPublicDir = "app.public_dir" // string

	KeySessionSecret     = "app.session_secret"      // string
	KeyAuthDuration      = "app.auth_duration"       // string (time.Duration)
	KeyAutoWhatsappLogin = "app.auto_whatsapp_login" // bool

	keyLogLevel  = "app.log_level"   // string
	keyIsJsonLog = "app.is_json_log" // bool

	keyDBDialect    = "db.dialect"        // string
	keySQLiteDBPath = "db.sqlite_db_path" // string
	KeyDBDisableLog = "db.disable_log"    // bool
	keyDBLogLevel   = "db.log_level"      // string

	keyWADBDialect      = "wa-db.dialect"          // string
	keyWASQLiteDBPath   = "wa-db.sqlite_db_path"   // string
	keyWADBLogLevel     = "wa-db.log_level"        // string
	keyWAClientLogLevel = "wa-db.client_log_level" // string
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

	// method request (string)
	KeyLogMethod = "method"
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

	RegisterGobs()

	if err := cfg.LoadYmlConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ParseFlags() (args *Args, err error) {
	args = &Args{}

	pflag.StringVarP(&args.Env, "env", "e", EnvProduction, "App's environment (production || development). Default: production")
	pflag.StringVarP(&args.CfgPath, "cfg", "c", "../config.yml", "Path to the config file (optional) e.g: ../config.yml")
	pflag.Parse()

	args.CfgPath, err = filepath.Abs(args.CfgPath)
	if err != nil {
		return nil, errors.New("invalid config path: " + args.CfgPath)
	}

	return args, nil
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
