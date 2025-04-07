# Getting Started (developers)

To run the project, please put the binary or cwd in {root project}/bin/
Overall, either for dev/prod the project looks like this:

```
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
            "program": "${workspaceFolder}/cmd/web/main.go",
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
            "program": "${workspaceFolder}/cmd/web/main.go",
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