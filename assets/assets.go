package assets

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// The simulated data for the requested path: assets/sprites/player/yellow.png
	// This represents a simple 32x32 blue square.
	AssetsSpritesPlayerYellowPNGData = "iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAAB4sgndAAAABmJLR0QA/wD/AP+gvaeTAAAAZElEQVRYw+3XQREAMAAA0P/l3Q1uT8g4nE/YVv2Q+gHQCQCXF9+1/zQAwGk/0PqfK75uNn7F+gEADgCw5/h0+l51AADYDwDYZ/b3E98XAQB+AIDz+Jp//gBwAQB+AAB2/QC398G2o89hLwAAAABJRU5ErkJggg=="

	// Background image data remains for the background.
	BackgroundPNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAIAAAACABAMAAAAxE/umAAAAGFBMVEUAAABmM+0oMO4nNO0nM+4mNuwkNOsnMuwQn6+YAAAAEklEQVRo3u3asQ0AAACDof6D5v8K3sA2mDGM4zhcI5qA1/0W0Qj1TjYJAAAAAAAAAIDQWw2Vqg7F4Xv3B3sAAAAAAAAAAMDfRk989pX4rAAAAABJRU5ErkJggg=="
)

//go:embed sprites/items/cheese.png
//go:embed sprites/items/donut.png
//go:embed sprites/player/yellow.png
//go:embed sprites/npc/slime.png
//go:embed sprites/world/wall.png
//go:embed sprites/world/floor.png
var folder embed.FS

func GetPlayerYellowSprite() (*ebiten.Image, error) {
	return GetSprite("player/yellow.png")
}

func GetSlimeSprite() (*ebiten.Image, error) {
	return GetSprite("npc/slime.png")
}

func GetFloorTileImage() (*ebiten.Image, error) {
	return GetSprite("world/floor.png")
}

func GetWallTileImage() (*ebiten.Image, error) {
	return GetSprite("world/wall.png")
}

func GetItem(item string) (*ebiten.Image, error) {
	return GetSprite("items/" + item + ".png")
}

func GetSprite(path string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFileSystem(folder, "sprites/"+path)
	if err != nil {
		return nil, err
	}
	return img, nil
}
