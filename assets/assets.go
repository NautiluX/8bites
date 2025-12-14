package assets

import (
	"bytes"
	"embed"
	"fmt"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// The simulated data for the requested path: assets/sprites/player/yellow.png
	// This represents a simple 32x32 blue square.
	AssetsSpritesPlayerYellowPNGData = "iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAAB4sgndAAAABmJLR0QA/wD/AP+gvaeTAAAAZElEQVRYw+3XQREAMAAA0P/l3Q1uT8g4nE/YVv2Q+gHQCQCXF9+1/zQAwGk/0PqfK75uNn7F+gEADgCw5/h0+l51AADYDwDYZ/b3E98XAQB+AIDz+Jp//gBwAQB+AAB2/QC398G2o89hLwAAAABJRU5ErkJggg=="

	// Background image data remains for the background.
	BackgroundPNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAIAAAACABAMAAAAxE/umAAAAGFBMVEUAAABmM+0oMO4nNO0nM+4mNuwkNOsnMuwQn6+YAAAAEklEQVRo3u3asQ0AAACDof6D5v8K3sA2mDGM4zhcI5qA1/0W0Qj1TjYJAAAAAAAAAIDQWw2Vqg7F4Xv3B3sAAAAAAAAAAMDfRk989pX4rAAAAABJRU5ErkJggg=="
)

//go:embed sprites/**/*.png
//go:embed sfx/*.wav
//go:embed maps/*.txt
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

func GetBackgroundMusic() (*audio.Player, error) {
	reader, err := folder.Open("sfx/backgroundmusic_1.wav")
	if err != nil {
		return nil, err
	}

	stream, err := wav.DecodeWithSampleRate(44100, reader)

	audioContext := audio.NewContext(44100)
	if err != nil {
		return nil, err
	}

	s := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// A simple wrapper type that embeds *bytes.Reader and adds a no-op Close method.
// By embedding *bytes.Reader, it automatically satisfies the io.Reader and io.Seeker interfaces.
type ByteReadSeekCloser struct {
	*bytes.Reader
}

// Close satisfies the io.Closer interface.
func (b ByteReadSeekCloser) Close() error {
	// No resources to clean up for an in-memory byte array, so we return nil.
	return nil
}

// NewReadSeekCloser creates an io.ReadSeekCloser from a byte slice.
func NewReadSeekCloser(data []byte) io.ReadSeekCloser {
	return ByteReadSeekCloser{
		Reader: bytes.NewReader(data),
	}
}

func GetMapTiles(name string) ([15][20]int, error) {
	file, err := folder.Open("maps/" + name + ".txt")
	if err != nil {
		return [15][20]int{}, err
	}
	defer file.Close()

	var tiles [15][20]int
	var buf = make([]byte, 1)
	for y := 0; y < 15; y++ {
		for x := 0; x < 20; x++ {
			_, err := file.Read(buf)
			if err != nil {
				return [15][20]int{}, fmt.Errorf("error reading map file: %w", err)
			}
			for buf[0] == '\n' || buf[0] == ' ' {
				_, err := file.Read(buf)
				if err != nil {
					return [15][20]int{}, fmt.Errorf("error reading map file: %w", err)
				}
			}
			fmt.Print(string(buf))
			var tile int
			_, err = fmt.Sscanf(string(buf), "%d", &tile)
			if err != nil {
				return [15][20]int{}, err
			}
			tiles[y][x] = tile
		}
	}
	return tiles, nil
}
