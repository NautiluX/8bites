package main

import (
	_ "image/png"
	"log"
	"math/rand/v2"
	"time"

	"github.com/NautiluX/8bites/assets"
	"github.com/NautiluX/8bites/pkg/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	mapWidth     = 640 / 32
	mapHeight    = 480 / 32
	playerSpeed  = 2
)

type Game struct {
	bgImage      *ebiten.Image
	playerSprite *sprites.CharacterSprite
	slimeSprite  *sprites.CharacterSprite
	mapTiles     [mapHeight][mapWidth]int
	wallTile     *ebiten.Image
	floorTile    *ebiten.Image
}

var theGame *Game

// init loads the assets before the game starts.
func init() {
	playerImg, err := assets.GetPlayerYellowSprite()
	if err != nil {
		log.Fatalf("Failed to load player sprite: %v", err)
	}
	playerSprite := sprites.NewCharacterSprite(playerImg, 32, 32, []sprites.Animation{
		{Name: "right", Frames: 12},
		{Name: "left", Frames: 12},
		{Name: "up", Frames: 12},
		{Name: "down", Frames: 12},
	})

	slimeImg, err := assets.GetSlimeSprite()
	if err != nil {
		log.Fatalf("Failed to load player sprite: %v", err)
	}
	slimeSprite := sprites.NewCharacterSprite(slimeImg, 32, 32, []sprites.Animation{
		{Name: "idle", Frames: 10},
	})

	wallTile, err := assets.GetWallTileImage()
	if err != nil {
		log.Fatalf("Failed to load wall tile image: %v", err)
	}

	floorTile, err := assets.GetFloorTileImage()
	if err != nil {
		log.Fatalf("Failed to load floor tile image: %v", err)
	}
	// Initialize the game struct (Global access for simplicity)
	theGame = &Game{
		bgImage:      nil,
		playerSprite: playerSprite,
		slimeSprite:  slimeSprite,
		mapTiles: [mapHeight][mapWidth]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 1, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1},
			{1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 1},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		wallTile:  wallTile,
		floorTile: floorTile,
	}

	playerSprite.X, playerSprite.Y = theGame.GetRandomFloorPosition()
	//select random tile to spawn slime
	slimeSprite.X, slimeSprite.Y = theGame.GetRandomFloorPosition()

}

func (g *Game) GetRandomFloorPosition() (int, int) {
	for {
		x := rand.IntN(mapWidth)
		y := rand.IntN(mapHeight)
		if g.mapTiles[y][x] == 0 {
			return x * 32, y * 32
		}
	}
}

// handleInputAndMovement processes keyboard input and updates the player's position,
// enforcing screen boundaries.
func (g *Game) handleInputAndMovement() {
	// Handle Up/W
	if g.playerSprite.Y%32 == 0 && g.playerSprite.X%32 == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			g.playerSprite.CurrentVy = -playerSpeed
			g.playerSprite.CurrentVx = 0
			g.playerSprite.SetAnimation("up")
		}
		// Handle Down/S
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			g.playerSprite.CurrentVy = playerSpeed
			g.playerSprite.CurrentVx = 0
			g.playerSprite.SetAnimation("down")
		}
		// Handle Left/A
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			g.playerSprite.CurrentVx = -playerSpeed
			g.playerSprite.CurrentVy = 0
			g.playerSprite.SetAnimation("left")
		}
		// Handle Right/D
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			g.playerSprite.CurrentVx = playerSpeed
			g.playerSprite.CurrentVy = 0
			g.playerSprite.SetAnimation("right")
		}
	}

	// Random movement for slime. 50% chance to change direction each update
	updateMovement := rand.IntN(101)
	if updateMovement > 75 && g.slimeSprite.CurrentVx == 0 && g.slimeSprite.X%32 == 0 && g.slimeSprite.Y%32 == 0 {
		g.slimeSprite.CurrentVx = -1 + rand.IntN(3)
		g.slimeSprite.CurrentVy = 0
	}
	if updateMovement < 25 && g.slimeSprite.CurrentVy == 0 && g.slimeSprite.Y%32 == 0 && g.slimeSprite.X%32 == 0 {
		g.slimeSprite.CurrentVx = 0
		g.slimeSprite.CurrentVy = -1 + rand.IntN(3)
	}

	if !g.checkWallCollision(g.slimeSprite) {
		g.slimeSprite.Move(screenWidth, screenHeight)
	}
	if !g.checkWallCollision(g.playerSprite) {
		g.playerSprite.Move(screenWidth, screenHeight)
	}
}

func (g *Game) checkWallCollision(s *sprites.CharacterSprite) bool {
	// figure out if the current movement would result in a collision with a wal.
	// we take the sprites current position width and height into account.
	newX := s.X + s.CurrentVx
	newY := s.Y + s.CurrentVy

	// Calculate the tile coordinates
	tileX1 := newX / 32
	tileY1 := newY / 32
	tileX2 := (newX + s.Width - 1) / 32
	tileY2 := (newY + s.Height - 1) / 32

	// Check all tiles the sprite would occupy
	for y := tileY1; y <= tileY2; y++ {
		for x := tileX1; x <= tileX2; x++ {
			if x < 0 || x >= mapWidth || y < 0 || y >= mapHeight {
				return true // Out of bounds is treated as a wall
			}
			if g.mapTiles[y][x] == 1 {
				return true // Collision with wall
			}
		}
	}
	return false
}

func (g *Game) animate() {
	g.playerSprite.Animate()
	g.slimeSprite.Animate()
}

var lastAnimationUpdate time.Time

// Update handles the game logic, primarily input and state changes.
func (g *Game) Update() error {
	// Call the separate function for handling movement
	g.handleInputAndMovement()

	if time.Since(lastAnimationUpdate) > 100*time.Millisecond {
		lastAnimationUpdate = time.Now()
		g.animate()
	}

	// Example: Exit on pressing Escape
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) drawMap(screen *ebiten.Image) {
	for y := range mapHeight {
		for x := range mapWidth {
			tile := g.mapTiles[y][x]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*32), float64(y*32))
			tileImg := g.floorTile
			if tile == 1 {
				tileImg = g.wallTile
			}
			screen.DrawImage(tileImg, op)
		}
	}
}

// Draw renders the game state to the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawMap(screen)

	// --- Draw Player ---
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(float64(g.playerSprite.X), float64(g.playerSprite.Y))
	playerImg := g.playerSprite.GetCurrentImage()
	screen.DrawImage(playerImg, playerOp)

	// --- Draw Slime ---
	slimeOp := &ebiten.DrawImageOptions{}
	slimeOp.GeoM.Translate(float64(g.slimeSprite.X), float64(g.slimeSprite.Y))
	slimeImg := g.slimeSprite.GetCurrentImage()
	screen.DrawImage(slimeImg, slimeOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	if theGame == nil {
		log.Fatal("Game initialization failed. Check the init function.")
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("8bites")
	if err := ebiten.RunGame(theGame); err != nil {
		log.Fatal(err)
	}
}
