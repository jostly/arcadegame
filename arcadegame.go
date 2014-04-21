package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
	"github.com/jackyb/go-sdl2/sdl_mixer"
	"github.com/jackyb/go-sdl2/sdl_ttf"
	"github.com/jostly/arcadegame/game"
)

const (
	SCREEN_WIDTH            = 640
	SCREEN_HEIGHT           = 480
	EnergyCostPerShot       = 15.
	EnergyRegainedPerSecond = 15.
)

var (
	x      = float64(SCREEN_WIDTH / 4)
	y      = float64(SCREEN_HEIGHT / 2)
	score  = 0
	energy = 100.
)

func error(message string) {
	log.Printf(message+": %v\n", sdl.GetError())
}

func main() {

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		error("Error initializing SDL")
		return
	}

	defer sdl.Quit()

	ttf.Init()

	window := sdl.CreateWindow("Hello World!", 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT,
		sdl.WINDOW_SHOWN)

	if window == nil {
		error("Error opening window")
		return
	}

	defer window.Destroy()

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	if renderer == nil {
		error("Error creating renderer")
		return
	}

	defer renderer.Destroy()

	texture := img.LoadTexture(renderer, "assets/images/exampletexture_2.png")

	if texture == nil {
		error("Error creating texture")
		return
	}

	defer texture.Destroy()

	font, err := ttf.OpenFont("assets/fonts/ComicNeue-Regular.ttf", 18)

	if err != nil {
		log.Printf("Error when loading font: %v", error)
		return
	}

	defer font.Close()

	ok := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 1024)
	if !ok {
		error("Can't open audio")
		return
	}

	defer mix.CloseAudio()

	mix.AllocateChannels(2)

	// Might have to convert sound effects to wav to load them as audio
	// and not music?
	shoot := mix.LoadMUS("assets/audio/shoot.ogg")
	if shoot == nil {
		error("Can't load ogg sound")
		return
	}

	defer shoot.Free()

	explode := mix.LoadMUS("assets/audio/invaderkilled.ogg")
	if explode == nil {
		error("Can't load ogg sound")
		return
	}

	defer explode.Free()

	game.RenderCallback = func(r *sdl.Renderer) {

		drawShip(r)
		drawMissiles(r)
		drawObstacles(r)
	}

	lastObstacleTick := sdl.GetTicks()

	game.UpdateCallback = func(delta float64) {

		energy += EnergyRegainedPerSecond * delta
		if energy > 100. {
			energy = 100.
		}

		moveSpeed := delta * 100

		keystate := sdl.GetKeyboardState()

		updateMissiles(moveSpeed * 3)
		updateObstacles(delta, explode)

		if sdl.GetTicks() > lastObstacleTick+500 && rand.Intn(20) == 0 {
			lastObstacleTick = sdl.GetTicks()
			size := rand.Float64()*20.0 + 10.0
			y := rand.Float64() * SCREEN_HEIGHT
			x := SCREEN_WIDTH + size
			speed := rand.Float64()*80.0 + 50.0
			obstacles = append(obstacles, Obstacle{x, y, speed, size, 0})
		}

		if keystate[sdl.SCANCODE_W] != 0 {
			y -= moveSpeed
		}
		if keystate[sdl.SCANCODE_S] != 0 {
			y += moveSpeed
		}
		if keystate[sdl.SCANCODE_A] != 0 {
			x -= moveSpeed
		}
		if keystate[sdl.SCANCODE_D] != 0 {
			x += moveSpeed
		}
		if keystate[sdl.SCANCODE_RSHIFT] != 0 {
			currentTick := sdl.GetTicks()
			if currentTick > lastFire+300 && energy > EnergyCostPerShot {
				lastFire = currentTick
				missiles = append(missiles, FloatPoint{x + 32, y})
				energy -= EnergyCostPerShot
				shoot.Play(1)
			}
		}
	}

	game.StatusCallback = func() string {
		return fmt.Sprintf("Energy: %3d   Score: %d", int(energy), score)
	}

	game.MainLoop(renderer, font)

}

type FloatPoint struct {
	X, Y float64
}

var lastFire = sdl.GetTicks()

var shipPoints = [...]sdl.Point{sdl.Point{30, 0}, sdl.Point{-20, -15}, sdl.Point{-20, 15}}

func drawShip(r *sdl.Renderer) {
	ship := make([]sdl.Point, 4)
	for i := 0; i < 4; i++ {
		p := shipPoints[i%3]
		ship[i] = sdl.Point{p.X + int32(x), p.Y + int32(y)}
	}

	r.SetDrawColor(255, 255, 255, 255)
	err := r.DrawLines(ship)
	if err != 0 {
		error("Draw error")
	}
}

var missiles = []FloatPoint{}

func drawMissiles(r *sdl.Renderer) {
	for _, p := range missiles {
		r.DrawPoint(int(p.X), int(p.Y))
	}
}

func updateMissiles(moveSpeed float64) {
	newMissiles := make([]FloatPoint, 0, len(missiles))
	for _, p := range missiles {
		p.X += moveSpeed

		if p.X <= SCREEN_WIDTH {
			newMissiles = append(newMissiles, p)
		}
	}
	missiles = newMissiles
}

type Obstacle struct {
	X, Y, Speed, Size, Angle float64
}

var obstacles = []Obstacle{}

func drawObstacles(r *sdl.Renderer) {
	r.SetDrawColor(255, 128, 100, 255)
	for _, o := range obstacles {
		rect := sdl.Rect{int32(o.X - o.Size), int32(o.Y - o.Size), int32(o.Size * 2), int32(o.Size * 2)}
		r.DrawRect(&rect)
	}
}

func updateObstacles(delta float64, explode *mix.Music) {
	newObstacles := make([]Obstacle, 0, len(obstacles))
	for _, o := range obstacles {
		o.X -= delta * o.Speed
		o.Angle += delta * o.Speed
		xmin := o.X - o.Size
		xmax := o.X + o.Size
		ymin := o.Y - o.Size
		ymax := o.Y + o.Size
		for i, m := range missiles {

			if m.X >= xmin && m.X <= xmax && m.Y >= ymin && m.Y <= ymax {
				explode.Play(1)
				o.X = -o.Size
				missiles[i].X = SCREEN_WIDTH * 10
				score += int(40 - o.Size)
				break
			}

		}

		if o.X > -o.Size {
			newObstacles = append(newObstacles, o)
		}
	}
	obstacles = newObstacles
}
