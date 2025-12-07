package sprites

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteId int

const (
	SpriteIdPlayer SpriteId = iota
	SpriteIdSlime
	SpriteIdCheese
	SpriteIdDonut
	SpriteIdPizza
	SpriteIdSushi
)

type CharacterSprite struct {
	Image            *ebiten.Image
	Width            int
	Height           int
	X                int
	Y                int
	Frames           int
	CurrentAnimation int
	Animations       []Animation
	CurrentFrame     int
	CurrentVx        int
	CurrentVy        int
	Id               SpriteId
}

func (s *CharacterSprite) CheckCollision(sprite *CharacterSprite) bool {
	ownRect := ebiten.NewImage(s.Width, s.Height).Bounds()
	ownRect = ownRect.Add(image.Point{X: s.X, Y: s.Y})

	spriteRect := ebiten.NewImage(sprite.Width, sprite.Height).Bounds()
	spriteRect = spriteRect.Add(image.Point{X: sprite.X, Y: sprite.Y})
	return ownRect.Overlaps(spriteRect)
}

type Animation struct {
	Name   string
	Frames int
}

func NewCharacterSprite(img *ebiten.Image, width, height int, animations []Animation, id SpriteId) *CharacterSprite {
	s := &CharacterSprite{
		Image:            img,
		Width:            width,
		Height:           height,
		Frames:           img.Bounds().Dx() / width,
		Animations:       animations,
		CurrentAnimation: 0,
		Id:               id,
	}
	return s
}

func (s *CharacterSprite) SetAnimation(animation string) {
	for i, anim := range s.Animations {
		if anim.Name == animation {
			if s.CurrentAnimation != i {
				s.CurrentAnimation = i
				s.CurrentFrame = 0
			}
			return
		}
	}
}

func (s *CharacterSprite) Move(screenWidth, screenHeight int) {
	s.X += s.CurrentVx
	s.Y += s.CurrentVy

	// Basic boundary check to keep the player on screen
	if s.X < 0 {
		s.X = 0
	}
	if s.Y < 0 {
		s.Y = 0
	}
	// Ensure player doesn't move past the right edge
	if s.X > screenWidth-s.Width {
		s.X = screenWidth - s.Width
	}
	// Ensure player doesn't move past the bottom edge
	if s.Y > screenHeight-s.Height {
		s.Y = screenHeight - s.Height
	}
}

func (s *CharacterSprite) Animate() {
	s.CurrentFrame++
	if s.CurrentFrame >= s.Animations[s.CurrentAnimation].Frames {
		s.CurrentFrame = 0
	}
}

func (s *CharacterSprite) GetCurrentImage() *ebiten.Image {
	x := s.CurrentFrame * s.Width
	y := s.CurrentAnimation * s.Height
	return s.Image.SubImage(image.Rect(x, y, x+s.Width, y+s.Height)).(*ebiten.Image)
}

func (s *CharacterSprite) GetFirstImage() *ebiten.Image {
	return s.Image.SubImage(image.Rect(0, 0, s.Width, s.Height)).(*ebiten.Image)
}
