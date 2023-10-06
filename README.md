# Course Backend Master Class

This repo contains source code produced during the ambitious course "Backend Master Class".

## Course objectives

Learn everything about backend web development: Golang, Postgres, Redis, Gin, gRPC, Docker, Kubernetes, AWS, CI/CD

## Content (so far)

* Docker, https://www.docker.com
  * multi stage builds
  * Docker in Docker
* golang, https://go.dev
  * Viper, https://github.com/spf13/viper
  * golang-migrate, https://github.com/golang-migrate/migrate
  * testcontainers-go, https://github.com/testcontainers/testcontainers-go
  * sqlc, https://sqlc.dev
* PostgreSQL, https://www.postgresql.org

## VSCode setup

settings.json

~~~json
{
    "go.toolsManagement.autoUpdate": true,
    "explorer.compactFolders": false,
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.alternateTools": {
        "go-langserver": "gopls", 
      },
    "go.languageServerFlags": [],
    "gopls": {
      "ui.semanticTokens": true,
      "formatting.gofumpt": true
    },
    "workbench.colorTheme": "Ayu Dark",
    "editor.fontLigatures": true,
    "editor.fontFamily": "JetBrains Mono",
    "editor.fontWeight": 500,
    "editor.fontSize": 13,
    "runOnSave.commands": [
      {
          "match": ".*\\.go$",
          "command": "golines ${file} -m 120 -w --ignore-generated --no-reformat-tags",
      },
  ],
  "editor.semanticTokenColorCustomizations": {},
  "editor.codeActionsOnSave": {},
  "go.buildFlags": [
    "-tags=unit,integration"
],
"go.testTags": "unit,integration",
}


~~~

# Some reading

Collected some other resources during the course.

https://eltonminetto.dev/en/post/2022-10-22-creating-api-using-go-sqlc/

# Troubleshooting tips

Update VS Code Go Tools. Command + Shift + P -> Go: Install/update tools Install all tools and restart VS Code.

