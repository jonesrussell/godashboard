# Development Guide

This guide provides information for developers who want to contribute to the Go Dashboard project.

## Development Environment

### Prerequisites

- Go 1.23 or later
- Task (task-based build tool)
- Git
- MinGW-w64 (for Windows builds)
- GoReleaser (for releases)
- golangci-lint (for code quality)
- Wire (for dependency injection)

### Setup

1. Install development tools:
```bash
task install-tools
```

2. Install dependencies:
```bash
task deps
```

3. Run tests:
```bash
task test
```

## Project Structure

```
.
├── cmd/
│   └── dashboard/     # Main application
├── docs/             # Documentation
│   ├── ARCHITECTURE.md
│   ├── API.md
│   └── DEVELOPMENT.md
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

## Build System

The project uses Task for build automation. Common tasks include:

- `task build` - Build the application
- `task test` - Run tests
- `task lint` - Run linter
- `task bench` - Run benchmarks
- `task run` - Run the application
- `task clean` - Clean build artifacts

### Build Output Structure

The build system creates an organized release structure:
```
build/
└── windows_amd64/          # or linux_amd64
    ├── bin/
    │   └── dashboard.exe   # Main executable
    ├── docs/               # Documentation
    │   ├── images/
    │   ├── API.md
    │   ├── ARCHITECTURE.md
    │   └── DEVELOPMENT.md
    ├── LICENSE
    └── README.md
```

## Testing

### Test Categories

1. Unit Tests
   - Widget behavior
   - State transitions
   - Message handling
   - Style application

2. Integration Tests
   - Widget interactions
   - Focus management
   - Layout behavior
   - Event propagation

3. Benchmark Tests
   - View rendering
   - Update performance
   - Memory allocations
   - Style operations

### Running Tests

```bash
# Run all tests
task test

# Run benchmarks
task bench

# Run tests with coverage
task test-coverage
```

## Release Process

### Local Testing

1. Test the release process locally:
```bash
task release-dry-run
```

This will:
- Build binaries for all platforms
- Create release archives
- Generate checksums
- Skip publishing

2. Verify the contents in `dist/`:
- Check binary functionality
- Verify documentation
- Validate archive structure

### Creating a Release

1. Update version information:
- Update version in code
- Update changelog
- Commit changes

2. Create and push a tag:
```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

3. The GitHub Actions workflow will automatically:
- Build binaries for all platforms
- Create release archives with documentation
- Generate checksums
- Create a GitHub release
- Upload artifacts

### Release Structure

Each release includes:
- Platform-specific binaries
- Complete documentation
- License information
- Checksums for verification

## Code Style

### General Guidelines

1. Follow Go conventions
2. Document exported symbols
3. Use provided test utilities
4. Add tests for new features
5. Run linter before committing

### Widget Development

1. Implement required interfaces:
   - Widget interface
   - Focusable interface (if needed)
   - Sizable interface (if needed)

2. Follow widget patterns:
   - Use composition over inheritance
   - Implement proper cleanup
   - Handle focus correctly
   - Manage dimensions

3. Add tests:
   - Basic functionality
   - Edge cases
   - Performance benchmarks

## Debugging

### Debug Mode

Run with debug output:
```bash
task run-debug
```

This enables:
- Detailed logging
- Performance metrics
- Widget state information
- Error details

### Logging

- Use structured logging with zap
- Include context in log entries
- Use appropriate log levels
- Add request IDs for tracking

## Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Add tests
5. Update documentation
6. Submit a pull request

### Pull Request Process

1. Ensure all tests pass
2. Update relevant documentation
3. Add entry to changelog
4. Get review from maintainers
5. Address feedback

## Support

- GitHub Issues for bug reports
- Discussions for questions
- Pull requests for contributions 