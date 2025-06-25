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
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	CellSize = 30
	GridWidth = 16
	GridHeight = 16
	MineCount = 40
	WindowWidth = GridWidth * CellSize
	WindowHeight = GridHeight*CellSize + 60
)

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
	HasMine      bool
	State        CellState
	NeighborMines int
}

type Game struct {
	board       [][]Cell
	gameState   GameState
	firstClick  bool
	startTime   time.Time
	minesLeft   int
}

func NewGame() *Game {
	g := &Game{
		board:     make([][]Cell, GridHeight),
		gameState: GamePlaying,
		firstClick: true,
		minesLeft: MineCount,
	}
	
	for i := range g.board {
		g.board[i] = make([]Cell, GridWidth)
	}
	
	return g
}

func (g *Game) placeMines(firstClickX, firstClickY int) {
	rand.Seed(time.Now().UnixNano())
	minesPlaced := 0
	
	for minesPlaced < MineCount {
		x := rand.Intn(GridWidth)
		y := rand.Intn(GridHeight)
		
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
	ebitenutil.DrawRect(screen, 0, 0, WindowWidth, 60, headerBg)
	
	minesText := fmt.Sprintf("Mines: %d", g.minesLeft)
	text.Draw(screen, minesText, basicfont.Face7x13, 10, 20, color.White)
	
	var elapsed time.Duration
	if !g.firstClick && g.gameState == GamePlaying {
		elapsed = time.Since(g.startTime)
	}
	timeText := fmt.Sprintf("Time: %d", int(elapsed.Seconds()))
	text.Draw(screen, timeText, basicfont.Face7x13, 10, 40, color.White)
	
	var statusText string
	var statusColor color.Color = color.White
	switch g.gameState {
	case GameWon:
		statusText = "YOU WIN! Press R to restart"
		statusColor = color.RGBA{0, 255, 0, 255}
	case GameLost:
		statusText = "GAME OVER! Press R to restart"
		statusColor = color.RGBA{255, 0, 0, 255}
	default:
		statusText = "Left click: open, Right click: flag"
	}
	text.Draw(screen, statusText, basicfont.Face7x13, WindowWidth-300, 20, statusColor)
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
	var textColor color.Color = color.Black
	
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
	
	ebitenutil.DrawRect(screen, screenX, screenY, CellSize-1, CellSize-1, cellColor)
	
	if cell.State == CellOpen {
		if cell.HasMine {
			text.Draw(screen, "*", basicfont.Face7x13, int(screenX)+10, int(screenY)+20, color.Black)
		} else if cell.NeighborMines > 0 {
			numColors := []color.Color{
				color.Black,
				color.RGBA{0, 0, 255, 255},   // 1: blue
				color.RGBA{0, 128, 0, 255},   // 2: green  
				color.RGBA{255, 0, 0, 255},   // 3: red
				color.RGBA{128, 0, 128, 255}, // 4: purple
				color.RGBA{128, 0, 0, 255},   // 5: maroon
				color.RGBA{0, 128, 128, 255}, // 6: teal
				color.RGBA{0, 0, 0, 255},     // 7: black
				color.RGBA{128, 128, 128, 255}, // 8: gray
			}
			if cell.NeighborMines < len(numColors) {
				textColor = numColors[cell.NeighborMines]
			}
			text.Draw(screen, fmt.Sprintf("%d", cell.NeighborMines), basicfont.Face7x13, int(screenX)+10, int(screenY)+20, textColor)
		}
	} else if cell.State == CellFlagged {
		text.Draw(screen, "F", basicfont.Face7x13, int(screenX)+10, int(screenY)+20, color.RGBA{255, 0, 0, 255})
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