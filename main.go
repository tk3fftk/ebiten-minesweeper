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
	GridWidth    = 16
	GridHeight   = 16
	MineCount    = 40
	WindowWidth  = GridWidth * CellSize
	WindowHeight = GridHeight*CellSize + 60
)

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
}

func NewGame() *Game {
	g := &Game{
		board:      make([][]Cell, GridHeight),
		gameState:  GamePlaying,
		firstClick: true,
		minesLeft:  MineCount,
	}

	for i := range g.board {
		g.board[i] = make([]Cell, GridWidth)
	}

	return g
}

func (g *Game) placeMines(firstClickX, firstClickY int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	minesPlaced := 0

	for minesPlaced < MineCount {
		x := rng.Intn(GridWidth)
		y := rng.Intn(GridHeight)

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
func drawLargeNumber(screen *ebiten.Image, num int, x, y int, c color.Color) {
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
						if px >= 0 && py >= 0 && px < WindowWidth && py < WindowHeight+60 {
							vector.DrawFilledRect(screen, float32(px), float32(py), 1, 1, c, false)
						}
					}
				}
			}
		}
	}
}

// drawLargeMine draws a large mine symbol
func drawLargeMine(screen *ebiten.Image, x, y int) {
	c := color.RGBA{255, 0, 0, 255} // Red mine
	pattern := []string{
		"  #  ",
		" ### ",
		"#####",
		" ### ",
		"  #  ",
	}

	scale := 2
	for row, line := range pattern {
		for col, char := range line {
			if char == '#' {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						px := x + col*scale + dx
						py := y + row*scale + dy + 2 // Offset slightly
						if px >= 0 && py >= 0 && px < WindowWidth && py < WindowHeight+60 {
							vector.DrawFilledRect(screen, float32(px), float32(py), 1, 1, c, false)
						}
					}
				}
			}
		}
	}
}

// drawLargeFlag draws a large flag symbol
func drawLargeFlag(screen *ebiten.Image, x, y int) {
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
						if px >= 0 && py >= 0 && px < WindowWidth && py < WindowHeight+60 {
							vector.DrawFilledRect(screen, float32(px), float32(py), 1, 1, c, false)
						}
					}
				}
			}
		}
	}
}

func (g *Game) calculateNeighborMines() {
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
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
	return x >= 0 && x < GridWidth && y >= 0 && y < GridHeight
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
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			if g.board[y][x].HasMine {
				g.board[y][x].State = CellOpen
			}
		}
	}
}

func (g *Game) checkWinCondition() {
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			if !g.board[y][x].HasMine && g.board[y][x].State != CellOpen {
				return
			}
		}
	}
	g.gameState = GameWon
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *NewGame()
		return nil
	}

	if g.gameState != GamePlaying {
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 60 {
			cellX := x / CellSize
			cellY := (y - 60) / CellSize
			g.openCell(cellX, cellY)
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		if y >= 60 {
			cellX := x / CellSize
			cellY := (y - 60) / CellSize
			g.toggleFlag(cellX, cellY)
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
	headerBg := color.RGBA{128, 128, 128, 255}
	vector.DrawFilledRect(screen, 0, 0, WindowWidth, 60, headerBg, false)

	minesText := fmt.Sprintf("Mines: %d", g.minesLeft)
	ebitenutil.DebugPrintAt(screen, minesText, 10, 20)

	var elapsed time.Duration
	if !g.firstClick && g.gameState == GamePlaying {
		elapsed = time.Since(g.startTime)
	}
	timeText := fmt.Sprintf("Time: %d", int(elapsed.Seconds()))
	ebitenutil.DebugPrintAt(screen, timeText, 10, 40)

	var statusText string
	switch g.gameState {
	case GameWon:
		statusText = "YOU WIN! Press R to restart"
	case GameLost:
		statusText = "GAME OVER! Press R to restart"
	default:
		statusText = "Left click: open, Right click: flag"
	}
	ebitenutil.DebugPrintAt(screen, statusText, WindowWidth-300, 20)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			g.drawCell(screen, x, y)
		}
	}
}

func (g *Game) drawCell(screen *ebiten.Image, x, y int) {
	screenX := float64(x * CellSize)
	screenY := float64(y*CellSize + 60)
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
			drawLargeMine(screen, int(screenX)+7, int(screenY)+8)
		} else if cell.NeighborMines > 0 {
			// Get color for this number
			var textColor color.Color
			if cell.NeighborMines < len(numberColors) {
				textColor = numberColors[cell.NeighborMines]
			} else {
				textColor = color.RGBA{0, 0, 0, 255} // Default black
			}
			drawLargeNumber(screen, cell.NeighborMines, int(screenX)+7, int(screenY)+8, textColor)
		}
	} else if cell.State == CellFlagged {
		drawLargeFlag(screen, int(screenX)+7, int(screenY)+8)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(WindowWidth*2, WindowHeight*2)
	ebiten.SetWindowTitle("Minesweeper")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
