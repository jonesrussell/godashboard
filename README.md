# Go Dashboard

A terminal-based system dashboard built with Go, featuring real-time system monitoring, todo management, and process monitoring.

## Features

- Real-time system monitoring
  - CPU usage with visual progress bars
  - Memory usage tracking
  - Disk space monitoring
- Grid-based widget layout
  - Flexible positioning
  - Dynamic resizing
  - Focus management
- Modern terminal UI
  - Smooth updates
  - Keyboard navigation
  - Color themes
- Performance optimized
  - Style caching
  - Content caching
  - Minimal allocations
- Structured logging
  - Zap-based logging
  - Log rotation
  - Debug output
- Comprehensive testing
  - Unit tests
  - Integration tests
  - Performance benchmarks
  - Test utilities

## Prerequisites

- Go 1.23 or later
- Task (task-based build tool)
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/jonesrussell/dashboard.git
cd dashboard
```

2. Install dependencies:
```bash
task deps
```

3. Build the project:
```bash
task build
```

## Usage

Run the dashboard:
```bash
task run
```

Run with debug output:
```bash
task run-debug
```

Run in external window:
```bash
task run-external
```

### Keyboard Controls

- `Tab` - Navigate between widgets
- `Enter` - Select/activate widget
- `q` or `Ctrl+C` - Quit
- `?` - Toggle help

## Development

### Requirements

- Go 1.23+
- Task
- Wire (dependency injection)
- golangci-lint

### Setup Development Environment

1. Install development tools:
```bash
task setup
```

2. Run tests:
```bash
task test
```

3. Run linter:
```bash
task lint
```

4. Run benchmarks:
```bash
task bench
```

### Project Structure

```
.
├── cmd/
│   └── dashboard/     # Main application
├── internal/
│   ├── logger/        # Logging package
│   ├── testutil/      # Test utilities
│   └── ui/           # User interface
│       ├── components/  # UI components
│       ├── container/   # Widget container
│       ├── styles/      # UI styling
│       └── widgets/     # Dashboard widgets
├── pkg/              # Public packages
└── test/            # Test utilities
```

### Testing

The project uses Go's testing framework with additional utilities:

- `testutil.NewTestLogger` - Creates a logger for tests
- `testutil.ReadLogFile` - Reads and verifies log output
- `testutil.NewUITest` - Helps test UI components

Example test:
```go
func TestMyFeature(t *testing.T) {
    log, logPath := testutil.NewTestLogger(t, "test-name")
    
    // Use the logger in your test
    log.Info("Test started")
    
    // Verify log output
    content, err := testutil.ReadLogFile(logPath)
    require.NoError(t, err)
    assert.Contains(t, content, "Test started")
}
```

### Logging

The dashboard uses structured logging with Zap:

- Debug output goes to `logs/dashboard-debug.log`
- Log rotation is configured
- Test logs are automatically cleaned up
- Log levels: debug, info, warn, error

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

### Code Style

- Follow Go conventions
- Use provided test utilities
- Add tests for new features
- Document exported symbols
- Run linter before committing

## License

This project is licensed under the MIT License - see the LICENSE file for details.

