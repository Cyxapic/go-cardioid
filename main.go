package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth          = 900
	screenHeight         = 900
	numLines             = 200
	translateW   float64 = screenWidth / 2
	translateH   float64 = screenHeight / 2
	dt                   = 16
)

func init() {
	rand.Seed(86)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cardioid")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

type Colors struct {
	colors []color.Color
	idx    int
}

func NewColors() *Colors {
	return &Colors{
		colors: []color.Color{
			color.RGBA{0, 255, 0, 150},
			color.RGBA{30, 250, 10, 150},
			color.RGBA{60, 240, 35, 150},
			color.RGBA{90, 230, 50, 150},
			color.RGBA{120, 210, 70, 150},
			color.RGBA{150, 100, 55, 150},
			color.RGBA{180, 80, 65, 150},
			color.RGBA{210, 30, 75, 150},
		},
	}
}

func (clr *Colors) NextColor() color.Color {
	if clr.idx >= len(clr.colors) {
		clr.idx = 0
	}
	c := clr.colors[clr.idx]
	clr.idx++
	return c
}

type Cardioid struct {
	tick          float64
	cardioidLines []*Line
	colors        *Colors
}

func NewCardiods() *Cardioid {
	lines := make([]*Line, numLines)
	for idxC := range lines {
		lines[idxC] = NewLine()
	}

	return &Cardioid{
		tick:          0,
		cardioidLines: lines,
		colors:        NewColors(),
	}
}

func (c *Cardioid) setNextCardiod() {
	c.tick = c.tick + dt
	factor := 1 + 0.0001*c.tick
	radius := 350 + 50*math.Abs(math.Sin(c.tick*0.004)-0.5)
	for idxC := range c.cardioidLines {
		c.cardioidLines[idxC].calcCoords(factor, radius, idxC)
		c.cardioidLines[idxC].setColor(c.colors.NextColor())
	}
}

func (cardioid *Cardioid) draw(screen *ebiten.Image) {
	for _, line := range cardioid.cardioidLines {
		line.draw(screen)
	}
}

type Line struct {
	x1  float32
	y1  float32
	x2  float32
	y2  float32
	clr color.Color
}

func NewLine() *Line {
	return &Line{}
}

func (line *Line) setColor(clr color.Color) {
	line.clr = clr
}

func (line *Line) calcCoords(factor, radiusFloat float64, currentLine int) {
	theta := (2.0 * math.Pi / float64(numLines)) * float64(currentLine)
	line.x1 = float32(radiusFloat*math.Cos(theta) + translateW)
	line.y1 = float32(radiusFloat*math.Sin(theta) + translateH)
	line.x2 = float32(radiusFloat*math.Cos(factor*theta) + translateW)
	line.y2 = float32(radiusFloat*math.Sin(factor*theta) + translateH)
}

func (line *Line) draw(screen *ebiten.Image) {
	vector.StrokeLine(screen, line.x1, line.y1, line.x2, line.y2, 1, line.clr, true)
}

type Game struct {
	cardioid *Cardioid
	tick     time.Time
}

func NewGame() *Game {
	cardiods := NewCardiods()
	g := &Game{cardioid: cardiods}
	g.tick = time.Now()
	g.cardioid.setNextCardiod()
	return g
}

func (g *Game) Update() error {
	if time.Since(g.tick) > time.Millisecond*50 {
		g.cardioid.setNextCardiod()
		g.tick = time.Now()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cardioid.draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
