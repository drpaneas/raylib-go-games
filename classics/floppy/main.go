/*******************************************************************************************
*
*   raylib-go - classic game: floppy
*
*   Sample game developed by Ian Eito, Albert Martos and Ramon Santamaria
*	Transliterated to Go by Panagiotis Georgiadis
*
*   This game has been created using using raylib-go -- Golang bindings for raylib
*   raylib-go is licensed under an unmodified zlib/libpng license
*
*   Copyright (c) 2022 Panagiotis Georgiadis (drpaneas)
*
********************************************************************************************/
package main

import (
	"fmt"
	. "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
	"image/color"
)

// ----------------------------------------------------------------------------------
// Some Defines
// ----------------------------------------------------------------------------------
const MAX_TUBES = 100
const FLOPPY_RADIUS = 24
const TUBES_WIDTH = 80

// ----------------------------------------------------------------------------------
// Types and Structures Definition
// ----------------------------------------------------------------------------------
type Floppy struct {
	position Vector2
	radius   int
	color    Color
}

type Tubes struct {
	rec    Rectangle
	color  Color
	active bool
}

// ------------------------------------------------------------------------------------
// Global Variables Declaration
// ------------------------------------------------------------------------------------
const screenWidth = 800
const screenHeight = 450

var gameOver bool = false
var pause bool = false
var score int = 0
var hiScore int = 0

var floppy Floppy
var tubes [MAX_TUBES * 2]Tubes
var tubesPos [MAX_TUBES]Vector2
var tubesSpeedX int
var superfx bool

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	InitWindow(screenWidth, screenHeight, "classic game: floppy")
	InitGame()
	SetTargetFPS(60)

	// Main game loop
	for !WindowShouldClose() { // Detect window close button or ESC key
		UpdateDrawFrame()
	}

	// De-Initialization
	CloseWindow() // Close window and OpenGL context
}

//------------------------------------------------------------------------------------
// Module Functions Definitions (local)
//------------------------------------------------------------------------------------

// Initialize game variables
func InitGame() {
	floppy.radius = FLOPPY_RADIUS
	floppy.position = Vector2{X: 80, Y: float32(screenHeight/2 - floppy.radius)}
	tubesSpeedX = 2

	for i := 0; i < MAX_TUBES; i++ {
		tubesPos[i].X = 400 + 280*float32(i)
		tubesPos[i].Y = float32(-GetRandomValue(0, 120))
	}

	for i := 0; i < MAX_TUBES*2; i += 2 {
		tubes[i].rec.X = tubesPos[i/2].X
		tubes[i].rec.Y = tubesPos[i/2].Y
		tubes[i].rec.Width = TUBES_WIDTH
		tubes[i].rec.Height = 255

		tubes[i+1].rec.X = tubesPos[i/2].X
		tubes[i+1].rec.Y = 600 + tubesPos[i/2].Y - 255
		tubes[i+1].rec.Width = TUBES_WIDTH
		tubes[i+1].rec.Height = 255

		tubes[i/2].active = true
	}

	score = 0

	gameOver = false
	superfx = false
	pause = false
}

// Update game (one frame)
func UpdateGame() {
	if !gameOver {
		if IsKeyPressed(KeyP) {
			pause = !pause
		}
		if !pause {
			for i := 0; i < MAX_TUBES; i++ {
				tubesPos[i].X -= float32(tubesSpeedX)
			}
			for i := 0; i < MAX_TUBES*2; i += 2 {
				tubes[i].rec.X = tubesPos[i/2].X
				tubes[i+1].rec.X = tubesPos[i/2].X
			}
			if IsKeyDown(KeySpace) && !gameOver {
				floppy.position.Y -= 3
			} else {
				floppy.position.Y += 1
			}
			// Check Collisions
			for i := 0; i < MAX_TUBES*2; i++ {
				if CheckCollisionCircleRec(floppy.position, float32(floppy.radius), tubes[i].rec) {
					gameOver = true
					pause = false
				} else if tubesPos[i/2].X < floppy.position.X && tubes[i/2].active && !gameOver {
					score += 100
					tubes[i/2].active = false
					superfx = true
					if score > hiScore {
						hiScore = score
					}
				}
			}
		}
	} else {
		if IsKeyPressed(KeyEnter) {
			InitGame()
			gameOver = false
		}
	}
}

// Draw game (one frame)
func DrawGame() {
	BeginDrawing()
	ClearBackground(RayWhite)

	if !gameOver {
		drawCircle(floppy.position.X, floppy.position.Y, floppy.radius, DarkGray)

		// Draw tubes
		for i := 0; i < MAX_TUBES; i++ {
			drawRectangle(tubes[i*2].rec.X, tubes[i*2].rec.Y, tubes[i*2].rec.Width, tubes[i*2].rec.Height, Gray)
			drawRectangle(tubes[i*2+1].rec.X, tubes[i*2+1].rec.Y, tubes[i*2+1].rec.Width, tubes[i*2+1].rec.Height, Gray)
		}

		// Draw flashing fx (one frame only)
		if superfx {
			drawRectangle(0, 0, screenWidth, screenHeight, White)
			superfx = false
		}

		drawText(fmt.Sprintf("%04d", score), 20, 20, 40, Gray)
		drawText(fmt.Sprintf("HI-SCORE: %04d", hiScore), 20, 70, 20, LightGray)

		if pause {
			drawText("GAME PAUSED", screenWidth/2-MeasureText("GAME PAUSED", 40)/2, screenHeight/2-40, 40, Gray)
		}
	} else {
		drawText("PRESS [ENTER] TO PLAY AGAIN", GetScreenWidth()/2-measureText("PRESS [ENTER] TO PLAY AGAIN", 20)/2, GetScreenHeight()/2-50, 20, Gray)
	}

	EndDrawing()
}

// Update and Draw (one frame)
func UpdateDrawFrame() {
	UpdateGame()
	DrawGame()
}

// ------------------------------------------------------------------------------------.
// Generics refactoring.
// ------------------------------------------------------------------------------------.

// Number is a constraint that permits any Integer and Floating-point type.
type Number interface {
	constraints.Integer | constraints.Float
}

// drawRectangle It's the same as rl.DrawRectangle but works with any Number type, to avoid type casting pollution.
func drawRectangle[T Number](posX, posY, width, height T, col color.RGBA) {
	DrawRectangle(int32(posX), int32(posY), int32(width), int32(height), col)
}

// drawCircle It's the same as rl.DrawCircle but works with any Number type, to avoid type casting pollution.
func drawCircle[T Number, N Number](centerX, centerY T, radius N, col color.RGBA) {
	DrawCircle(int32(centerX), int32(centerY), float32(radius), col)
}

// drawText It's the same as rl.DrawText but works with any Number type, to avoid type casting pollution.
func drawText[T Number](text string, posX, posY, fontSize T, col color.RGBA) {
	DrawText(text, int32(posX), int32(posY), int32(fontSize), col)
}

// measureText It's the same as rl.MeasureText but works with any Number type, to avoid type casting pollution.
func measureText[T Number](text string, fontSize T) T {
	return T(MeasureText(text, int32(fontSize)))
}
