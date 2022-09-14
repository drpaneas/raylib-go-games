/*******************************************************************************************
*
*   raylib-go - classic game: gold fever
*
*   Sample game developed by Ian Eito, Albert Martos and Ramon Santamaria
*   Transliterated to Go by Panagiotis Georgiadis
*
*   This game has been created using raylib-go -- Golang bindings for raylib
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
// Types and Structures Definition
// ----------------------------------------------------------------------------------
type Player struct {
	position Vector2
	speed    Vector2
	radius   int
}

type Enemy struct {
	position     Vector2
	speed        Vector2
	radius       int
	radiusBounds int
	moveRight    bool
}

type Points struct {
	position Vector2
	radius   int
	value    int
	active   bool
}

type Home struct {
	rec    Rectangle
	active bool
	save   bool
	color  Color
}

// ------------------------------------------------------------------------------------
// Global Variables Declaration
// ------------------------------------------------------------------------------------
const (
	screenWidth  = 800
	screenHeight = 450
)

var gameOver bool
var pause bool
var score int
var hiScore int

var player Player
var enemy Enemy
var points Points
var home Home
var follow bool

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	InitWindow(screenWidth, screenHeight, "classic game: gold fever")
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
	pause = false
	score = 0

	player.position = Vector2{X: 50, Y: 50}
	player.radius = 20
	player.speed = Vector2{X: 5, Y: 5}

	enemy.position = Vector2{X: screenWidth - 50, Y: screenHeight / 2}
	enemy.radius = 20
	enemy.radiusBounds = 150
	enemy.speed = Vector2{X: 3, Y: 3}
	enemy.moveRight = true
	follow = false

	points.radius = 10
	points.position = Vector2{
		X: float32(GetRandomValue(int32(points.radius), int32(screenWidth-points.radius))),
		Y: float32(GetRandomValue(int32(points.radius), int32(screenHeight-points.radius)))}
	points.value = 100
	points.active = true

	home.rec.Width = 50
	home.rec.Height = 50
	home.rec.X = float32(GetRandomValue(0, int32(screenWidth-home.rec.Width)))
	home.rec.Y = float32(GetRandomValue(0, int32(screenHeight-home.rec.Height)))
	home.active = false
	home.save = false
}

// Update game (one frame)
func UpdateGame() {
	if !gameOver {
		if IsKeyPressed(KeyP) {
			pause = !pause
		}
		if !pause {
			if IsKeyDown(KeyRight) {
				player.position.X += player.speed.X
			}
			if IsKeyDown(KeyLeft) {
				player.position.X -= player.speed.X
			}
			if IsKeyDown(KeyUp) {
				player.position.Y -= player.speed.Y
			}
			if IsKeyDown(KeyDown) {
				player.position.Y += player.speed.Y
			}

			if int(player.position.X)-player.radius <= 0 {
				player.position.X = float32(player.radius)
			}
			if int(player.position.X)+player.radius >= screenWidth {
				player.position.X = float32(screenWidth - player.radius)
			}
			if int(player.position.Y)-player.radius <= 0 {
				player.position.Y = float32(player.radius)
			}
			if int(player.position.Y)+player.radius >= screenHeight {
				player.position.Y = float32(screenHeight - player.radius)
			}

			if (follow || CheckCollisionCircles(player.position, float32(player.radius), enemy.position, float32(enemy.radiusBounds))) && !home.save {
				if player.position.X > enemy.position.X {
					enemy.position.X += enemy.speed.X
				}
				if player.position.X < enemy.position.X {
					enemy.position.X -= enemy.speed.X
				}
				if player.position.Y > enemy.position.Y {
					enemy.position.Y += enemy.speed.Y
				}
				if player.position.Y < enemy.position.Y {
					enemy.position.Y -= enemy.speed.Y
				}
			} else {
				if enemy.moveRight {
					enemy.position.X += enemy.speed.X
				} else {
					enemy.position.X -= enemy.speed.X
				}
			}

			if int(enemy.position.X)-enemy.radius <= 0 {
				enemy.moveRight = true
			}
			if int(enemy.position.X)+enemy.radius >= screenWidth {
				enemy.moveRight = false
			}
			if int(enemy.position.X)-enemy.radius <= 0 {
				enemy.position.X = float32(enemy.radius)
			}
			if int(enemy.position.X)+enemy.radius >= screenWidth {
				enemy.position.X = float32(screenWidth - enemy.radius)
			}
			if int(enemy.position.Y)-enemy.radius <= 0 {
				enemy.position.Y = float32(enemy.radius)
			}
			if int(enemy.position.Y)+enemy.radius >= screenHeight {
				enemy.position.Y = float32(screenHeight - enemy.radius)
			}

			if CheckCollisionCircles(player.position, float32(player.radius), points.position, float32(points.radius)) && points.active {
				follow = true
				points.active = false
				home.active = true
			}
			if CheckCollisionCircles(player.position, float32(player.radius), enemy.position, float32(enemy.radius)) && !home.save {
				gameOver = true
				if hiScore < score {
					hiScore = score
				}
			}
			if CheckCollisionCircleRec(player.position, float32(player.radius), home.rec) {
				follow = false
				if !points.active {
					score += points.value
					points.active = true
					enemy.speed.X += 0.5
					enemy.speed.Y += 0.5
					points.position = Vector2{
						X: float32(GetRandomValue(int32(points.radius), int32(screenWidth-points.radius))),
						Y: float32(GetRandomValue(int32(points.radius), int32(screenHeight-points.radius)))}
				}
				home.save = true
			} else {
				home.save = false
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
		if follow {
			DrawRectangle(0, 0, screenWidth, screenHeight, Red)
			DrawRectangle(10, 10, screenWidth-20, screenHeight-20, RayWhite)
		}

		drawRectangleLines(home.rec.X, home.rec.Y, home.rec.Width, home.rec.Height, Blue)

		drawCircleLines(enemy.position.X, enemy.position.Y, enemy.radiusBounds, Red)
		drawCircleV(enemy.position, enemy.radius, Maroon)

		drawCircleV(player.position, player.radius, Gray)
		if points.active {
			drawCircleV(points.position, points.radius, Gold)
		}

		drawText(fmt.Sprintf("SCORE: %04d", score), 20, 15, 20, Gray)
		drawText(fmt.Sprintf("HI-SCORE: %04d", hiScore), 300, 15, 20, Gray)

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

// drawCircleV It's the same as rl.DrawRectangle but works with any Number type, to avoid type casting pollution.
func drawCircleV[T Number](center Vector2, radius T, col color.RGBA) {
	DrawCircleV(center, float32(radius), col)
}

// drawText It's the same as rl.DrawText but works with any Number type, to avoid type casting pollution.
func drawText[T Number](text string, posX, posY, fontSize T, col color.RGBA) {
	DrawText(text, int32(posX), int32(posY), int32(fontSize), col)
}

// measureText It's the same as rl.MeasureText but works with any Number type, to avoid type casting pollution.
func measureText[T Number](text string, fontSize T) T {
	return T(MeasureText(text, int32(fontSize)))
}

// drawRectangleLines It's the same as rl.DrawRectangleLines but works with any Number type, to avoid type casting pollution.
func drawRectangleLines[T, N, U, W Number](posX T, posY N, width U, height W, col color.RGBA) {
	DrawRectangleLines(int32(posX), int32(posY), int32(width), int32(height), col)
}

// drawCircleLines It's the same as rl.DrawCircleLines but works with any Number type, to avoid type casting pollution.
func drawCircleLines[T, N, U Number](centerX T, centerY N, radius U, col color.RGBA) {
	DrawCircleLines(int32(centerX), int32(centerY), float32(radius), col)
}
