// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	ebitenImage       *ebiten.Image
	ebitenImageWidth  = 0
	ebitenImageHeight = 0
)

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
}

func (s *Sprite) Update() {
	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if screenWidth <= s.x+s.imageWidth {
		s.x = 2*(screenWidth-s.imageWidth) - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if screenHeight <= s.y+s.imageHeight {
		s.y = 2*(screenHeight-s.imageHeight) - s.y
		s.vy = -s.vy
	}
}

type Sprites struct {
	sprites []*Sprite
	num     int
}

func (s *Sprites) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
}

const (
	MinSprites = 0
	MaxSprites = 50000
)

var sprites = &Sprites{make([]*Sprite, MaxSprites), 500}

var op *ebiten.DrawImageOptions

func init() {
	op = &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1.0, 1.0, 1.0, 0.5)
}

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		sprites.num -= 20
		if sprites.num < MinSprites {
			sprites.num = MinSprites
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		sprites.num += 20
		if MaxSprites < sprites.num {
			sprites.num = MaxSprites
		}
	}
	sprites.Update()

	if ebiten.IsRunningSlowly() {
		return nil
	}
	for i := 0; i < sprites.num; i++ {
		s := sprites.sprites[i]
		op.GeoM = ebiten.GeoM{}
		op.GeoM.Translate(float64(s.x), float64(s.y))
		screen.DrawImage(ebitenImage, op)
	}
	msg := fmt.Sprintf(`FPS: %0.2f
Num of sprites: %d
Press <- or -> to change the number of sprites`, ebiten.CurrentFPS(), sprites.num)
	if err := ebitenutil.DebugPrint(screen, msg); err != nil {
		return err
	}
	return nil
}

func main() {
	var err error
	ebitenImage, _, err = ebitenutil.NewImageFromFile("_resources/images/ebiten.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	ebitenImageWidth, ebitenImageHeight = ebitenImage.Size()
	for i := range sprites.sprites {
		w, h := ebitenImage.Size()
		x, y := rand.Intn(screenWidth-w), rand.Intn(screenHeight-h)
		vx, vy := 2*rand.Intn(2)-1, 2*rand.Intn(2)-1
		sprites.sprites[i] = &Sprite{
			imageWidth:  w,
			imageHeight: h,
			x:           x,
			y:           y,
			vx:          vx,
			vy:          vy,
		}
	}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Sprites (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
