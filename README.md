# Minesweeper

A classic Minesweeper game implementation in Go using the [Ebiten](https://ebiten.org/) 2D game engine.

![Minesweeper Game](https://img.shields.io/badge/Go-1.22+-blue.svg)
![Ebiten](https://img.shields.io/badge/Ebiten-v2.8.8-green.svg)

## Features

### 🎮 Classic Gameplay

- Traditional Minesweeper rules and mechanics
- Mine-free first click guarantee (3x3 safe area)
- Automatic chain reaction for empty cells
- Flag system for marking suspected mines
- Win/lose detection with visual feedback

### 🏆 Three Difficulty Levels

- **初級 (Beginner)**: 9×9 grid with 10 mines
- **中級 (Intermediate)**: 16×16 grid with 40 mines
- **上級 (Expert)**: 30×16 grid with 99 mines

### 🎨 Enhanced Graphics

- Large, color-coded numbers (1-8) for mine counts
- Custom bomb graphics with spikes for realistic appearance
- Red flag symbols for marked cells
- Responsive header layout that adapts to window size
- Adaptive window scaling based on difficulty level

### ⌨️ Intuitive Controls

- **Left Click**: Open cell
- **Right Click**: Flag/unflag cell
- **R Key**: Restart current difficulty
- **1 Key**: Switch to Beginner level
- **2 Key**: Switch to Intermediate level
- **3 Key**: Switch to Expert level

### 📊 Game Information

- Real-time mine counter
- Game timer (starts on first click)
- Current difficulty display
- Win/lose status with restart instructions

## Installation

### Prerequisites

- Go 1.22 or higher
- Git

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/ebiten-tour.git
cd ebiten-tour

# Download dependencies
go mod tidy

# Build the game
go build -o minesweeper

# Run the game
./minesweeper
```

### Quick Run

```bash
go run main.go
```

## Game Rules

1. **Objective**: Clear all cells that don't contain mines
2. **Numbers**: Show count of mines in adjacent 8 cells
3. **Flags**: Right-click to mark suspected mine locations
4. **First Click**: Always safe - no mine in clicked cell or surrounding area
5. **Chain Reaction**: Clicking empty cell (0 adjacent mines) opens all connected empty areas
6. **Win Condition**: Open all non-mine cells
7. **Lose Condition**: Click on a mine

## Technical Details

### Architecture

- **Single-file design**: Simple, self-contained implementation
- **Dynamic sizing**: All game logic adapts to different grid dimensions
- **Responsive UI**: Header layout adjusts to window width
- **Custom graphics**: Hand-drawn mine and number displays

### Performance

- **60 FPS rendering**: Smooth gameplay experience
- **Efficient drawing**: Vector-based graphics for crisp display
- **Memory optimized**: Dynamic board allocation per difficulty

### Dependencies

- [Ebiten v2.8.8](https://ebiten.org/) - 2D game engine
- Go standard library

## Development

### Code Formatting

```bash
gofmt -w .
```

### Project Structure

```
ebiten-tour/
├── main.go           # Complete game implementation
├── go.mod           # Go module dependencies
├── go.sum           # Dependency checksums
├── CLAUDE.md        # Development documentation
└── README.md        # This file
```

### Key Components

- **Difficulty System**: Configurable grid sizes and mine counts
- **Game Logic**: Mine placement, cell opening, win/lose detection
- **Rendering Engine**: Custom drawing functions for numbers, mines, flags
- **Input Handling**: Mouse and keyboard event processing
- **UI Layout**: Responsive header with game information

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards (`gofmt`, `go vet`)
- Maintain the single-file architecture for simplicity
- Test changes across all three difficulty levels
- Ensure responsive behavior for different window sizes

## Screenshots

### Beginner Level (9x9)

Compact layout with large, clear numbers for easy gameplay.

### Intermediate Level (16x16)

Balanced grid size with optimal visibility and control.

### Expert Level (30x16)

Wide grid that fits perfectly on screen with appropriate scaling.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Classic Minesweeper game design by Microsoft
- [Ebiten](https://ebiten.org/) game engine by Hajime Hoshi
- Inspired by traditional Windows Minesweeper

## Support

If you encounter any issues or have suggestions:

1. Check existing [Issues](https://github.com/yourusername/ebiten-tour/issues)
2. Create a new issue with detailed description
3. Include your Go version and operating system

---

**Enjoy playing Minesweeper! 💣🎮**
