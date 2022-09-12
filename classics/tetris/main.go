/*******************************************************************************************
*
*   raylib - classic game: tetris
*
*   Sample game developed by Marc Palau and Ramon Santamaria
*   Transliterated to Go by Panagiotis Georgiadis
*
*   This game has been created using raylib-go -- Golang bindings for raylib
*   raylib is licensed under an unmodified zlib/libpng license (View raylib.h for details)
*
*   Copyright (c) 2022 Panagiotis Georgiadis (drpaneas)
*
********************************************************************************************/
package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ----------------------------------------------------------------------------------
// Some Defines
// ----------------------------------------------------------------------------------
const (
	SquareSize           = 20
	GridHorizontalSize   = 12
	VerticalSize         = 20
	LateralSpeed         = 10
	TurningSpeed         = 12
	FastFallAwaitCounter = 30
	FadingTime           = 33
)

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------

type GridSquare int

const (
	EMPTY GridSquare = iota
	MOVING
	FULL
	BLOCK
	FADING
)

// ------------------------------------------------------------------------------------
// Global Variables Declaration
// ------------------------------------------------------------------------------------
const screenWidth = 800
const screenHeight = 450

var (
	gameOver = false
	pause    = false
)

// Matrices
var grid [GridHorizontalSize][VerticalSize]GridSquare

var (
	piece         [4][4]GridSquare
	incomingPiece [4][4]GridSquare
)

// These variables keep track of the active piece position
var piecePositionX = 0
var piecePositionY = 0

// Game parameters
var fadingColor rl.Color

var (
	beginPlay    = true // This var is only true at the beginning of the game, used for the first matrix creations
	pieceActive  = false
	detection    = false
	lineToDelete = false
)

// Statistics
var level = 1
var score = 0

// Counters
var gravityMovementCounter = 0

var (
	lateralMovementCounter  = 0
	turnMovementCounter     = 0
	fastFallMovementCounter = 0
	fadeLineCounter         = 0
)

// Based on level
var gravitySpeed = level * 30

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	// Initialization (Note windowTitle is unused on Android)
	//---------------------------------------------------------
	rl.InitWindow(screenWidth, screenHeight, "classic game: tetris")

	InitGame()

	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		// Update and Draw
		//----------------------------------------------------------------------------------
		UpdateDrawFrame()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.CloseWindow() // Close window and OpenGL context
	//--------------------------------------------------------------------------------------
}

//--------------------------------------------------------------------------------------
// Game Module Functions Definition
//--------------------------------------------------------------------------------------

// InitGame Initialize game variables
func InitGame() {
	// Initialize game statistics
	level = 1
	score = 0

	fadingColor = rl.Gray

	piecePositionX = 0
	piecePositionY = 0

	pause = false

	beginPlay = true
	pieceActive = false
	detection = false
	lineToDelete = false

	// Counters
	gravityMovementCounter = 0
	lateralMovementCounter = 0
	turnMovementCounter = 0
	fastFallMovementCounter = 0

	fadeLineCounter = 0
	gravitySpeed = 30

	// Initialize grid matrices
	for i := 0; i < GridHorizontalSize; i++ {
		for j := 0; j < VerticalSize; j++ {
			if (j == VerticalSize-1) || (i == 0) || (i == GridHorizontalSize-1) {
				grid[i][j] = BLOCK
			} else {
				grid[i][j] = EMPTY
			}
		}
	}

	// Initialize incoming piece matrices
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			incomingPiece[i][j] = EMPTY
		}
	}
}

// UpdateGame Update game (one frame)
func UpdateGame() {
	if !gameOver {
		if rl.IsKeyPressed(rl.KeyP) {
			pause = !pause
		}

		if !pause {
			if !lineToDelete {
				if !pieceActive {
					// Get another piece
					pieceActive = CreatePiece()

					// We leave a little time before starting the fast falling down
					fastFallMovementCounter = 0
				} else { // Piece falling
					// Counters update
					fastFallMovementCounter++
					gravityMovementCounter++
					lateralMovementCounter++
					turnMovementCounter++

					// We make sure to move if we've pressed the key this frame
					if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight) {
						lateralMovementCounter = LateralSpeed
					}
					if rl.IsKeyPressed(rl.KeyUp) {
						turnMovementCounter = TurningSpeed
					}

					// Fall down
					if rl.IsKeyDown(rl.KeyDown) && (fastFallMovementCounter >= FastFallAwaitCounter) {
						// We make sure the piece is going to fall this frame
						gravityMovementCounter += gravitySpeed
					}

					if gravityMovementCounter >= gravitySpeed {
						// Basic falling movement
						CheckDetection(&detection)

						// Check if the piece has collided with another piece or with the bounding
						ResolveFallingMovement(&detection, &pieceActive)

						// Check if we fulfilled a line and if so, erase the line and pull down the line above
						CheckCompletion(&lineToDelete)

						gravityMovementCounter = 0
					}

					// Move laterally at player's will
					if lateralMovementCounter >= LateralSpeed {
						// Update the lateral movement and if success, reset the lateral counter
						if !ResolveLateralMovement() {
							lateralMovementCounter = 0
						}
					}

					// Turn the piece at player's will
					if turnMovementCounter >= TurningSpeed {
						// Update the turning movement and reset the turning counter
						if ResolveTurnMovement() {
							turnMovementCounter = 0
						}
					}
				}

				// Game over logic
				for j := 0; j < 2; j++ {
					for i := 1; i < GridHorizontalSize-1; i++ {
						if grid[i][j] == FULL {
							gameOver = true
						}
					}
				}

			} else {
				// Animation when deleting score
				fadeLineCounter++

				if fadeLineCounter%8 < 4 {
					fadingColor = rl.Maroon
				} else {
					fadingColor = rl.Gray
				}

				if fadeLineCounter >= FadingTime {
					deletedLines := DeleteCompleteLines()
					fadeLineCounter = 0
					lineToDelete = false

					score += deletedLines
				}
			}
		}
	} else {
		if rl.IsKeyPressed(rl.KeyEnter) {
			InitGame()
			gameOver = false
		}
	}
}

// DrawGame Draw game (one frame)
func DrawGame() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	if !gameOver {
		// Draw gameplay area
		offset := rl.Vector2{
			X: screenWidth/2 - (GridHorizontalSize * SquareSize) - 50,
			Y: screenHeight/2 - ((VerticalSize - 1) * SquareSize / 2) + SquareSize*2,
		}

		offset.Y -= 50 // NOTE: Hardcoded position!

		controller := offset.X

		for j := 0; j < VerticalSize; j++ {
			for i := 0; i < GridHorizontalSize; i++ {
				// Draw each square of the grid
				switch grid[i][j] {
				case EMPTY:
					DrawLine(offset.X, offset.Y, offset.X+SquareSize, offset.Y, rl.LightGray)
					DrawLine(offset.X, offset.Y, offset.X, offset.Y+SquareSize, rl.LightGray)
					DrawLine(offset.X+SquareSize, offset.Y, offset.X+SquareSize, offset.Y+SquareSize, rl.DarkGray)
					DrawLine(offset.X, offset.Y+SquareSize, offset.X+SquareSize, offset.Y+SquareSize, rl.DarkGray)
					offset.X += SquareSize
				case FULL:
					DrawRectangle(offset.X, offset.Y, SquareSize, SquareSize, rl.Gray)
					offset.X += SquareSize
				case MOVING:
					DrawRectangle(offset.X, offset.Y, SquareSize, SquareSize, rl.DarkGray)
					offset.X += SquareSize
				case BLOCK:
					DrawRectangle(offset.X, offset.Y, SquareSize, SquareSize, rl.LightGray)
					offset.X += SquareSize
				case FADING:
					DrawRectangle(offset.X, offset.Y, SquareSize, SquareSize, fadingColor)
					offset.X += SquareSize
				}
			}

			offset.X = controller
			offset.Y += SquareSize
		}

		// Draw incoming piece (hardcoded)
		offset.X = 500
		offset.Y = 45

		controller = offset.X

		for j := 0; j < 4; j++ {
			for i := 0; i < 4; i++ {
				if incomingPiece[i][j] == EMPTY {
					DrawLine(offset.X, offset.Y, offset.X+SquareSize, offset.Y, rl.LightGray)
					DrawLine(offset.X, offset.Y, offset.X, offset.Y+SquareSize, rl.LightGray)
					DrawLine(offset.X+SquareSize, offset.Y, offset.X+SquareSize, offset.Y+SquareSize, rl.LightGray)
					DrawLine(offset.X, offset.Y+SquareSize, offset.X+SquareSize, offset.Y+SquareSize, rl.LightGray)
					offset.X += SquareSize
				} else if incomingPiece[i][j] == MOVING {
					DrawRectangle(offset.X, offset.Y, SquareSize, SquareSize, rl.Gray)
					offset.X += SquareSize
				}
			}

			offset.X = controller
			offset.Y += SquareSize
		}

		DrawText("INCOMING:", offset.X, offset.Y-100, 10, rl.Gray)
		DrawText(fmt.Sprintf("LINES: %04d", score), 500, 250, 20, rl.Gray)

		if pause {
			rl.DrawText("GAME PAUSED", screenWidth/2-rl.MeasureText("GAME PAUSED", 40)/2, screenHeight/2-40, 40, rl.Gray)
		}
	} else {
		const replayMsg = "PRESS [ENTER] TO PLAY AGAIN"
		DrawText(replayMsg, rl.GetScreenWidth()/2-MeasureText(replayMsg, 20)/2, rl.GetScreenHeight()/2-50, 20, rl.Gray)
	}

	rl.EndDrawing()
}

// UpdateDrawFrame Update and Draw (one frame)
func UpdateDrawFrame() {
	UpdateGame()
	DrawGame()
}

//--------------------------------------------------------------------------------------
// Additional module functions
//--------------------------------------------------------------------------------------

func CreatePiece() bool {
	piecePositionX = (GridHorizontalSize - 4) / 2
	piecePositionY = 0

	// If the game is starting, and you are going to create the first piece, we create an extra one
	if beginPlay {
		GetRandomPiece()
		beginPlay = false
	}

	// We assign the incoming piece to the actual piece
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			piece[i][j] = incomingPiece[i][j]
		}
	}

	// We assign a random piece to the incoming one
	GetRandomPiece()

	// Assign the piece to the grid
	for i := piecePositionX; i < piecePositionX+4; i++ {
		for j := 0; j < 4; j++ {
			if piece[i-piecePositionX][j] == MOVING {
				grid[i][j] = MOVING
			}
		}
	}

	return true
}

// GetRandomPiece Get a random piece
func GetRandomPiece() {
	random := rl.GetRandomValue(0, 6)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			incomingPiece[i][j] = EMPTY
		}
	}

	switch random {
	case 0: // O Square
		incomingPiece[1][1] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING
	case 1: // L
		incomingPiece[1][0] = MOVING
		incomingPiece[1][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING
	case 2: // L inverse
		incomingPiece[1][2] = MOVING
		incomingPiece[2][0] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[2][2] = MOVING
	case 3: // I Line
		incomingPiece[0][1] = MOVING
		incomingPiece[1][1] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[3][1] = MOVING
	case 4: // T Piece
		incomingPiece[1][0] = MOVING
		incomingPiece[1][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][1] = MOVING
	case 5: // Z
		incomingPiece[1][1] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[2][2] = MOVING
		incomingPiece[3][2] = MOVING
	case 6: // S or Z inverse
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[3][1] = MOVING
	}
}

func CheckDetection(detection *bool) {
	for j := VerticalSize - 1; j >= 0; j-- {
		for i := 0; i < GridHorizontalSize-1; i++ {
			if grid[i][j] == MOVING && (grid[i][j+1] == FULL || grid[i][j+1] == BLOCK) {
				*detection = true
			}
		}
	}
}

func CheckCompletion(lineToDelete *bool) {
	var calculator int

	for j := VerticalSize - 2; j >= 0; j-- {
		calculator = 0
		for i := 1; i < GridHorizontalSize-1; i++ {
			// Count each square of the line
			if grid[i][j] == FULL {
				calculator++
			}

			// Check if we completed the whole line
			if calculator == GridHorizontalSize-2 {
				*lineToDelete = true
				calculator = 0

				// Mark the completed line
				for z := 1; z < GridHorizontalSize-1; z++ {
					grid[z][j] = FADING
				}
			}
		}
	}
}

func ResolveFallingMovement(detection, pieceActive *bool) {
	// If we finished moving this piece, we stop it
	if *detection {
		for j := VerticalSize - 2; j >= 0; j-- {
			for i := 1; i < GridHorizontalSize-1; i++ {
				if grid[i][j] == MOVING {
					grid[i][j] = FULL
					*detection = false
					*pieceActive = false
				}
			}
		}
	} else { // We move down the piece
		for j := VerticalSize - 2; j >= 0; j-- {
			for i := 1; i < GridHorizontalSize-1; i++ {
				if grid[i][j] == MOVING {
					grid[i][j+1] = MOVING
					grid[i][j] = EMPTY
				}
			}
		}

		piecePositionY++
	}
}

func DeleteCompleteLines() int {
	var deletedLines int

	// Erase the completed line
	for j := VerticalSize - 2; j >= 0; j-- {
		for grid[1][j] == FADING {
			for i := 1; i < GridHorizontalSize-1; i++ {
				grid[i][j] = EMPTY
			}

			// Erase the completed line by relocating all the current lines of the grid down
			// otherwise there will be a gap with EMPTY cells
			for j2 := j - 1; j2 >= 0; j2-- {
				for i2 := 1; i2 < GridHorizontalSize-1; i2++ {
					if grid[i2][j2] == FULL {
						grid[i2][j2+1] = FULL
						grid[i2][j2] = EMPTY
					} else if grid[i2][j2] == FADING {
						grid[i2][j2+1] = FADING
						grid[i2][j2] = EMPTY
					}
				}
			}

			deletedLines++
		}
	}

	return deletedLines
}

func ResolveLateralMovement() bool {
	collision := false

	// Piece movement
	if rl.IsKeyDown(rl.KeyLeft) { // Move left
		// Check if is possible to move to the left
		for j := VerticalSize - 2; j >= 0; j-- {
			for i := 1; i < GridHorizontalSize-1; i++ {
				if grid[i][j] == MOVING {
					// Check if we are touching the left wall, or we have a full square at the left
					if i-1 == 0 || grid[i-1][j] == FULL {
						collision = true
					}
				}
			}
		}

		// If able, move left
		if !collision {
			for j := VerticalSize - 2; j >= 0; j-- {
				for i := 1; i < GridHorizontalSize-1; i++ { // We check the matrix from left to right
					// Move everything to the left
					if grid[i][j] == MOVING {
						grid[i-1][j] = MOVING
						grid[i][j] = EMPTY
					}
				}
			}

			piecePositionX--
		}
	} else if rl.IsKeyDown(rl.KeyRight) { // Move Right
		// Check if is possible to move to right
		for j := VerticalSize - 2; j >= 0; j-- {
			for i := 1; i < GridHorizontalSize-1; i++ {
				if grid[i][j] == MOVING {
					// Check if we are touching the right wall, or we have a full square at the right
					if i+1 == GridHorizontalSize-1 || grid[i+1][j] == FULL {
						collision = true
					}
				}
			}
		}

		// If able, move right
		if !collision {
			for j := VerticalSize - 2; j >= 0; j-- {
				for i := GridHorizontalSize - 1; i >= 1; i-- { // We check the matrix from right to left
					// Move everything to the right
					if grid[i][j] == MOVING {
						grid[i+1][j] = MOVING
						grid[i][j] = EMPTY
					}
				}
			}

			piecePositionX++
		}
	}

	return collision
}

func ResolveTurnMovement() bool {
	// Input for turning the piece
	if rl.IsKeyDown(rl.KeyUp) {
		var aux GridSquare
		var checker bool

		// Check all turning possibilities
		if (grid[piecePositionX+3][piecePositionY] == MOVING) &&
			(grid[piecePositionX][piecePositionY] != EMPTY) &&
			(grid[piecePositionX][piecePositionY] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+3][piecePositionY+3] == MOVING) &&
			(grid[piecePositionX+3][piecePositionY] != EMPTY) &&
			(grid[piecePositionX+3][piecePositionY] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX][piecePositionY+3] == MOVING) &&
			(grid[piecePositionX+3][piecePositionY+3] != EMPTY) &&
			(grid[piecePositionX+3][piecePositionY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX][piecePositionY] == MOVING) &&
			(grid[piecePositionX][piecePositionY+3] != EMPTY) &&
			(grid[piecePositionX][piecePositionY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+1][piecePositionY] == MOVING) &&
			(grid[piecePositionX][piecePositionY+2] != EMPTY) &&
			(grid[piecePositionX][piecePositionY+2] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+3][piecePositionY+1] == MOVING) &&
			(grid[piecePositionX+1][piecePositionY] != EMPTY) &&
			(grid[piecePositionX+1][piecePositionY] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+2][piecePositionY+3] == MOVING) &&
			(grid[piecePositionX+3][piecePositionY+1] != EMPTY) &&
			(grid[piecePositionX+3][piecePositionY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX][piecePositionY+2] == MOVING) &&
			(grid[piecePositionX+2][piecePositionY+3] != EMPTY) &&
			(grid[piecePositionX+2][piecePositionY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+2][piecePositionY] == MOVING) &&
			(grid[piecePositionX][piecePositionY+1] != EMPTY) &&
			(grid[piecePositionX][piecePositionY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+3][piecePositionY+2] == MOVING) &&
			(grid[piecePositionX+2][piecePositionY] != EMPTY) &&
			(grid[piecePositionX+2][piecePositionY] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+1][piecePositionY+3] == MOVING) &&
			(grid[piecePositionX+3][piecePositionY+2] != EMPTY) &&
			(grid[piecePositionX+3][piecePositionY+2] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX][piecePositionY+1] == MOVING) &&
			(grid[piecePositionX+1][piecePositionY+3] != EMPTY) &&
			(grid[piecePositionX+1][piecePositionY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+1][piecePositionY+1] == MOVING) &&
			(grid[piecePositionX+1][piecePositionY+2] != EMPTY) &&
			(grid[piecePositionX+1][piecePositionY+2] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+2][piecePositionY+1] == MOVING) &&
			(grid[piecePositionX+1][piecePositionY+1] != EMPTY) &&
			(grid[piecePositionX+1][piecePositionY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+2][piecePositionY+2] == MOVING) &&
			(grid[piecePositionX+2][piecePositionY+1] != EMPTY) &&
			(grid[piecePositionX+2][piecePositionY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePositionX+1][piecePositionY+2] == MOVING) &&
			(grid[piecePositionX+2][piecePositionY+2] != EMPTY) &&
			(grid[piecePositionX+2][piecePositionY+2] != MOVING) {
			checker = true
		}

		if !checker {
			aux = piece[0][0]
			piece[0][0] = piece[3][0]
			piece[3][0] = piece[3][3]
			piece[3][3] = piece[0][3]
			piece[0][3] = aux

			aux = piece[1][0]
			piece[1][0] = piece[3][1]
			piece[3][1] = piece[2][3]
			piece[2][3] = piece[0][2]
			piece[0][2] = aux

			aux = piece[2][0]
			piece[2][0] = piece[3][2]
			piece[3][2] = piece[1][3]
			piece[1][3] = piece[0][1]
			piece[0][1] = aux

			aux = piece[1][1]
			piece[1][1] = piece[2][1]
			piece[2][1] = piece[2][2]
			piece[2][2] = piece[1][2]
			piece[1][2] = aux
		}

		for j := VerticalSize - 2; j >= 0; j-- {
			for i := 1; i < GridHorizontalSize-1; i++ {
				if grid[i][j] == MOVING {
					grid[i][j] = EMPTY
				}
			}
		}

		for i := piecePositionX; i < piecePositionX+4; i++ {
			for j := piecePositionY; j < piecePositionY+4; j++ {
				if piece[i-piecePositionX][j-piecePositionY] == MOVING {
					grid[i][j] = MOVING
				}
			}
		}

		return true
	}

	return false
}

// ------------------------------------------------------------------------------------.
// Generics refactoring.
// ------------------------------------------------------------------------------------.

// Number is a constraint that permits any Integer and Floating-point type.
type Number interface {
	constraints.Integer | constraints.Float
}

// DrawLine It's the same as rl.DrawLine but works with any Number type, to avoid type casting pollution.
func DrawLine[T Number](startPosX, startPosY, endPosX, endPosY T, col color.RGBA) {
	rl.DrawLine(int32(startPosX), int32(startPosY), int32(endPosX), int32(endPosY), col)
}

// DrawRectangle It's the same as rl.DrawRectangle but works with any Number type, to avoid type casting pollution.
func DrawRectangle[T Number](posX, posY, width, height T, col color.RGBA) {
	rl.DrawRectangle(int32(posX), int32(posY), int32(width), int32(height), col)
}

// DrawText It's the same as rl.DrawText but works with any Number type, to avoid type casting pollution.
func DrawText[T Number](text string, posX, posY, fontSize T, col color.RGBA) {
	rl.DrawText(text, int32(posX), int32(posY), int32(fontSize), col)
}

// MeasureText It's the same as rl.MeasureText but works with any Number type, to avoid type casting pollution.
func MeasureText[T Number](text string, fontSize T) T {
	return T(rl.MeasureText(text, int32(fontSize)))
}
