# Development Guide

## Setup Development Environment

### Prerequisites
- Go 1.23+
- Task (task runner)
- Wire (dependency injection)
- golangci-lint

### Installation
1. Clone repository:
```bash
git clone https://github.com/jonesrussell/dashboard.git
cd dashboard
```

2. Install dependencies:
```bash
task deps
```

3. Install development tools:
```bash
task install-tools
```

## Development Workflow

### Building
```bash
task build        # Build application
task run          # Build and run
task run-debug    # Run with debug output
task run-external # Run in new window
```

### Testing
```bash
task test   # Run tests with coverage
task bench  # Run benchmarks
task lint   # Run linter
```

### Code Generation
```bash
task wire   # Generate dependency injection code
```

## Widget Development

### Creating a New Widget

1. Create widget package:
```
internal/ui/widgets/mywidget/
├── widget.go      # Widget implementation
└── widget_test.go # Tests
```

2. Implement Widget interface:
```go
type Widget struct {
    components.BaseWidget
    // Widget-specific fields
}

func New() *Widget {
    return &Widget{}
}

// Implement required methods
func (w *Widget) Init() tea.Cmd { ... }
func (w *Widget) Update(msg tea.Msg) (Widget, tea.Cmd) { ... }
func (w *Widget) View() string { ... }
```

### Best Practices

#### Time Handling
- Always check for both nil and zero time values
- Use pointer to time.Time for optional timestamps
- Validate time values before using them
- Use time.IsZero() for zero time checks
- Document time-related edge cases

#### State Management
- Keep widget state minimal
- Use proper types for state
- Handle all state transitions

#### Error Handling
- Use error wrapping
- Provide user feedback
- Log errors appropriately

#### Performance
- Cache expensive computations
- Minimize allocations
- Profile when needed

#### Testing
- Test all state transitions
- Mock external dependencies
- Use test utilities

## Code Style

### General Guidelines
- Follow Go conventions
- Document exported symbols
- Keep functions focused
- Use meaningful names

### Logging
- Use structured logging
- Include context
- Appropriate log levels

### Error Handling
- Return errors, don't panic
- Wrap errors with context
- Log at appropriate level

## Pull Request Process

1. Create feature branch
2. Write tests
3. Update documentation
4. Run linter
5. Submit PR

## Release Process

1. Update version
2. Update changelog
3. Run tests
4. Create release
5. Update documentation 