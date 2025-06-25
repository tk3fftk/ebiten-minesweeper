# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Minesweeper game implemented in Go using the Ebiten 2D game engine. The entire game is contained in a single `main.go` file with a monolithic architecture suitable for a simple game project.

## Development Commands

### Building and Running

```bash
# Build the game
go build -o minesweeper

# Run directly with Go
go run main.go

# Update dependencies
go mod tidy

# format
gofmt -w .
```

### Game Controls

- Left click: Open cell
- Right click: Flag/unflag cell  
- R key: Restart current difficulty
- 1 key: Switch to Beginner (9x9, 10 mines)
- 2 key: Switch to Intermediate (16x16, 40 mines)  
- 3 key: Switch to Expert (30x16, 99 mines)

## Code Architecture

### Core Data Structures

- `Cell`: Represents individual grid cells with mine status, visibility state, and neighbor count
- `Game`: Main game state containing the 2D board array, game status, timing, and mine counter
- `CellState`: Enum for cell visibility (Closed, Open, Flagged)
- `GameState`: Enum for overall game status (Playing, Won, Lost)

### Key Game Logic Flow

1. **Initialization**: `NewGame()` creates empty board, mines placed on first click
2. **Mine Placement**: `placeMines()` ensures first click area (3x3) is mine-free
3. **Cell Opening**: `openCell()` handles recursive opening of empty areas
4. **Win/Loss Detection**: Automatic checking after each move

### Rendering Architecture

The game uses Ebiten's immediate mode rendering:

- `Update()`: Handles input and game logic (60 FPS)
- `Draw()`: Renders header and board each frame
- `Layout()`: Defines fixed window dimensions

### Configuration Constants

Game supports three difficulty levels with dynamic sizing:

- **Beginner**: 9x9 grid with 10 mines
- **Intermediate**: 16x16 grid with 40 mines (default)
- **Expert**: 30x16 grid with 99 mines
- Cell size: 30 pixels (constant)
- Window dimensions calculated dynamically based on selected difficulty
- Header height: 60 pixels for UI controls

## Ebiten-Specific Considerations

- Uses deprecated `ebiten/v2/text` and `ebitenutil.DrawRect` APIs
- Mouse input handled through `inpututil` package for precise click detection
- Color rendering uses `color.RGBA` values throughout
- Font rendering uses `basicfont.Face7x13` for simplicity

## Architecture Notes

The single-file design prioritizes simplicity over modularity. For future expansion (multiple difficulty levels, high scores, etc.), consider refactoring into separate packages for game logic, rendering, and input handling.
