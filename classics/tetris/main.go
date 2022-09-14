/*******************************************************************************************
*
*   raylib-go - classic game: tetris
*
*   Sample game developed by Marc Palau and Ramon Santamaria
*   Ported to Go by Panagiotis Georgiadis
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
	"image/color"

	"golang.org/x/exp/constraints"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ----------------------------------------------------------------------------------
// Some Defines
// ----------------------------------------------------------------------------------
const (
	squareSize           = 20 // Size of the squares that compose the pieces
	gridSizeX            = 12 // 10 + 2 (left and right walls)
	gridSizeY            = 20 // 18 + 2 (top and bottom walls)
	speedTurn            = 12
	fastFallAwaitCounter = 30
	timeToFade           = 33

	// Deceleration: the higher the value, the slower the piece moves
	// It is the number of frames that must pass before the piece moves down one cell
	// NOTE: The higher the value, the slower the piece moves (60 frames --> 1 second)
	framesToWaitBeforeMoveDown        = 30
	framesToWaitBeforeLateralMovement = 10
)

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------

type gridSquare int

const (
	EMPTY gridSquare = iota
	MOVING
	FULL
	BLOCK
	FADING
)

// ------------------------------------------------------------------------------------
// Global Variables Declaration
// ------------------------------------------------------------------------------------

// Resolution of the screen
const (
	screenWidth  = 800
	screenHeight = 450
)

var (
	// Toggle flags
	isGameover      bool // has the game ended?
	isPaused        bool // has the game been paused?
	isFirst         bool // is this the first tetromino piece of the game?
	isPieceFalling  bool // used to know if a piece is active or not.
	isDownCollided  bool // has a piece reached the bottom of the grid or another piece.
	hasLineToDelete bool // used to know if a line has to be deleted.

	// Counters
	verticalMoveCounter   int // Counter used to move the piece down.
	horizontalMoveCounter int // Counter used to move the piece left or right
	turnMovementCounter   int // Counter used to turn the piece.
	fastFallMoveCounter   int // Counter used to move the piece down faster.
	fadeLineCounter       int // Counter used to fade a line.

	// The tetrominos (pieces, bricks, blocks, whatever you want to call them) are stored in a 4x4 matrix.
	piece         [4][4]gridSquare // geometric shape, composed of 4 squares connected orthogonally.
	incomingPiece [4][4]gridSquare // Next tetromino to be created (non-active piece).
	piecePosX     int              // X position of the current active tetromino (in grid squares, not pixels).
	piecePosY     int              // Y position of the current active tetromino (in grid squares, not pixels).

	// Statistics
	score int

	// Grid
	grid [gridSizeX][gridSizeY]gridSquare // Grid area matrix.

	// Colors
	fadingColor rl.Color
)

// init initializes the game the first time
func init() {
	reset()
}

// ------------------------------------------------------------------------------------
// Program main entry point
// ------------------------------------------------------------------------------------
func main() {
	rl.InitWindow(screenWidth, screenHeight, "classic game: tetris")
	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() { // Detect window close button or ESC key
		UpdateDrawFrame()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}

//--------------------------------------------------------------------------------------
// Game Module Functions Definition
//--------------------------------------------------------------------------------------

// reset all the required global variables to their original values.
// It's called when the game starts, and when the player loses (gameover).
func reset() {
	score = 0
	fadingColor = rl.Gray

	// Keep track of the piece that is falling down
	piecePosX = 0
	piecePosY = 0

	// Toggle flags
	isPaused = false
	isFirst = true
	isPieceFalling = false
	isDownCollided = false
	hasLineToDelete = false

	// Counters
	verticalMoveCounter = 0
	horizontalMoveCounter = 0
	turnMovementCounter = 0
	fastFallMoveCounter = 0
	fadeLineCounter = 0

	// We use a 12x20 grid, but we only play the game inside a smaller 10x18 grid, so we have to leave 2 squares
	// empty on each side of the grid.
	//
	//  The grid is composed of 5 types of squares:
	//
	//    1. EMPTY : empty square
	//    2. MOVING: square that is part of the moving tetromino
	//    3. FULL  : square that is part of a tetromino that has reached the bottom of the grid
	//    4. BLOCK : square that is part of the wall perimeter of the grid
	//    5. FADING: square that is part of a tetromino that has reached the bottom of the grid,
	//   			 and it is going to be deleted
	//
	// Initialize the main gaming grid area with empty squares and surrounding walls
	for i := 0; i < gridSizeX; i++ {
		for j := 0; j < gridSizeY; j++ {
			isBottomWall := j == gridSizeY-1
			isLeftWall := i == 0
			isRightWall := i == gridSizeX-1

			if isBottomWall || isLeftWall || isRightWall {
				grid[i][j] = BLOCK // Surrounding Walls
			} else {
				grid[i][j] = EMPTY // Gaming area
			}
		}
	}

	// Initialize preview area grid for incoming piece
	// NOTE: We could use the same grid for the preview area, but we prefer to use a different one
	//       to avoid possible problems with the main grid
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			incomingPiece[i][j] = EMPTY
		}
	}
}

// UpdateGame Update game logic (one frame)
func UpdateGame() {
	// 1. Check if the game is over (if the player has lost)
	if !isGameover {
		// 2. If the game is _not_ over, then check if the game is paused,
		//    and if so, wait for the user to press P to continue
		if rl.IsKeyPressed(rl.KeyP) {
			isPaused = !isPaused
		}
		// 3. If the game is _not_ paused, then proceed to the next step.
		if !isPaused {
			// 4. Check if a line has been completed, and if so, we have to delete it.
			if !hasLineToDelete {
				// 5. If there is no line to delete, then check if the current piece has reached the bottom of the grid.
				if !isPieceFalling {
					// 5a.1 A piece has reached the bottom of the grid, so we have to create a new one.
					isPieceFalling = CreatePiece()
					// 5a.2 In case the user had previously pressed the down key, we have to reset the fastFallMoveCounter
					//     to avoid the piece to fall down too fast.
					fastFallMoveCounter = 0
				} else {
					// 5b.1 If the piece has not reached the bottom of the grid, then check for user input.

					// Counters update
					fastFallMoveCounter++
					verticalMoveCounter++
					horizontalMoveCounter++
					turnMovementCounter++

					// 5b.2 Check if the user has pressed the left or right key to move the piece horizontally.
					if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyRight) {
						horizontalMoveCounter = framesToWaitBeforeLateralMovement
					}

					// 5b.3 Check if the user has pressed the up key to turn the piece.
					if rl.IsKeyPressed(rl.KeyUp) {
						turnMovementCounter = speedTurn
					}

					// 5b.4 Check if the user has pressed the down key to move the piece down faster.
					if rl.IsKeyDown(rl.KeyDown) && (fastFallMoveCounter >= fastFallAwaitCounter) {
						verticalMoveCounter += framesToWaitBeforeMoveDown // Move the piece down faster
					}

					// 5b.5 Check if the number of frames (verticalMoveCounter) has reached the limit, and if so, move the piece down.
					if verticalMoveCounter >= framesToWaitBeforeMoveDown {
						verticalMoveCounter = 0 // Reset the counter

						// 5b.5.1 Check if the piece has collided with the bottom of the grid or with another piece
						isDownCollided = checkCollisionY()
						if isDownCollided {
							stopMovingDown()
						} else {
							moveDown()
						}
						// 5b.5.2 Check if the player has completed a line, and if so, mark it (FADING) to be deleted in the next frame
						CheckCompletion(&hasLineToDelete)
					}

					// 5b.6 Move laterally at player's will
					if horizontalMoveCounter >= framesToWaitBeforeLateralMovement {
						// Update the lateral movement and if success, reset the lateral counter
						if !ResolveLateralMovement() {
							horizontalMoveCounter = 0
						}
					}

					// Turn the piece at player's will
					if turnMovementCounter >= speedTurn {
						// Update the turning movement and reset the turning counter
						if ResolveTurnMovement() {
							turnMovementCounter = 0
						}
					}
				}

				// If the piece has reached the top of the grid, then the game is over.
				for j := 0; j < 2; j++ {
					for i := 1; i < gridSizeX-1; i++ {
						if grid[i][j] == FULL {
							isGameover = true
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

				if fadeLineCounter >= timeToFade {
					deletedLines := DeleteCompleteLines()
					fadeLineCounter = 0
					hasLineToDelete = false

					score += deletedLines
				}
			}
		}
	} else {
		if rl.IsKeyPressed(rl.KeyEnter) {
			reset()
			isGameover = false
		}
	}
}

// moveDown moves the piece down by one cell
// It does that by doing two things:
//  1. Tags the (i, j) current position (MOVING squares) into EMPTY
//  2. Tags the (i, j+1) position equivalent (EMPTY squares) into MOVING
func moveDown() {
	for j := gridSizeY - 2; j >= 0; j-- { // Start from the bottom
		for i := 1; i < gridSizeX-1; i++ { // Start from the left
			if grid[i][j] == MOVING { // If the current cell is MOVING
				grid[i][j+1] = MOVING // Tag the cell below as MOVING
				grid[i][j] = EMPTY    // Tag the current cell as EMPTY
			}
		}
	}

	piecePosY++ // Update piece position information, one cell down
}

// stopMovingDown converts the MOVING squares to FULL and resets the related boolean flags
func stopMovingDown() {
	for j := gridSizeY - 2; j >= 0; j-- { // We start from the bottom of the grid
		for i := 1; i < gridSizeX-1; i++ { // We start from the left side of the grid
			if grid[i][j] == MOVING { // If the square is part of the moving piece
				grid[i][j] = FULL      // Convert it to FULL
				isDownCollided = false // Reset the ground flag
				isPieceFalling = false // Reset the falling flag
			}
		}
	}
}

// DrawGame Draw game (one frame)
func DrawGame() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	if !isGameover {
		// Draw gameplay area
		offset := rl.Vector2{
			// X Offset the grid to the center of the screen
			X: screenWidth/2 - (gridSizeX * squareSize),
			// Y Offset the grid to the bottom of the screen
			Y: screenHeight/2 - ((gridSizeY - 1) * squareSize / 2) + squareSize*2,
		}

		offset.X -= 50 // offset to the left
		offset.Y -= 50 // NOTE: Hardcoded position! Places the bottom of the grid a bit higher

		controller := offset.X

		for j := 0; j < gridSizeY; j++ {
			for i := 0; i < gridSizeX; i++ {
				// Draw each square of the grid
				switch grid[i][j] {
				case EMPTY:
					DrawLine(offset.X, offset.Y, offset.X+squareSize, offset.Y, rl.LightGray)
					DrawLine(offset.X, offset.Y, offset.X, offset.Y+squareSize, rl.LightGray)
					DrawLine(offset.X+squareSize, offset.Y, offset.X+squareSize, offset.Y+squareSize, rl.DarkGray)
					DrawLine(offset.X, offset.Y+squareSize, offset.X+squareSize, offset.Y+squareSize, rl.DarkGray)
					offset.X += squareSize
				case FULL:
					DrawRectangle(offset.X, offset.Y, squareSize, squareSize, rl.Gray)
					offset.X += squareSize
				case MOVING:
					DrawRectangle(offset.X, offset.Y, squareSize, squareSize, rl.DarkGray)
					offset.X += squareSize
				case BLOCK:
					DrawRectangle(offset.X, offset.Y, squareSize, squareSize, rl.LightGray)
					offset.X += squareSize
				case FADING:
					DrawRectangle(offset.X, offset.Y, squareSize, squareSize, fadingColor)
					offset.X += squareSize
				}
			}

			offset.X = controller
			offset.Y += squareSize
		}

		// Draw incoming piece (hardcoded)
		offset.X = 500
		offset.Y = 45

		controller = offset.X

		for j := 0; j < 4; j++ {
			for i := 0; i < 4; i++ {
				if incomingPiece[i][j] == EMPTY {
					DrawLine(offset.X, offset.Y, offset.X+squareSize, offset.Y, rl.LightGray)                       // top line
					DrawLine(offset.X, offset.Y, offset.X, offset.Y+squareSize, rl.LightGray)                       // left line
					DrawLine(offset.X+squareSize, offset.Y, offset.X+squareSize, offset.Y+squareSize, rl.LightGray) // right line
					DrawLine(offset.X, offset.Y+squareSize, offset.X+squareSize, offset.Y+squareSize, rl.LightGray) // bottom line
					offset.X += squareSize
				} else if incomingPiece[i][j] == MOVING {
					DrawRectangle(offset.X, offset.Y, squareSize, squareSize, rl.Gray)
					offset.X += squareSize
				}
			}

			offset.X = controller
			offset.Y += squareSize
		}

		DrawText("INCOMING:", offset.X, offset.Y-100, 10, rl.Gray)
		DrawText(fmt.Sprintf("LINES: %04d", score), 500, 250, 20, rl.Gray)

		if isPaused {
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
	piecePosX = (gridSizeX - 4) / 2 // Centerpiece in X axis
	piecePosY = 0                   // Start piece at top of the grid

	// If the game is starting, and you are going to create the first piece, we create an extra one
	if isFirst {
		getRandomPiece()

		isFirst = false
	}

	// We assign the incoming piece to the actual piece
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			piece[i][j] = incomingPiece[i][j]
		}
	}

	// We assign a random piece to the incoming one
	getRandomPiece()

	// Assign the piece to the grid
	for i := piecePosX; i < piecePosX+4; i++ {
		for j := 0; j < 4; j++ {
			if piece[i-piecePosX][j] == MOVING {
				grid[i][j] = MOVING
			}
		}
	}

	return true
}

// getRandomPiece Get a random piece
func getRandomPiece() {
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

// checkCollisionY Check if the current moving piece is colliding with the ground or another piece
func checkCollisionY() bool {
	for j := gridSizeY - 2; j >= 0; j-- { // We check from the bottom playable line to the top
		for i := 1; i < gridSizeX-1; i++ { // We check from left to right (playable area)
			// Check if any square of the piece is colliding with the ground wall (that is NOT playable area) (BLOCK)
			// or with another piece that is already in the grid (FULL)
			if grid[i][j] == MOVING && (grid[i][j+1] == FULL || grid[i][j+1] == BLOCK) {
				return true
			}
		}
	}

	return false
}

func CheckCompletion(lineToDelete *bool) {
	var calculator int

	for j := gridSizeY - 2; j >= 0; j-- {
		calculator = 0

		for i := 1; i < gridSizeX-1; i++ {
			// Count each square of the line
			if grid[i][j] == FULL {
				calculator++
			}

			// Check if we completed the whole line
			if calculator == gridSizeX-2 {
				*lineToDelete = true
				calculator = 0

				// Mark the completed line
				for z := 1; z < gridSizeX-1; z++ {
					grid[z][j] = FADING
				}
			}
		}
	}
}

func DeleteCompleteLines() int {
	var deletedLines int

	// Erase the completed line
	for j := gridSizeY - 2; j >= 0; j-- {
		for grid[1][j] == FADING {
			for i := 1; i < gridSizeX-1; i++ {
				grid[i][j] = EMPTY
			}

			// Erase the completed line by relocating all the current lines of the grid down
			// otherwise there will be a gap with EMPTY cells
			for j2 := j - 1; j2 >= 0; j2-- {
				for i2 := 1; i2 < gridSizeX-1; i2++ {
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
		for j := gridSizeY - 2; j >= 0; j-- {
			for i := 1; i < gridSizeX-1; i++ {
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
			for j := gridSizeY - 2; j >= 0; j-- {
				for i := 1; i < gridSizeX-1; i++ { // We check the matrix from left to right
					// Move everything to the left
					if grid[i][j] == MOVING {
						grid[i-1][j] = MOVING
						grid[i][j] = EMPTY
					}
				}
			}

			piecePosX--
		}
	} else if rl.IsKeyDown(rl.KeyRight) { // Move Right
		// Check if is possible to move to right
		for j := gridSizeY - 2; j >= 0; j-- {
			for i := 1; i < gridSizeX-1; i++ {
				if grid[i][j] == MOVING {
					// Check if we are touching the right wall, or we have a full square at the right
					if i+1 == gridSizeX-1 || grid[i+1][j] == FULL {
						collision = true
					}
				}
			}
		}

		// If able, move right
		if !collision {
			for j := gridSizeY - 2; j >= 0; j-- {
				for i := gridSizeX - 1; i >= 1; i-- { // We check the matrix from right to left
					// Move everything to the right
					if grid[i][j] == MOVING {
						grid[i+1][j] = MOVING
						grid[i][j] = EMPTY
					}
				}
			}

			piecePosX++
		}
	}

	return collision
}

func ResolveTurnMovement() bool {
	// Input for turning the piece
	if rl.IsKeyDown(rl.KeyUp) {
		var (
			aux     gridSquare
			checker bool
		)

		// Check all turning possibilities
		if (grid[piecePosX+3][piecePosY] == MOVING) &&
			(grid[piecePosX][piecePosY] != EMPTY) &&
			(grid[piecePosX][piecePosY] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+3][piecePosY+3] == MOVING) &&
			(grid[piecePosX+3][piecePosY] != EMPTY) &&
			(grid[piecePosX+3][piecePosY] != MOVING) {
			checker = true
		}

		if (grid[piecePosX][piecePosY+3] == MOVING) &&
			(grid[piecePosX+3][piecePosY+3] != EMPTY) &&
			(grid[piecePosX+3][piecePosY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePosX][piecePosY] == MOVING) &&
			(grid[piecePosX][piecePosY+3] != EMPTY) &&
			(grid[piecePosX][piecePosY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+1][piecePosY] == MOVING) &&
			(grid[piecePosX][piecePosY+2] != EMPTY) &&
			(grid[piecePosX][piecePosY+2] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+3][piecePosY+1] == MOVING) &&
			(grid[piecePosX+1][piecePosY] != EMPTY) &&
			(grid[piecePosX+1][piecePosY] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+2][piecePosY+3] == MOVING) &&
			(grid[piecePosX+3][piecePosY+1] != EMPTY) &&
			(grid[piecePosX+3][piecePosY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePosX][piecePosY+2] == MOVING) &&
			(grid[piecePosX+2][piecePosY+3] != EMPTY) &&
			(grid[piecePosX+2][piecePosY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+2][piecePosY] == MOVING) &&
			(grid[piecePosX][piecePosY+1] != EMPTY) &&
			(grid[piecePosX][piecePosY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+3][piecePosY+2] == MOVING) &&
			(grid[piecePosX+2][piecePosY] != EMPTY) &&
			(grid[piecePosX+2][piecePosY] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+1][piecePosY+3] == MOVING) &&
			(grid[piecePosX+3][piecePosY+2] != EMPTY) &&
			(grid[piecePosX+3][piecePosY+2] != MOVING) {
			checker = true
		}

		if (grid[piecePosX][piecePosY+1] == MOVING) &&
			(grid[piecePosX+1][piecePosY+3] != EMPTY) &&
			(grid[piecePosX+1][piecePosY+3] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+1][piecePosY+1] == MOVING) &&
			(grid[piecePosX+1][piecePosY+2] != EMPTY) &&
			(grid[piecePosX+1][piecePosY+2] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+2][piecePosY+1] == MOVING) &&
			(grid[piecePosX+1][piecePosY+1] != EMPTY) &&
			(grid[piecePosX+1][piecePosY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+2][piecePosY+2] == MOVING) &&
			(grid[piecePosX+2][piecePosY+1] != EMPTY) &&
			(grid[piecePosX+2][piecePosY+1] != MOVING) {
			checker = true
		}

		if (grid[piecePosX+1][piecePosY+2] == MOVING) &&
			(grid[piecePosX+2][piecePosY+2] != EMPTY) &&
			(grid[piecePosX+2][piecePosY+2] != MOVING) {
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

		for j := gridSizeY - 2; j >= 0; j-- {
			for i := 1; i < gridSizeX-1; i++ {
				if grid[i][j] == MOVING {
					grid[i][j] = EMPTY
				}
			}
		}

		for i := piecePosX; i < piecePosX+4; i++ {
			for j := piecePosY; j < piecePosY+4; j++ {
				if piece[i-piecePosX][j-piecePosY] == MOVING {
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
