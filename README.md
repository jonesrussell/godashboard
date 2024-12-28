# Terminal Dashboard

A modern terminal user interface (TUI) dashboard built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- 🎨 Beautiful terminal UI styling with Lipgloss
- 🔄 Interactive components with Bubbletea
- ⌨️ Intuitive keyboard controls
- 📊 Modular widget system
- 💡 Built-in help system

## Prerequisites

- Go 1.23 or higher

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/dashboard
cd dashboard

# Install dependencies
go mod tidy
```

## Running the Dashboard

```bash
go run cmd/dashboard/main.go
```

## Keyboard Controls

- `q` or `ctrl+c` - Quit the application
- `?` - Toggle help menu
- `tab` - Navigate between widgets
- `enter` - Select/activate current widget

## Project Structure

```
.
├── cmd/
│   └── dashboard/        # Main application entry point
├── internal/
│   ├── ui/              # UI components and layouts
│   │   ├── styles/      # Lipgloss styles
│   │   └── components/  # Reusable UI components
│   └── models/          # Data models
└── pkg/                 # Public packages
```

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
- The Go community for inspiration and support
