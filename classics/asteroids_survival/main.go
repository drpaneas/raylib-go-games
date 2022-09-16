/*******************************************************************************************
*
*   raylib-go - classic game: asteroids survival
*
*   Sample game developed by Ian Eito, Albert Martos and Ramon Santamaria
*	Translated to Go by Panagiotis Georgiadis
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
	"math"
)

// ----------------------------------------------------------------------------------
// Some Defines
// ----------------------------------------------------------------------------------
const (
	PLAYER_BASE_SIZE   = 20.0
	PLAYER_SPEED       = 6.0
	METEORS_SPEED      = 2
	MAX_MEDIUM_METEORS = 8
	MAX_SMALL_METEORS  = 16
)

// ----------------------------------------------------------------------------------
// Types and Structures Definition
// ----------------------------------------------------------------------------------
type Player struct {
	position     Vector2
	speed        Vector2
	acceleration float32
	rotation     float32
	collider     Vector3
	color        Color
}
type Meteor struct {
	position Vector2
	speed    Vector2
	radius   float32
	active   bool
	color    Color
}

// ------------------------------------------------------------------------------------
// Global Variables Declaration
// ------------------------------------------------------------------------------------
const screenWidth = 800
const screenHeight = 450

var framesCounter = 0
var gameOver = false
var pause = false

// NOTE: Defined triangle is isosceles with common angles of 70 degrees.
var shipHeight float32 = 0.0

var player Player
var mediumMeteor [MAX_MEDIUM_METEORS]Meteor
var smallMeteor [MAX_SMALL_METEORS]Meteor

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization (Note windowTitle is unused on Android)
	//---------------------------------------------------------
	InitWindow(screenWidth, screenHeight, "classic game: asteroids survival")

	InitGame()
	SetTargetFPS(60)

	// Main game loop
	for !WindowShouldClose() { // Detect window close button or ESC key
		// Update and Draw
		//----------------------------------------------------------------------------------
		UpdateDrawFrame()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}

//------------------------------------------------------------------------------------
// Module Functions Definitions (local)
//------------------------------------------------------------------------------------

// Initialize game variables
func InitGame() {
	var posx, posy int
	var velx, vely int
	var correctRange bool

	pause = false

	framesCounter = 0

	shipHeight = float32((PLAYER_BASE_SIZE / 2) / math.Tan(20*Deg2rad))
	// Initialization player
	player.position = newVector2(screenWidth/2, screenHeight/2-shipHeight/2)
	player.speed = Vector2{}
	player.acceleration = 0
	player.rotation = 0
	player.collider = newVector3(
		player.position.X+sin(player.rotation*Deg2rad)*(shipHeight/2.5),
		player.position.Y-cos(player.rotation*Deg2rad)*(shipHeight/2.5),
		12)
	player.color = LightGray

	for i := 0; i < MAX_MEDIUM_METEORS; i++ {
		posx = getRandomValue(0, screenWidth)

		for !correctRange {
			if posx > screenWidth/2-150 && posx < screenWidth/2+150 {
				posx = getRandomValue(0, screenWidth)
			} else {
				correctRange = true
			}
		}

		correctRange = false

		posy = getRandomValue(0, screenHeight)

		for !correctRange {
			if posy > screenHeight/2-150 && posy < screenHeight/2+150 {
				posy = getRandomValue(0, screenHeight)
			} else {
				correctRange = true
			}
		}

		correctRange = false
		velx = getRandomValue(-METEORS_SPEED, METEORS_SPEED)
		vely = getRandomValue(-METEORS_SPEED, METEORS_SPEED)

		for !correctRange {
			if velx == 0 && vely == 0 {
				velx = getRandomValue(-METEORS_SPEED, METEORS_SPEED)
				vely = getRandomValue(-METEORS_SPEED, METEORS_SPEED)
			} else {
				correctRange = true
			}
		}
		mediumMeteor[i].position = newVector2(posx, posy)
		mediumMeteor[i].speed = newVector2(velx, vely)
		mediumMeteor[i].radius = 20
		mediumMeteor[i].active = true
		mediumMeteor[i].color = Green
	}

	for i := 0; i < MAX_SMALL_METEORS; i++ {
		posx := getRandomValue(0, screenWidth)
		for !correctRange {
			if posx > screenWidth/2-150 && posx < screenWidth/2+150 {
				posx = getRandomValue(0, screenWidth)
			} else {
				correctRange = true
			}
		}
		correctRange = false
		posy := getRandomValue(0, screenHeight)
		for !correctRange {
			if posy > screenHeight/2-150 && posy < screenHeight/2+150 {
				posy = getRandomValue(0, screenHeight)
			} else {
				correctRange = true
			}
		}
		correctRange = false
		velx := getRandomValue(-METEORS_SPEED, METEORS_SPEED)
		vely := getRandomValue(-METEORS_SPEED, METEORS_SPEED)
		for !correctRange {
			if velx == 0 && vely == 0 {
				velx = getRandomValue(-METEORS_SPEED, METEORS_SPEED)
				vely = getRandomValue(-METEORS_SPEED, METEORS_SPEED)
			} else {
				correctRange = true
			}
		}
		smallMeteor[i].position = newVector2(posx, posy)
		smallMeteor[i].speed = newVector2(velx, vely)
		smallMeteor[i].radius = 10
		smallMeteor[i].active = true
		smallMeteor[i].color = Yellow
	}
}

// Update game (one frame)
func UpdateGame() {
	if !gameOver {
		if IsKeyPressed(KeyP) {
			pause = !pause
		}

		if !pause {
			framesCounter++

			// Player logic

			// Rotation
			if IsKeyDown(KeyLeft) {
				player.rotation -= 5
			}

			if IsKeyDown(KeyRight) {
				player.rotation += 5
			}

			// Speed
			player.speed.X = sin(player.rotation*Deg2rad) * PLAYER_SPEED
			player.speed.Y = cos(player.rotation*Deg2rad) * PLAYER_SPEED

			// Controller
			if IsKeyDown(KeyUp) {
				if player.acceleration < 1 {
					player.acceleration += 0.04
				}
			} else {
				if player.acceleration > 0 {
					player.acceleration -= 0.02
				} else if player.acceleration < 0 {
					player.acceleration = 0
				}
			}

			if IsKeyDown(KeyDown) {
				if player.acceleration > 0 {
					player.acceleration -= 0.04
				} else if player.acceleration < 0 {
					player.acceleration = 0
				}
			}

			// Movement
			player.position.X += player.speed.X * player.acceleration
			player.position.Y -= player.speed.Y * player.acceleration

			// Wall behaviour for player
			if player.position.X > screenWidth+shipHeight {
				player.position.X = -shipHeight
			} else if player.position.X < -shipHeight {
				player.position.X = screenWidth + shipHeight
			}

			if player.position.Y > screenHeight+shipHeight {
				player.position.Y = -shipHeight
			} else if player.position.Y < -shipHeight {
				player.position.Y = screenHeight + shipHeight
			}

			// Collision Player to meteors
			player.collider = newVector3(
				player.position.X+sin(player.rotation*Deg2rad)*(shipHeight/2.5),
				player.position.Y-cos(player.rotation*Deg2rad)*(shipHeight/2.5),
				12)

			for a := 0; a < MAX_MEDIUM_METEORS; a++ {
				if CheckCollisionCircles(newVector2(player.collider.X, player.collider.Y), player.collider.Z, mediumMeteor[a].position, mediumMeteor[a].radius) && mediumMeteor[a].active {
					gameOver = true
				}
			}

			for a := 0; a < MAX_SMALL_METEORS; a++ {
				if CheckCollisionCircles(newVector2(player.collider.X, player.collider.Y), player.collider.Z, smallMeteor[a].position, smallMeteor[a].radius) && smallMeteor[a].active {
					gameOver = true
				}
			}

			// Meteor logic
			for i := 0; i < MAX_MEDIUM_METEORS; i++ {
				if mediumMeteor[i].active {
					// Movement
					mediumMeteor[i].position.X += mediumMeteor[i].speed.X
					mediumMeteor[i].position.Y += mediumMeteor[i].speed.Y

					// wall behaviour
					if mediumMeteor[i].position.X > screenWidth+mediumMeteor[i].radius {
						mediumMeteor[i].position.X = -mediumMeteor[i].radius
					} else if mediumMeteor[i].position.X < 0-mediumMeteor[i].radius {
						mediumMeteor[i].position.X = screenWidth + mediumMeteor[i].radius
					}

					if mediumMeteor[i].position.Y > screenHeight+mediumMeteor[i].radius {
						mediumMeteor[i].position.Y = -mediumMeteor[i].radius
					} else if mediumMeteor[i].position.Y < 0-mediumMeteor[i].radius {
						mediumMeteor[i].position.Y = screenHeight + mediumMeteor[i].radius
					}
				}
			}

			for i := 0; i < MAX_SMALL_METEORS; i++ {
				if smallMeteor[i].active {
					// Movement
					smallMeteor[i].position.X += smallMeteor[i].speed.X
					smallMeteor[i].position.Y += smallMeteor[i].speed.Y

					// wall behaviour
					if smallMeteor[i].position.X > screenWidth+smallMeteor[i].radius {
						smallMeteor[i].position.X = -smallMeteor[i].radius
					} else if smallMeteor[i].position.X < 0-smallMeteor[i].radius {
						smallMeteor[i].position.X = screenWidth + smallMeteor[i].radius
					}

					if smallMeteor[i].position.Y > screenHeight+smallMeteor[i].radius {
						smallMeteor[i].position.Y = -smallMeteor[i].radius
					} else if smallMeteor[i].position.Y < 0-smallMeteor[i].radius {
						smallMeteor[i].position.Y = screenHeight + smallMeteor[i].radius
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
		// Draw Spaceship
		v1 := newVector2(player.position.X+sin(player.rotation*Deg2rad)*(shipHeight), player.position.Y-cos(player.rotation*Deg2rad)*(shipHeight))
		v2 := newVector2(player.position.X-cos(player.rotation*Deg2rad)*(PLAYER_BASE_SIZE/2), player.position.Y-sin(player.rotation*Deg2rad)*(PLAYER_BASE_SIZE/2))
		v3 := newVector2(player.position.X+cos(player.rotation*Deg2rad)*(PLAYER_BASE_SIZE/2), player.position.Y+sin(player.rotation*Deg2rad)*(PLAYER_BASE_SIZE/2))
		DrawTriangle(v1, v2, v3, Maroon)

		// Draw meteor
		for i := 0; i < MAX_MEDIUM_METEORS; i++ {
			if mediumMeteor[i].active {
				DrawCircleV(mediumMeteor[i].position, mediumMeteor[i].radius, Gray)
			} else {
				DrawCircleV(mediumMeteor[i].position, mediumMeteor[i].radius, Fade(LightGray, 0.3))
			}
		}

		for i := 0; i < MAX_SMALL_METEORS; i++ {
			if smallMeteor[i].active {
				DrawCircleV(smallMeteor[i].position, smallMeteor[i].radius, DarkGray)
			} else {
				DrawCircleV(smallMeteor[i].position, smallMeteor[i].radius, Fade(LightGray, 0.3))
			}
		}

		DrawText(fmt.Sprintf("TIME: %.02v", framesCounter/60), 10, 10, 20, Black)

		if pause {
			DrawText("GAME PAUSED", screenWidth/2-MeasureText("GAME PAUSED", 40)/2, screenHeight/2-40, 40, Gray)
		}
	} else {
		drawText("PRESS [ENTER] TO PLAY AGAIN", GetScreenWidth()/2-measureText("PRESS [ENTER] TO PLAY AGAIN", 20)/2, GetScreenHeight()/2-50, 20, Gray)
	}

	EndDrawing()
	//----------------------------------------------------------------------------------
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

// drawText It's the same as rl.DrawText but works with any Number type, to avoid type casting pollution.
func drawText[T Number](text string, posX, posY, fontSize T, col color.RGBA) {
	DrawText(text, int32(posX), int32(posY), int32(fontSize), col)
}

// measureText It's the same as rl.MeasureText but works with any Number type, to avoid type casting pollution.
func measureText[T Number](text string, fontSize T) T {
	return T(MeasureText(text, int32(fontSize)))
}

// newVector2 - It's the same as NewVector2 but works with any Number type, to avoid type casting pollution.
func newVector2[T, U Number](x T, y U) Vector2 {
	return Vector2{X: float32(x), Y: float32(y)}
}

// newVector3 -  It's the same as NewVector3 but works with any Number type, to avoid type casting pollution.
func newVector3[T, U, S Number](X T, Y U, Z S) Vector3 {
	return Vector3{X: float32(X), Y: float32(Y), Z: float32(Z)}
}

func sin[T Number](x T) T {
	return T(math.Sin(float64(x)))
}

func cos[T Number](x T) T {
	return T(math.Cos(float64(x)))
}

func getRandomValue[T Number](min, max T) T {
	return T(GetRandomValue(int32(min), int32(max)))
}
