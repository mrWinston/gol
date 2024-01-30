package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var colorBG color.Color = color.RGBA{
	R: 128,
	G: 128,
	B: 128,
	A: 255,
}

var colorAlive color.Color = color.Black
var colorDead color.Color = color.White

var borderSize = 5

var updateInterval = 250 * time.Millisecond
var sinceUpdate = 2

var deadPixelImg *ebiten.Image
var alivePixelImg *ebiten.Image
var tileImageSize int = 100

func NewGameOfLife(width int, height int) *GameOfLife {
	gol := &GameOfLife{
		cells:      [][]bool{},
		lastUpdate: time.Now(),
	}
	gol.cells = make([][]bool, width)
	for i := range gol.cells {
		gol.cells[i] = make([]bool, height)
		for j := range gol.cells[i] {
			gol.cells[i][j] = rand.Intn(3) == 1
		}
	}
  
  deadPixelImg = ebiten.NewImage(tileImageSize,tileImageSize)
  deadPixelImg.Fill(colorDead)
  alivePixelImg = ebiten.NewImage(tileImageSize,tileImageSize)
  alivePixelImg.Fill(colorAlive)
  for i := 0; i < tileImageSize; i++ {
    for j := 0; j < borderSize; j++ {
      deadPixelImg.Set(i, j, colorBG)
      deadPixelImg.Set(i, tileImageSize - j, colorBG)
      deadPixelImg.Set(j, i, colorBG)
      deadPixelImg.Set(tileImageSize - j, i, colorBG)
      alivePixelImg.Set(i, j, colorBG)
      alivePixelImg.Set(i, tileImageSize - j, colorBG)
      alivePixelImg.Set(j, i, colorBG)
      alivePixelImg.Set(tileImageSize - j, i, colorBG)
    }
  }

	return gol
}

type GameOfLife struct {
	cells      [][]bool
	lastUpdate time.Time
  sinceUpdate int
}

// Update updates a game by one tick. The given argument represents a screen image.
// Update updates only the game logic and Draw draws the screen.
//
// You can assume that Update is always called TPS-times per second (60 by default), and you can assume
// that the time delta between two Updates is always 1 / TPS [s] (1/60[s] by default). As Ebitengine already
// adjusts the number of Update calls, you don't have to measure time deltas in Update by e.g. OS timers.
//
// An actual TPS is available by ActualTPS(), and the result might slightly differ from your expected TPS,
// but still, your game logic should stick to the fixed time delta and should not rely on ActualTPS() value.
// This API is for just measurement and/or debugging. In the long run, the number of Update calls should be
// adjusted based on the set TPS on average.
//
// An actual time delta between two Updates might be bigger than expected. In this case, your game's
// Update or Draw takes longer than they should. In this case, there is nothing other than optimizing
// your game implementation.
//
// In the first frame, it is ensured that Update is called at least once before Draw. You can use Update
// to initialize the game state.
//
// After the first frame, Update might not be called or might be called once
// or more for one frame. The frequency is determined by the current TPS (tick-per-second).
//
// If the error returned is nil, game execution proceeds normally.
// If the error returned is Termination, game execution halts, but does not return an error from RunGame.
// If the error returned is any other non-nil value, game execution halts and the error is returned from RunGame.
func (ga *GameOfLife) Update() error {
  if time.Since(ga.lastUpdate) < updateInterval {
    return nil
  }

	ga.lastUpdate = time.Now()
	cellsCopy := make([][]bool, len(ga.cells))
  fmt.Printf("Update called, %f\n", ebiten.ActualTPS())
  
	for x := range ga.cells {
    cellsCopy[x] = make([]bool, len(ga.cells[x]))
		for y := range ga.cells[x] {
			neighbors := ga.AliveNeighbors(x, y)
      cellsCopy[x][y] = ga.cells[x][y]
			if !ga.cells[x][y] {
				cellsCopy[x][y] = neighbors == 3
			}
      if ga.cells[x][y] {
        cellsCopy[x][y] = neighbors >= 2 && neighbors <= 3
      }

		}
	}
  ga.cells = cellsCopy

	return nil
}

func (ga *GameOfLife) AliveNeighbors(x int, y int) int {
	count := 0
  minX := max(x-1, 0)
  maxX := min(x+1, len(ga.cells) - 1)
  minY := max(y-1, 0)
  maxY := min(y+1, len(ga.cells[0]) - 1)
 
	for i := minX; i <= maxX; i++ {
		for j := minY; j <= maxY; j++ {
			if j-minY == 1 && i-minX == 1 {
				continue
			}
      if ga.cells[i][j] {	
        count++ 
      }
		}
	}
	return count
}

// Draw draws the game screen by one frame.
//
// The give argument represents a screen image. The updated content is adopted as the game screen.
//
// The frequency of Draw calls depends on the user's environment, especially the monitors refresh rate.
// For portability, you should not put your game logic in Draw in general.
func (ga *GameOfLife) Draw(screen *ebiten.Image) {
	cellsX := len(ga.cells)
	cellsY := len(ga.cells[0])

	screenX := screen.Bounds().Dx()
	screenY := screen.Bounds().Dy()
	pixPerCellX := screenX / cellsX
	pixPerCellY := screenY / cellsY

	pixPerCell := min(pixPerCellX, pixPerCellY)

	screen.Fill(colorBG)
	for x := range ga.cells {
		for y := range ga.cells[x] {
      dio := ebiten.DrawImageOptions{}
      scaleFactor := float64(pixPerCell)/float64(tileImageSize)
      dio.GeoM.Scale(scaleFactor, scaleFactor)
      dio.GeoM.Translate(float64(x * pixPerCell), float64(y * pixPerCell))
      if ga.cells[x][y] {
        screen.DrawImage(alivePixelImg, &dio)
      } else {
        screen.DrawImage(deadPixelImg, &dio)
      }
		}
	}
}


// Layout accepts a native outside size in device-independent pixels and returns the game's logical screen
// size.
//
// On desktops, the outside is a window or a monitor (fullscreen mode). On browsers, the outside is a body
// element. On mobiles, the outside is the view's size.
//
// Even though the outside size and the screen size differ, the rendering scale is automatically adjusted to
// fit with the outside.
//
// Layout is called almost every frame.
//
// It is ensured that Layout is invoked before Update is called in the first frame.
//
// If Layout returns non-positive numbers, the caller can panic.
//
// You can return a fixed screen size if you don't care, or you can also return a calculated screen size
// adjusted with the given outside size.
//
// If the game implements the interface LayoutFer, Layout is never called and LayoutF is called instead.
func (ga *GameOfLife) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}
