# Getting Started (developers)

## Quick Start
To run the project, please put the binary or cwd in {root project}/bin/
Overall, either for dev/prod the project looks like this:

```text
bin/
    web_app (binary)

public/
    ...

pages/
    *.tmpl

config.yml (config file)
```

launch.json example for vscode users:
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "dev",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/app/main.go",
            "output": "${workspaceFolder}/bin/debug_web",
            "cwd": "${workspaceFolder}/bin/",
            "args": [
                "--env=development",
                "--cfg=../config.yml"
            ]
        },
        {
            "name": "prod",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/app/main.go",
            "output": "${workspaceFolder}/bin/debug_web",
            "cwd": "${workspaceFolder}/bin/",
            "args": [
                "--env=production",
                "--cfg=../config.yml"
            ]
        }
    ]
}
```


## Development Dependencies

Before starting development, please install the following dependencies:

1. **Mockery v3** - For auto-generating mocks
   - Install from: https://vektra.github.io/mockery/v3/
   - Used for generating mock objects in tests

2. **YQ v4** - YAML processor used in Makefile
   - Install from: https://github.com/mikefarah/yq/releases
   - Required for parsing config.yml (used in makefile)

3. **Goose** - Database migration tool
   - Install from: https://github.com/pressly/goose
   - Used for managing SQLite database migrations

4. **GCC Compiler**
   - Required for compiling the sqlite3 package (https://github.com/mattn/go-sqlite3)
   - Install using your system's package manager:
     - Ubuntu/Debian: `sudo apt install gcc`
     - macOS: Install Xcode Command Line Tools
     - Windows: Install MinGW or MSYS2