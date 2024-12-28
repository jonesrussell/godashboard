# Terminal Dashboard

A modern terminal user interface (TUI) dashboard built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- 🎨 Beautiful terminal UI styling with Lipgloss
- 🔄 Interactive components with Bubbletea
- ⌨️ Intuitive keyboard controls
- 📊 Modular widget system
- 💡 Built-in help system
- 📝 Structured logging with Zap
- 🔌 Dependency injection with Wire

## Prerequisites

- Go 1.23 or higher
- Task (taskfile.dev) for development commands

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/dashboard
cd dashboard

# Install development tools
task install-tools

# Install dependencies
task deps
```

## Development Commands

```bash
# Run the application
task run

# Run in external window
task run-external

# Format code
task fmt

# Run linter
task lint

# Run tests
task test

# Run all checks
task all

# Clean build artifacts
task clean

# Watch for changes
task watch
```

## Project Structure

```
.
├── cmd/
│   └── dashboard/        # Main application entry point
├── internal/
│   ├── ui/              # UI components and layouts
│   │   ├── styles/      # Lipgloss styles
│   │   └── components/  # Reusable UI components
│   └── logger/          # Structured logging system
├── build/               # Build artifacts
└── coverage/            # Test coverage reports
```

## Keyboard Controls

- `q` or `ctrl+c` - Quit the application
- `?` - Toggle help menu
- `tab` - Navigate between widgets
- `enter` - Select/activate current widget

## Development

### Adding New Widgets

To create a new widget, implement the `Widget` interface in `internal/ui/components/widget.go`:

```go
type Widget interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Widget, tea.Cmd)
    View() string
    SetSize(width, height int)
}
```

### Logging

The application uses [uber-go/zap](https://github.com/uber-go/zap) for structured logging with the following features:

- Multiple log levels (Debug, Info, Warn, Error, Fatal)
- Structured field-based logging
- Automatic log rotation
- Request ID tracking
- Context-aware logging

### Dependency Injection

We use [google/wire](https://github.com/google/wire) for compile-time dependency injection:

- Providers are defined in `provider.go` files
- Wire is automatically run before builds
- See `cmd/dashboard/wire.go` for the main injection setup

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for the amazing Bubbletea and Lipgloss libraries
- [Uber](https://github.com/uber-go/zap) for the Zap logging library
- [Google](https://github.com/google/wire) for the Wire dependency injection tool
- The Go community for inspiration and support
