package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	CellSize     = 30
	HeaderHeight = 60
)

type Difficulty int

const (
	Beginner Difficulty = iota
	Intermediate
	Expert
)

type DifficultyConfig struct {
	Width     int
	Height    int
	MineCount int
	Name      string
}

var difficultyConfigs = map[Difficulty]DifficultyConfig{
	Beginner:     {Width: 9, Height: 9, MineCount: 10, Name: "初級 (9x9, 10 mines)"},
	Intermediate: {Width: 16, Height: 16, MineCount: 40, Name: "中級 (16x16, 40 mines)"},
	Expert:       {Width: 30, Height: 16, MineCount: 99, Name: "上級 (30x16, 99 mines)"},
}

// Color mapping for mine count numbers
var numberColors = []color.Color{
	color.RGBA{0, 0, 0, 255},       // 0: black (unused)
	color.RGBA{0, 0, 255, 255},     // 1: blue
	color.RGBA{0, 128, 0, 255},     // 2: green
	color.RGBA{255, 0, 0, 255},     // 3: red
	color.RGBA{128, 0, 128, 255},   // 4: purple
	color.RGBA{128, 0, 0, 255},     // 5: maroon
	color.RGBA{0, 128, 128, 255},   // 6: cyan
	color.RGBA{0, 0, 0, 255},       // 7: black
	color.RGBA{128, 128, 128, 255}, // 8: gray
}

type CellState int

const (
	CellClosed CellState = iota
	CellOpen
	CellFlagged
)

type GameState int

const (
	GamePlaying GameState = iota
	GameWon
	GameLost
)

type Cell struct {
	HasMine       bool
	State         CellState
	NeighborMines int
}

type Game struct {
	board      [][]Cell
	gameState  GameState
	firstClick bool
	startTime  time.Time
	minesLeft  int
	difficulty Difficulty
	config     DifficultyConfig
}

func NewGame() *Game {
	return NewGameWithDifficulty(Intermediate) // Default to intermediate
}

func NewGameWithDifficulty(difficulty Difficulty) *Game {
	config := difficultyConfigs[difficulty]
	g := &Game{
		board:      make([][]Cell, config.Height),
		gameState:  GamePlaying,
		firstClick: true,
		minesLeft:  config.MineCount,
		difficulty: difficulty,
		config:     config,
	}

	for i := range g.board {
		g.board[i] = make([]Cell, config.Width)
	}

	return g
}

func (g *Game) placeMines(firstClickX, firstClickY int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	minesPlaced := 0

	for minesPlaced < g.config.MineCount {
		x := rng.Intn(g.config.Width)
		y := rng.Intn(g.config.Height)

		if !g.board[y][x].HasMine && !g.isFirstClickArea(x, y, firstClickX, firstClickY) {
			g.board[y][x].HasMine = true
			minesPlaced++
		}
	}

	g.calculateNeighborMines()
}

func (g *Game) isFirstClickArea(x, y, firstX, firstY int) bool {
	return abs(x-firstX) <= 1 && abs(y-firstY) <= 1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// drawLargeNumber draws a large colored number using simple pixel patterns
func (g *Game) drawLargeNumber(screen *ebiten.Image, num int, x, y int, c color.Color) {
	// Simple 5x7 bitmap patterns for numbers 1-8
	patterns := map[int][]string{
		1: {
			"  #  ",
			" ##  ",
			"  #  ",
			"  #  ",
			"  #  ",
			"  #  ",
			" ### ",
		},
		2: {
			" ### ",
			"#   #",
			"    #",
			"   # ",
			"  #  ",
			" #   ",
			"#####",
		},
		3: {
			" ### ",
			"#   #",
			"    #",
			" ### ",
			"    #",
			"#   #",
			" ### ",
		},
		4: {
			"   # ",
			"  ## ",
			" # # ",
			"#  # ",
			"#####",
			"   # ",
			"   # ",
		},
		5: {
			"#####",
			"#    ",
			"#    ",
			"#### ",
			"    #",
			"#   #",
			" ### ",
		},
		6: {
			" ### ",
			"#   #",
			"#    ",
			"#### ",
			"#   #",
			"#   #",
			" ### ",
		},
		7: {
			"#####",
			"    #",
			"   # ",
			"  #  ",
			" #   ",
			" #   ",
			" #   ",
		},
		8: {
			" ### ",
			"#   #",
			"#   #",
			" ### ",
			"#   #",
			"#   #",
			" ### ",
		},
	}

	pattern, exists := patterns[num]
	if !exists {
		return
	}

	// Scale factor for larger text
	scale := 2

	for row, line := range pattern {
		for col, char := range line {
			if char == '#' {
				// Draw a 2x2 pixel block for each '#'
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						px := x + col*scale + dx
						py := y + row*scale + dy
						if px >= 0 && py >= 0 && px < g.config.Width*CellSize && py < g.config.Height*CellSize+HeaderHeight {
							vector.DrawFilledRect(screen, float32(px), float32(py), 1, 1, c, false)
						}
					}
				}
			}
		}
	}
}

// drawMine draws a mine symbol that looks like a bomb
func (g *Game) drawMine(screen *ebiten.Image, x, y int) {
	// Draw a black circle for the bomb body
	centerX := x + 8
	centerY := y + 8
	radius := 6

	// Draw filled circle using small rectangles
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				px := centerX + dx
				py := centerY + dy
				if px >= 0 && py >= 0 && px < g.config.Width*CellSize && py < g.config.Height*CellSize+HeaderHeight {
					vector.DrawFilledRect(screen, float32(px), float32(py), 1, 1, color.RGBA{0, 0, 0, 255}, false)
				}
			}
		}
	}

	// Draw spikes around the bomb
	spikeColor := color.RGBA{0, 0, 0, 255}

	// Top spike
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(screen, float32(centerX), float32(centerY-radius-3+i), 1, 1, spikeColor, false)
	}
	// Bottom spike
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(screen, float32(centerX), float32(centerY+radius+1+i), 1, 1, spikeColor, false)
	}
	// Left spike
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(screen, float32(centerX-radius-3+i), float32(centerY), 1, 1, spikeColor, false)
	}
	// Right spike
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(screen, float32(centerX+radius+1+i), float32(centerY), 1, 1, spikeColor, false)
	}

	// Diagonal spikes
	for i := 0; i < 2; i++ {
		// Top-left
		vector.DrawFilledRect(screen, float32(centerX-4-i), float32(centerY-4-i), 1, 1, spikeColor, false)
		// Top-right
		vector.DrawFilledRect(screen, float32(centerX+4+i), float32(centerY-4-i), 1, 1, spikeColor, false)
		// Bottom-left
		vector.DrawFilledRect(screen, float32(centerX-4-i), float32(centerY+4+i), 1, 1, spikeColor, false)
		// Bottom-right
		vector.DrawFilledRect(screen, float32(centerX+4+i), float32(centerY+4+i), 1, 1, spikeColor, false)
	}
}

// drawLargeFlag draws a large flag symbol
func (g *Game) drawLargeFlag(screen *ebiten.Image, x, y int) {
	c := color.RGBA{255, 0, 0, 255} // Red flag
	pattern := []string{
		"##   ",
		"###  ",
		"##   ",
		"#    ",
		"#    ",
		"#    ",
		"#    ",
	}

	scale := 2
	for row, line := range pattern {
		for col, char := range line {
			if char == '#' {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						px := x + col*scale + dx
						py := y + row*scale + dy
						if px >= 0 && py >= 0 && px < g.config.Width*CellSize && py < g.config.Height*CellSize+HeaderHeight {
							vector.DrawFilledRect(screen, float32(px), float32(py), 1, 1, c, false)
						}
					}
				}
			}
		}
	}
}

func (g *Game) calculateNeighborMines() {
	for y := 0; y < g.config.Height; y++ {
		for x := 0; x < g.config.Width; x++ {
			if !g.board[y][x].HasMine {
				g.board[y][x].NeighborMines = g.countNeighborMines(x, y)
			}
		}
	}
}

func (g *Game) countNeighborMines(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if g.isValidPosition(nx, ny) && g.board[ny][nx].HasMine {
				count++
			}
		}
	}
	return count
}

func (g *Game) isValidPosition(x, y int) bool {
	return x >= 0 && x < g.config.Width && y >= 0 && y < g.config.Height
}

func (g *Game) openCell(x, y int) {
	if !g.isValidPosition(x, y) || g.board[y][x].State != CellClosed {
		return
	}

	if g.firstClick {
		g.placeMines(x, y)
		g.firstClick = false
		g.startTime = time.Now()
	}

	g.board[y][x].State = CellOpen

	if g.board[y][x].HasMine {
		g.gameState = GameLost
		g.revealAllMines()
		return
	}

	if g.board[y][x].NeighborMines == 0 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				g.openCell(x+dx, y+dy)
			}
		}
	}

	g.checkWinCondition()
}

func (g *Game) toggleFlag(x, y int) {
	if !g.isValidPosition(x, y) || g.board[y][x].State == CellOpen {
		return
	}

	if g.board[y][x].State == CellClosed {
		g.board[y][x].State = CellFlagged
		g.minesLeft--
	} else {
		g.board[y][x].State = CellClosed
		g.minesLeft++
	}
}

func (g *Game) revealAllMines() {
	for y := 0; y < g.config.Height; y++ {
		for x := 0; x < g.config.Width; x++ {
			if g.board[y][x].HasMine {
				g.board[y][x].State = CellOpen
			}
		}
	}
}

func (g *Game) checkWinCondition() {
	for y := 0; y < g.config.Height; y++ {
		for x := 0; x < g.config.Width; x++ {
			if !g.board[y][x].HasMine && g.board[y][x].State != CellOpen {
				return
			}
		}
	}
	g.gameState = GameWon
}

func (g *Game) Update() error {
	// Restart current difficulty
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *NewGameWithDifficulty(g.difficulty)
		return nil
	}

	// Difficulty switching
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		*g = *NewGameWithDifficulty(Beginner)
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		*g = *NewGameWithDifficulty(Intermediate)
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		*g = *NewGameWithDifficulty(Expert)
		return nil
	}

	if g.gameState != GamePlaying {
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= HeaderHeight {
			cellX := x / CellSize
			cellY := (y - HeaderHeight) / CellSize
			if cellX < g.config.Width && cellY < g.config.Height {
				g.openCell(cellX, cellY)
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		if y >= HeaderHeight {
			cellX := x / CellSize
			cellY := (y - HeaderHeight) / CellSize
			if cellX < g.config.Width && cellY < g.config.Height {
				g.toggleFlag(cellX, cellY)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{192, 192, 192, 255})

	g.drawHeader(screen)
	g.drawBoard(screen)
}

func (g *Game) drawHeader(screen *ebiten.Image) {
	windowWidth := g.config.Width * CellSize
	headerBg := color.RGBA{128, 128, 128, 255}
	vector.DrawFilledRect(screen, 0, 0, float32(windowWidth), HeaderHeight, headerBg, false)

	// Adapt layout based on window width
	isNarrow := windowWidth < 400 // Small window (beginner level)

	if isNarrow {
		// Compact layout for narrow windows
		minesText := fmt.Sprintf("M:%d", g.minesLeft)
		ebitenutil.DebugPrintAt(screen, minesText, 5, 15)

		var elapsed time.Duration
		if !g.firstClick && g.gameState == GamePlaying {
			elapsed = time.Since(g.startTime)
		}
		timeText := fmt.Sprintf("T:%d", int(elapsed.Seconds()))
		ebitenutil.DebugPrintAt(screen, timeText, 60, 15)

		// Show short difficulty name
		var shortDifficulty string
		switch g.difficulty {
		case Beginner:
			shortDifficulty = "Easy"
		case Intermediate:
			shortDifficulty = "Med"
		case Expert:
			shortDifficulty = "Hard"
		}
		ebitenutil.DebugPrintAt(screen, shortDifficulty, 115, 15)

		// Compact status text
		var statusText string
		switch g.gameState {
		case GameWon:
			statusText = "WIN! R:restart"
		case GameLost:
			statusText = "LOSE! R:restart"
		default:
			statusText = "1:E 2:M 3:H R:restart"
		}
		ebitenutil.DebugPrintAt(screen, statusText, 5, 45)
	} else {
		// Full layout for wider windows
		minesText := fmt.Sprintf("Mines: %d", g.minesLeft)
		ebitenutil.DebugPrintAt(screen, minesText, 10, 15)

		var elapsed time.Duration
		if !g.firstClick && g.gameState == GamePlaying {
			elapsed = time.Since(g.startTime)
		}
		timeText := fmt.Sprintf("Time: %d", int(elapsed.Seconds()))
		ebitenutil.DebugPrintAt(screen, timeText, 10, 30)

		// Show current difficulty
		difficultyText := g.config.Name
		ebitenutil.DebugPrintAt(screen, difficultyText, 10, 45)

		var statusText string
		switch g.gameState {
		case GameWon:
			statusText = "YOU WIN! Press R to restart"
		case GameLost:
			statusText = "GAME OVER! Press R to restart"
		default:
			statusText = "1:Beginner 2:Intermediate 3:Expert R:Restart"
		}
		maxX := windowWidth - 10
		if len(statusText)*7 < maxX {
			ebitenutil.DebugPrintAt(screen, statusText, maxX-len(statusText)*7, 15)
		} else {
			ebitenutil.DebugPrintAt(screen, statusText, maxX-300, 15)
		}
	}
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	for y := 0; y < g.config.Height; y++ {
		for x := 0; x < g.config.Width; x++ {
			g.drawCell(screen, x, y)
		}
	}
}

func (g *Game) drawCell(screen *ebiten.Image, x, y int) {
	screenX := float64(x * CellSize)
	screenY := float64(y*CellSize + HeaderHeight)
	cell := &g.board[y][x]

	var cellColor color.Color

	switch cell.State {
	case CellClosed:
		cellColor = color.RGBA{220, 220, 220, 255}
	case CellFlagged:
		cellColor = color.RGBA{255, 255, 0, 255}
	case CellOpen:
		if cell.HasMine {
			cellColor = color.RGBA{255, 0, 0, 255}
		} else {
			cellColor = color.RGBA{240, 240, 240, 255}
		}
	}

	vector.DrawFilledRect(screen, float32(screenX), float32(screenY), CellSize-1, CellSize-1, cellColor, false)

	if cell.State == CellOpen {
		if cell.HasMine {
			g.drawMine(screen, int(screenX)+7, int(screenY)+8)
		} else if cell.NeighborMines > 0 {
			// Get color for this number
			var textColor color.Color
			if cell.NeighborMines < len(numberColors) {
				textColor = numberColors[cell.NeighborMines]
			} else {
				textColor = color.RGBA{0, 0, 0, 255} // Default black
			}
			g.drawLargeNumber(screen, cell.NeighborMines, int(screenX)+7, int(screenY)+8, textColor)
		}
	} else if cell.State == CellFlagged {
		g.drawLargeFlag(screen, int(screenX)+7, int(screenY)+8)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.Width * CellSize, g.config.Height*CellSize + HeaderHeight
}

func main() {
	game := NewGame()
	windowWidth := game.config.Width * CellSize
	windowHeight := game.config.Height*CellSize + HeaderHeight
	ebiten.SetWindowSize(windowWidth*2, windowHeight*2)
	ebiten.SetWindowTitle("Minesweeper")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
