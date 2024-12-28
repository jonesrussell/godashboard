# Architecture Overview

## Core Components

### Dashboard
The main application container that manages the overall UI layout and widget coordination.

### Widget System
- **Base Widget**: Provides common functionality for all widgets
  - Focus management
  - Size constraints
  - Style handling
  - Event propagation

### Component Hierarchy
```
Dashboard
├── Container
│   ├── System Info Widget
│   │   ├── CPU Monitor
│   │   ├── Memory Monitor
│   │   └── Disk Monitor
│   └── Tasks Widget
│       ├── Task List
│       └── Task Controls
└── Help System
```

## Key Design Patterns

### Model-View-Update (MVU)
The UI follows the Elm architecture pattern through Bubbletea:
1. **Model**: Represents application state
2. **View**: Renders state to terminal UI
3. **Update**: Handles messages and updates state

### Dependency Injection
Uses Google Wire for compile-time dependency injection:
- Promotes testability
- Manages component lifecycle
- Reduces coupling

### Widget Interface
All widgets implement the core Widget interface:
```go
type Widget interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Widget, tea.Cmd)
    View() string
    SetSize(width, height int)
    GetDimensions() (width, height int)
    Focus()
    Blur()
    IsFocused() bool
}
```

## Data Flow

### Event Handling
1. User input or system event occurs
2. Event propagates through widget hierarchy
3. Appropriate widget handles event
4. State updates trigger re-render

### API Integration
- RESTful API client for task management
- Environment-based configuration
- Proper error handling and retry logic

## Performance Considerations

### Caching
- Style caching for consistent performance
- Content caching to reduce calculations
- Size caching for layout stability

### Memory Management
- Minimal allocations in render path
- Efficient string building
- Smart update batching

## Testing Strategy

### Levels
1. Unit Tests: Individual components
2. Integration Tests: Component interaction
3. Performance Tests: Benchmarks

### Utilities
- Test logger configuration
- UI testing helpers
- Assertion utilities 