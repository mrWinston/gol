package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var width int = 88
var height int = 66

func main() {
  gol := NewGameOfLife(width, height)

  ebiten.SetWindowSize(640*2, 480*2)
  ebiten.SetWindowTitle("Game of Life")
  ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
  ebiten.SetFullscreen(false)
  err := ebiten.RunGame(gol)
  if err != nil {
    log.Panic(err)
  }
}
