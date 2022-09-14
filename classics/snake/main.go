/*******************************************************************************************
*
*   raylib-go - classic game: snake
*
*   Sample game developed by Ian Eito, Albert Martos and Ramon Santamaria
*	Transliterating to Go by Panagiotis Georgiadis
*
*   This game has been created using raylib-go -- Golang bindings for raylib
*   raylib-go is licensed under an unmodified zlib/libpng license
*
*   Copyright (c) 2022 Panagiotis Georgiadis (drpaneas)
*
********************************************************************************************/
package main

import (
	"image/color"

	. "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

// ----------------------------------------------------------------------------------
// Some Defines
// ----------------------------------------------------------------------------------
const SNAKE_LENGTH int = 256
const SQUARE_SIZE int = 31

// ----------------------------------------------------------------------------------
// Types and Structures Definition
// ----------------------------------------------------------------------------------
type Snake struct {
	position Vector2
	size     Vector2
	speed    Vector2
	color    Color
}

type Food struct {
	position Vector2
	size     Vector2
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

var fruit = Food{}
var snake = [SNAKE_LENGTH]Snake{}
var snakePosition = [SNAKE_LENGTH]Vector2{}
var allowMove = false
var offset = Vector2{}
var counterTail = 0

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization (Note windowTitle is unused on Android)
	//---------------------------------------------------------
	InitWindow(screenWidth, screenHeight, "classic game: snake")

	InitGame()

	SetTargetFPS(60)

	for !WindowShouldClose() {
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
	framesCounter = 0
	gameOver = false
	pause = false

	counterTail = 1
	allowMove = false

	offset.X = float32(screenHeight % SQUARE_SIZE)
	offset.Y = float32(screenHeight % SQUARE_SIZE)

	for i := 0; i < SNAKE_LENGTH; i++ {
		snake[i].position = Vector2{X: offset.X / 2, Y: offset.Y / 2}
		snake[i].size = Vector2{X: float32(SQUARE_SIZE), Y: float32(SQUARE_SIZE)}
		snake[i].speed = Vector2{X: float32(SQUARE_SIZE)}

		if i == 0 {
			snake[i].color = DarkBlue
		} else {
			snake[i].color = Blue
		}
	}

	for i := 0; i < SNAKE_LENGTH; i++ {
		snakePosition[i] = Vector2{}
	}

	fruit.size = Vector2{X: float32(SQUARE_SIZE), Y: float32(SQUARE_SIZE)}
	fruit.color = SkyBlue
	fruit.active = false
}

// Update game (one frame)
func UpdateGame() {
	if !gameOver {
		if IsKeyPressed(KeyP) {
			pause = !pause
		}

		if !pause {
			// Player control
			if IsKeyPressed(KeyRight) && (snake[0].speed.X == 0) && allowMove {
				snake[0].speed = Vector2{X: float32(SQUARE_SIZE)}
				allowMove = false
			}
			if IsKeyPressed(KeyLeft) && (snake[0].speed.X == 0) && allowMove {
				snake[0].speed = Vector2{X: float32(-SQUARE_SIZE)}
				allowMove = false
			}
			if IsKeyPressed(KeyUp) && (snake[0].speed.Y == 0) && allowMove {
				snake[0].speed = Vector2{Y: float32(-SQUARE_SIZE)}
				allowMove = false
			}
			if IsKeyPressed(KeyDown) && (snake[0].speed.Y == 0) && allowMove {
				snake[0].speed = Vector2{Y: float32(SQUARE_SIZE)}
				allowMove = false
			}

			// Snake movement
			for i := 0; i < counterTail; i++ {
				snakePosition[i] = snake[i].position
			}

			if (framesCounter % 5) == 0 {
				for i := 0; i < counterTail; i++ {
					if i == 0 {
						snake[0].position.X += snake[0].speed.X
						snake[0].position.Y += snake[0].speed.Y
						allowMove = true
					} else {
						snake[i].position = snakePosition[i-1]
					}
				}
			}

			// Wall behaviour
			if ((snake[0].position.X) > (screenWidth - offset.X)) || ((snake[0].position.Y) > (screenHeight - offset.Y)) || (snake[0].position.X < 0) || (snake[0].position.Y < 0) {
				gameOver = true
			}

			// Collision with yourself
			for i := 1; i < counterTail; i++ {
				if (snake[0].position.X == snake[i].position.X) && (snake[0].position.Y == snake[i].position.Y) {
					gameOver = true
				}
			}

			// Fruit position calculation
			if !fruit.active {
				fruit.active = true
				fruit.position = Vector2{
					X: float32(GetRandomValue(0, (screenWidth/int32(SQUARE_SIZE))-1)*int32(SQUARE_SIZE) + int32(offset.X/2)),
					Y: float32(GetRandomValue(0, (screenHeight/int32(SQUARE_SIZE))-1)*int32(SQUARE_SIZE) + int32(offset.Y/2)),
				}

				for i := 0; i < counterTail; i++ {
					for (fruit.position.X == snake[i].position.X) && (fruit.position.Y == snake[i].position.Y) {
						fruit.position = Vector2{
							X: float32(GetRandomValue(0, (screenWidth/int32(SQUARE_SIZE))-1)*int32(SQUARE_SIZE) + int32(offset.X/2)),
							Y: float32(GetRandomValue(0, (screenHeight/int32(SQUARE_SIZE))-1)*int32(SQUARE_SIZE) + int32(offset.Y/2)),
						}
						i = 0
					}
				}
			}

			// Collision
			if (snake[0].position.X < (fruit.position.X+fruit.size.X) && (snake[0].position.X+snake[0].size.X) > fruit.position.X) && (snake[0].position.Y < (fruit.position.Y+fruit.size.Y) && (snake[0].position.Y+snake[0].size.Y) > fruit.position.Y) {
				snake[counterTail].position = snakePosition[counterTail-1]
				counterTail += 1
				fruit.active = false
			}

			framesCounter++
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
		// Draw grid lines
		for i := 0; i < screenWidth/SQUARE_SIZE+1; i++ {
			DrawLineV(
				Vector2{X: float32(SQUARE_SIZE*i) + offset.X/2, Y: offset.Y / 2},
				Vector2{X: float32(SQUARE_SIZE*i) + offset.X/2, Y: screenHeight - offset.Y/2},
				LightGray)
		}

		for i := 0; i < screenHeight/SQUARE_SIZE+1; i++ {
			DrawLineV(
				Vector2{X: offset.X / 2, Y: float32(SQUARE_SIZE*i) + offset.Y/2},
				Vector2{X: screenWidth - offset.X/2, Y: float32(SQUARE_SIZE*i) + offset.Y/2},
				LightGray)
		}

		// Draw snake
		for i := 0; i < counterTail; i++ {
			DrawRectangleV(snake[i].position, snake[i].size, snake[i].color)
		}

		// Draw fruit to pick
		DrawRectangleV(fruit.position, fruit.size, fruit.color)

		if pause {
			DrawText("GAME PAUSED", screenWidth/2-MeasureText("GAME PAUSED", 40)/2, screenHeight/2-40, 40, Gray)
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

// drawText It's the same as rl.DrawText but works with any Number type, to avoid type casting pollution.
func drawText[T Number](text string, posX, posY, fontSize T, col color.RGBA) {
	DrawText(text, int32(posX), int32(posY), int32(fontSize), col)
}

// measureText It's the same as rl.MeasureText but works with any Number type, to avoid type casting pollution.
func measureText[T Number](text string, fontSize T) T {
	return T(MeasureText(text, int32(fontSize)))
}
