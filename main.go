package main

import (
	"bytes"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/NautiluX/8bites/assets"
	"github.com/NautiluX/8bites/pkg/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	playerSprite *sprites.Player
	enemies      []sprites.CharacterSprite
	bites        []*sprites.CharacterSprite
	eatenBites   []sprites.CharacterSprite
	currentBite  *sprites.CharacterSprite
	mapTiles     [mapHeight][mapWidth]int
	wallTile     *ebiten.Image
	floorTile    *ebiten.Image
	title        GameTitle
	Ended        bool
}

type GameTitle struct {
	Visible       bool
	Duration      time.Duration
	StartTime     time.Time
	LastShakeTime time.Time
	WordsVisible  int
	ShakeX        int
	ShakeY        int
	Text          string
}

var (
	theGame *Game
	font    *text.GoTextFaceSource
	bites   []*sprites.CharacterSprite
	enemies []*sprites.CharacterSprite
)

// init loads the assets before the game starts.
func init() {
	var err error
	font, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}

	bites = []*sprites.CharacterSprite{
		loadSprite("cheese", sprites.SpriteIdCheese, 9),
		loadSprite("pizza", sprites.SpriteIdPizza, 10),
		loadSprite("donut", sprites.SpriteIdDonut, 23),
		loadSprite("sushi", sprites.SpriteIdSushi, 12),
		loadSprite("orange", sprites.SpriteIdOrange, 8),
		loadSprite("avocado", sprites.SpriteIdAvocado, 21),
		loadSprite("apple", sprites.SpriteIdApple, 20),
		loadSprite("banana", sprites.SpriteIdBanana, 21),
	}

	slimeImg, err := assets.GetSlimeSprite()
	if err != nil {
		log.Fatalf("failed to load player sprite: %v", err)
	}
	slimeSprite := sprites.NewCharacterSprite(slimeImg, 32, 32, []sprites.Animation{
		{Name: "idle", Frames: 10},
	}, sprites.SpriteIdSlime)

	enemies = []*sprites.CharacterSprite{slimeSprite}
	theGame = &Game{}
	err = ResetGame()
	if err != nil {
		log.Fatal(err)
	}
}

func loadSprite(spriteName string, spriteId sprites.SpriteId, frames int) *sprites.CharacterSprite {
	img, err := assets.GetItem(spriteName)
	if err != nil {
		log.Fatalf("failed to load donut image: %v", err)
	}
	sprite := sprites.NewCharacterSprite(img, 32, 32, []sprites.Animation{
		{Name: "idle", Frames: frames},
	}, spriteId)
	return sprite
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
	// Move queued
	if g.playerSprite.Y%32 == 0 && g.playerSprite.X%32 == 0 && (g.playerSprite.NextVx != 0 || g.playerSprite.NextVy != 0) {
		g.playerSprite.CurrentVx = g.playerSprite.NextVx
		g.playerSprite.CurrentVy = g.playerSprite.NextVy
		g.playerSprite.CurrentAnimation = g.playerSprite.NextAnimation
		g.playerSprite.NextVx = 0
		g.playerSprite.NextVy = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.playerSprite.CurrentVy != 0 {
			g.playerSprite.CurrentVy = -playerSpeed
			g.playerSprite.SetAnimation("up")
		} else {
			g.playerSprite.NextVy = -playerSpeed
			g.playerSprite.SetNextAnimation("up")
		}
	}
	// Handle Down/S
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.playerSprite.CurrentVy != 0 {
			g.playerSprite.CurrentVy = playerSpeed
			g.playerSprite.SetAnimation("down")
		} else {
			g.playerSprite.NextVy = playerSpeed
			g.playerSprite.SetNextAnimation("down")
		}
	}
	// Handle Left/A
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if g.playerSprite.CurrentVx != 0 {
			g.playerSprite.CurrentVx = -playerSpeed
			g.playerSprite.SetAnimation("left")
		} else {
			g.playerSprite.NextVx = -playerSpeed
			g.playerSprite.SetNextAnimation("left")
		}
	}
	// Handle Right/D
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if g.playerSprite.CurrentVx != 0 {
			g.playerSprite.CurrentVx = playerSpeed
			g.playerSprite.SetAnimation("right")
		} else {
			g.playerSprite.NextVx = playerSpeed
			g.playerSprite.SetNextAnimation("right")
		}
	}

	for i := range g.enemies {
		// Random movement for slime. 50% chance to change direction each update
		slimeSprite := &g.enemies[i]
		updateMovement := rand.IntN(101)
		if updateMovement > 75 && slimeSprite.CurrentVx == 0 && slimeSprite.X%32 == 0 && slimeSprite.Y%32 == 0 {
			slimeSprite.CurrentVx = -1 + rand.IntN(3)
			slimeSprite.CurrentVy = 0
		}
		if updateMovement < 25 && slimeSprite.CurrentVy == 0 && slimeSprite.Y%32 == 0 && slimeSprite.X%32 == 0 {
			slimeSprite.CurrentVx = 0
			slimeSprite.CurrentVy = -1 + rand.IntN(3)
		}

		if !g.checkWallCollision(slimeSprite) {
			slimeSprite.Move(screenWidth, screenHeight)
		}
	}
	if !g.checkWallCollision(&g.playerSprite.CharacterSprite) {
		g.playerSprite.Move(screenWidth, screenHeight)
	}
}

func (g *Game) checkGameEnd() {
	if len(g.eatenBites) >= 8 {
		// Show title
		g.title.Visible = true
		g.title.StartTime = time.Now()
		g.title.WordsVisible = 0
		g.title.Text = "YOU WIN!"
		g.Ended = true
		return
	}
	for _, enemy := range g.enemies {
		if g.playerSprite.CheckCollision(&enemy) {
			// Show title
			g.title.Visible = true
			g.title.StartTime = time.Now()
			g.title.WordsVisible = 0
			g.title.Text = "GAME OVER!"
			g.Ended = true
		}
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
	for i := range g.enemies {
		g.enemies[i].Animate()
	}
	g.currentBite.Animate()
}

var lastAnimationUpdate time.Time

// Update handles the game logic, primarily input and state changes.
func (g *Game) Update() error {

	if time.Since(lastAnimationUpdate) > 100*time.Millisecond {
		lastAnimationUpdate = time.Now()
		g.animate()
	}

	// Example: Exit on pressing Escape
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if g.Ended {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			err := ResetGame()
			if err != nil {
				log.Fatalf("Failed to reset game: %v", err)
			}
		}
		return nil
	}
	g.checkGameEnd()
	g.checkBiteEaten()

	g.handleInputAndMovement()

	return nil
}

func (g *Game) checkBiteEaten() {
	if g.playerSprite.CheckCollision(g.currentBite) {
		alreadyEaten := false
		for _, bite := range g.eatenBites {
			if g.currentBite.Id == bite.Id {
				alreadyEaten = true
				g.playerSprite.Points += 100 * len(g.enemies)
				g.placeNewEnemy()
				break
			}
		}
		if !alreadyEaten {
			g.eatenBites = append(g.eatenBites, *g.currentBite)
			g.playerSprite.Points += 500 + 100*len(g.enemies)
		}
		g.placeNewBite()
	}
}

func ResetGame() error {
	playerImg, err := assets.GetPlayerYellowSprite()
	if err != nil {
		log.Fatalf("failed to load player sprite: %v", err)
	}
	playerSprite := sprites.NewPlayerSprite(playerImg, 32, 32, []sprites.Animation{
		{Name: "right", Frames: 12},
		{Name: "left", Frames: 12},
		{Name: "up", Frames: 12},
		{Name: "down", Frames: 12},
	})

	wallTile, err := assets.GetWallTileImage()
	if err != nil {
		return fmt.Errorf("failed to load wall tile image: %w", err)
	}

	floorTile, err := assets.GetFloorTileImage()
	if err != nil {
		return fmt.Errorf("failed to load floor tile image: %w", err)
	}

	theGame.bites = bites
	theGame.eatenBites = []sprites.CharacterSprite{}
	theGame.bgImage = nil
	theGame.playerSprite = playerSprite
	theGame.enemies = []sprites.CharacterSprite{}
	theGame.mapTiles = [mapHeight][mapWidth]int{
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
	}
	theGame.placeNewBite()
	theGame.placeNewEnemy()
	theGame.wallTile = wallTile
	theGame.floorTile = floorTile
	theGame.title = GameTitle{
		Visible:      true,
		Duration:     5 * time.Second,
		StartTime:    time.Now(),
		WordsVisible: 0,
		ShakeX:       0,
		ShakeY:       0,
		Text:         "8 BITES TO WIN!",
	}
	theGame.Ended = false

	playerSprite.X, playerSprite.Y = theGame.GetRandomFloorPosition()
	//select random tile to spawn slime
	return nil
}

func (g *Game) placeNewBite() {
	g.currentBite = g.bites[rand.IntN(len(theGame.bites))]

	g.currentBite.X, g.currentBite.Y = theGame.GetRandomFloorPosition()
}

func (g *Game) placeNewEnemy() {
	slimeSprite := enemies[rand.IntN(len(enemies))]
	slimeSprite.X, slimeSprite.Y = theGame.GetRandomFloorPosition()
	g.enemies = append(g.enemies, *slimeSprite)
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

	// --- Draw Bites ---
	biteOp := &ebiten.DrawImageOptions{}
	biteOp.GeoM.Translate(float64(g.currentBite.X), float64(g.currentBite.Y))
	biteImg := g.currentBite.GetCurrentImage()
	screen.DrawImage(biteImg, biteOp)

	// --- Draw Player ---
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(float64(g.playerSprite.X), float64(g.playerSprite.Y))
	playerImg := g.playerSprite.GetCurrentImage()
	screen.DrawImage(playerImg, playerOp)

	for _, enemy := range g.enemies {
		slimeOp := &ebiten.DrawImageOptions{}
		slimeOp.GeoM.Translate(float64(enemy.X), float64(enemy.Y))
		slimeImg := enemy.GetCurrentImage()
		screen.DrawImage(slimeImg, slimeOp)
	}

	for i, bite := range g.eatenBites {
		eatenBiteOp := &ebiten.DrawImageOptions{}
		eatenBiteOp.GeoM.Translate(float64(i)*32, 0)
		eatenBiteImg := bite.GetFirstImage()
		screen.DrawImage(eatenBiteImg, eatenBiteOp)
	}

	g.drawScore(screen)

	// Draw title on new level
	if g.title.Visible {
		g.drawTitle(screen)
	}
}

func (g *Game) drawScore(screen *ebiten.Image) {
	t := text.GoTextFace{
		Source: font,
		Size:   16,
	}

	numToDraw := g.playerSprite.Points
	if g.playerSprite.Points > g.playerSprite.LastPoints {
		numToDraw = g.playerSprite.LastPoints
		g.playerSprite.LastPoints += (g.playerSprite.Points-g.playerSprite.LastPoints)/10 + 1
	}
	// draw score with 10 leading zeros
	pointsText := fmt.Sprintf("Score: %010d", numToDraw)
	op := &text.DrawOptions{}
	op.GeoM.Translate(32, float64(screenHeight-24))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, pointsText, &t, op)
}

func (g *Game) drawTitle(screen *ebiten.Image) {

	t := text.GoTextFace{
		Source: font,
		Size:   24,
	}

	words := strings.Split(g.title.Text, " ")
	tw, th := text.Measure(g.title.Text, &t, 0)
	x, y := screenWidth/2-tw/2, screenHeight/2-th/2

	//Draw white block around text
	bgOp := &ebiten.DrawImageOptions{}
	bgOp.GeoM.Translate(float64(x-10), float64(y-10))
	bgRect := ebiten.NewImage(int(tw)+20, int(th)+20)
	//half-transparent background
	bgRect.Fill(color.RGBA{220, 220, 225, 0})
	screen.DrawImage(bgRect, bgOp)

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.Gray{})
	if time.Since(g.title.LastShakeTime) > 50*time.Millisecond {
		g.title.ShakeX = rand.IntN(6) - 3
		g.title.ShakeY = rand.IntN(6) - 3
		g.title.LastShakeTime = time.Now()
	}
	op.GeoM.Translate(float64(g.title.ShakeX), float64(g.title.ShakeY))
	for i := 0; i < g.title.WordsVisible && i < len(words); i++ {
		word := words[i]
		wordWidth, _ := text.Measure(word+" ", &t, 0)
		text.Draw(screen, word+" ", &t, op)
		op.GeoM.Translate(float64(wordWidth), 0)
	}
	if time.Since(g.title.StartTime) > time.Second*time.Duration(g.title.WordsVisible) && g.title.WordsVisible < len(words) {
		g.title.WordsVisible++
	}
	if time.Since(g.title.StartTime) > g.title.Duration {
		g.title.Visible = false
	}
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
