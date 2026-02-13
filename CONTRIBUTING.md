# Contributing

Contributions are welcome! Here's how to get started.

## Development Setup

**Requirements:** Go 1.25+

```sh
git clone https://github.com/ggfevans/endorse.git
cd endorse
make build
make test
```

## Code Style

- Run `gofmt` on all Go files
- Run `go vet ./...` before submitting
- Follow existing patterns in the codebase

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-change`)
3. Make your changes
4. Run tests: `make test`
5. Run vet: `make vet`
6. Commit with a clear message
7. Open a pull request against `main`

## Issue Guidelines

- Check existing issues before opening a new one
- Include steps to reproduce for bug reports
- For feature requests, describe the use case

## Project Structure

```
cmd/endorse/         Main entry point
internal/
  app/               Root application model and update loop
  config/            Configuration and credential storage
  linkedin/          LinkedIn API client (wraps mautrix-linkedin)
  ui/
    compose/         Message compose textarea
    convlist/        Conversation list panel
    header/          Top header bar
    layout/          Layout calculations
    modal/           Auth and confirm modals
    statusbar/       Bottom status bar
    styles/          Theme and style definitions
    thread/          Message thread panel
  util/              Text and time utilities
```
