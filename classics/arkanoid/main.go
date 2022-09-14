package main

import (
	. "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
	"image/color"
	"math"
)

// ----------------------------------------------------------------------------------
// Some Defines
// ----------------------------------------------------------------------------------
const PLAYER_MAX_LIFE = 5
const LINES_OF_BRICKS = 5
const BRICKS_PER_LINE = 20

type Player struct {
	position Vector2
	size     Vector2
	life     int
}

type Ball struct {
	position Vector2
	speed    Vector2
	radius   int
	active   bool
}
type Brick struct {
	position Vector2
	active   bool
}

// ------------------------------------------------------------------------------------
// Global Variables Declaration
// ------------------------------------------------------------------------------------
const screenWidth = 800
const screenHeight = 450

var gameOver bool = false
var pause bool = false

var player Player
var ball Ball
var brick [LINES_OF_BRICKS][BRICKS_PER_LINE]Brick
var brickSize Vector2

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization (Note windowTitle is unused on Android)
	//---------------------------------------------------------
	InitWindow(screenWidth, screenHeight, "classic game: arkanoid")

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
	brickSize = Vector2{X: float32(GetScreenWidth() / BRICKS_PER_LINE), Y: 40}

	// Initialize player
	player = Player{
		position: Vector2{X: screenWidth / 2, Y: screenHeight * 7 / 8},
		size:     Vector2{X: screenWidth / 10, Y: 20},
		life:     PLAYER_MAX_LIFE,
	}

	// Initialize ball
	ball = Ball{
		position: Vector2{X: screenWidth / 2, Y: screenHeight*7/8 - 30},
		speed:    Vector2{},
		radius:   7,
		active:   false,
	}

	// Initialize bricks
	initialDownPosition := 50

	for i := 0; i < LINES_OF_BRICKS; i++ {
		for j := 0; j < BRICKS_PER_LINE; j++ {
			brick[i][j] = Brick{
				position: Vector2{
					X: float32(j)*brickSize.X + brickSize.X/2,
					Y: float32(i)*brickSize.Y + float32(initialDownPosition),
				},
				active: true,
			}
		}
	}
}

// Update game (one frame)
func UpdateGame() {
	if !gameOver {
		if IsKeyPressed('P') {
			pause = !pause
		}

		if !pause {
			// Player movement logic
			if IsKeyDown(KeyLeft) {
				player.position.X -= 5
			}

			if (player.position.X - player.size.X/2) <= 0 {
				player.position.X = player.size.X / 2
			}

			if IsKeyDown(KeyRight) {
				player.position.X += 5
			}

			if (player.position.X + player.size.X/2) >= screenWidth {
				player.position.X = screenWidth - player.size.X/2
			}

			// Ball launching logic
			if !ball.active {
				if IsKeyPressed(KeySpace) {
					ball.active = true
					ball.speed = Vector2{Y: -5}
				}
			}

			// Ball movement logic
			if ball.active {
				ball.position.X += ball.speed.X
				ball.position.Y += ball.speed.Y
			} else {
				ball.position = Vector2{X: player.position.X, Y: screenHeight*7/8 - 30}
			}

			// Collision logic: ball vs walls
			if ((int(ball.position.X) + ball.radius) >= screenWidth) || ((int(ball.position.X) - ball.radius) <= 0) {
				ball.speed.X *= -1
			}
			if (int(ball.position.Y) - ball.radius) <= 0 {
				ball.speed.Y *= -1
			}
			if (int(ball.position.Y) + ball.radius) >= screenHeight {
				ball.speed = Vector2{}
				ball.active = false

				player.life--
			}

			// Collision logic: ball vs player
			if CheckCollisionCircleRec(ball.position, float32(ball.radius),
				Rectangle{X: player.position.X - player.size.X/2, Y: player.position.Y - player.size.Y/2, Width: player.size.X, Height: player.size.Y}) {
				if ball.speed.Y > 0 {
					ball.speed.Y *= -1
					ball.speed.X = (ball.position.X - player.position.X) / (player.size.X / 2) * 5
				}
			}

			// Collision logic: ball vs bricks
			for i := 0; i < LINES_OF_BRICKS; i++ {
				for j := 0; j < BRICKS_PER_LINE; j++ {
					if brick[i][j].active {

						if ((int(ball.position.Y) - ball.radius) <= int(brick[i][j].position.Y+brickSize.Y/2)) &&
							((int(ball.position.Y) - ball.radius) > int(brick[i][j].position.Y+brickSize.Y/2+ball.speed.Y)) &&
							(int(fabs(ball.position.X-brick[i][j].position.X)) < (int(brickSize.X)/2 + ball.radius*2/3)) && (ball.speed.Y < 0) {
							// Hit below
							brick[i][j].active = false
							ball.speed.Y *= -1
						} else if ((int(ball.position.Y) + ball.radius) >= int(brick[i][j].position.Y-brickSize.Y/2)) &&
							((int(ball.position.Y) + ball.radius) < int(brick[i][j].position.Y-brickSize.Y/2+ball.speed.Y)) &&
							(int(fabs(ball.position.X-brick[i][j].position.X)) < (int(brickSize.X)/2 + ball.radius*2/3)) && (ball.speed.Y > 0) {
							// Hit above
							brick[i][j].active = false
							ball.speed.Y *= -1
						} else if ((int(ball.position.X) + ball.radius) >= int(brick[i][j].position.X-brickSize.X/2)) &&
							((int(ball.position.X) + ball.radius) < int(brick[i][j].position.X-brickSize.X/2+ball.speed.X)) &&
							(int(fabs(ball.position.Y-brick[i][j].position.Y)) < (int(brickSize.Y)/2 + ball.radius*2/3)) && (ball.speed.X > 0) {
							// Hit Left
							brick[i][j].active = false
							ball.speed.X *= -1
						} else if ((int(ball.position.X) - ball.radius) <= int(brick[i][j].position.X+brickSize.X/2)) &&
							((int(ball.position.X) - ball.radius) > int(brick[i][j].position.X+brickSize.X/2+ball.speed.X)) &&
							(int(fabs(ball.position.Y-brick[i][j].position.Y)) < (int(brickSize.Y)/2 + ball.radius*2/3)) && (ball.speed.X < 0) {
							// Hit Right
							brick[i][j].active = false
							ball.speed.X *= -1
						}
					}
				}
			}

			// Game over logic
			if player.life <= 0 {
				gameOver = true
			} else {
				gameOver = true

				for i := 0; i < LINES_OF_BRICKS; i++ {
					for j := 0; j < BRICKS_PER_LINE; j++ {
						if brick[i][j].active {
							gameOver = false
						}
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
		// Draw player bar
		drawRectangle(player.position.X-player.size.X/2, player.position.Y-player.size.Y/2, player.size.X, player.size.Y, Black)

		// Draw player lives
		for i := 0; i < player.life; i++ {
			drawRectangle(20+40*i, screenHeight-30, 35, 10, LightGray)
		}

		// Draw ball
		drawCircleV(ball.position, ball.radius, Maroon)

		// Draw bricks
		for i := 0; i < LINES_OF_BRICKS; i++ {
			for j := 0; j < BRICKS_PER_LINE; j++ {
				if brick[i][j].active {
					posX := brick[i][j].position.X - brickSize.X/2
					posY := brick[i][j].position.Y - brickSize.Y/2
					width := brickSize.X
					height := brickSize.Y

					if (i+j)%2 == 0 {
						drawRectangle(posX, posY, width, height, Gray)
					} else {
						drawRectangle(posX, posY, width, height, DarkGray)
					}
				}
			}
		}

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

// fabs It's the same as math.Abs but works with any Number type, to avoid type casting pollution.
func fabs[T Number](n T) T {
	return T(math.Abs(float64(n)))
}

// drawRectangle It's the same as rl.DrawRectangle but works with any Number type, to avoid type casting pollution.
func drawRectangle[T Number](posX, posY, width, height T, col color.RGBA) {
	DrawRectangle(int32(posX), int32(posY), int32(width), int32(height), col)
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
