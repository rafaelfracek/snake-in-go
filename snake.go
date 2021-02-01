/**
 * Copyright (c) 2021
 *
 * @author      Mentisimo Rafael Frącek
 * @license     GNU General Public License version 3 (GPLv3)
 */

package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gonutz/w32"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

const (
	WINDOW_TITLE     = "Go Snake!"
	ARROW_LEFT       = "ArrowLeft"
	ARROW_RIGHT      = "ArrowRight"
	ARROW_UP         = "ArrowUp"
	ARROW_DOWN       = "ArrowDown"
	ARROW_PREFIX     = "Arrow"
	BACKGROUND_COLOR = "#000"
	SNAKE_COLOR      = "#f00"
	APPLE_COLOR      = "#0f0"
)

type Position struct {
	x int16
	y int16
}

type GameState struct {
	snake    []Position
	move     Position
	nextMove Position
	points   int16
	apple    *Position
	running  bool
}

func main() {
	hideConsole()
	rand.Seed(time.Now().UnixNano())

	wnd, cv, err := sdlcanvas.CreateWindow(400, 400, WINDOW_TITLE)
	if err != nil {
		panic(err)
	}
	wnd.Window.SetResizable(false)
	defer wnd.Destroy()

	gameState := newDefaultGameState()

	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				makeStep(gameState)
			}
		}
	}()

	wnd.KeyDown = createKeyDownFunc(gameState)

	wnd.MainLoop(createRenderFrameFunc(cv, gameState))
}

func makeStep(state *GameState) {
	if state.running {
		if state.nextMove.x != -state.move.x || state.nextMove.y != -state.move.y {
			state.move = state.nextMove
		}
		state.snake = append(state.snake, movePosition(getHead(state), &state.move))
		collision := len(*filterPositions(&state.snake, func(position *Position) bool {
			return position.x == getHead(state).x && position.y == getHead(state).y
		})) >= 2
		if collision {
			*state = *newDefaultGameState()
		} else {
			if getHead(state).x == state.apple.x && getHead(state).y == state.apple.y {
				state.points++
				*state.apple = *generateApplePosition(state)
			}
			if state.points <= 0 {
				state.snake = state.snake[1:]
			} else {
				state.points--
			}
		}
	}
}

func getHead(gameState *GameState) *Position {
	return &gameState.snake[len(gameState.snake)-1]
}

func processBound(number int16) int16 {
	if number > 19 {
		return 0
	} else if number < 0 {
		return 19
	}
	return number
}

func filterPositions(positions *[]Position, test func(*Position) bool) *[]Position {
	result := []Position{}
	for _, position := range *positions {
		if test(&position) {
			result = append(result, position)
		}
	}
	return &result
}

func generateApplePosition(state *GameState) *Position {
	var applePosition Position
	for {
		applePosition = newPosition(generateRandomNumber(19), generateRandomNumber(19))
		if len(*filterPositions(&state.snake, isTheSamePosition(&applePosition))) == 0 {
			break
		}
	}
	return &applePosition
}

func generateRandomNumber(max int16) int16 {
	return int16(rand.Intn(int(max + 1)))
}

func isTheSamePosition(position *Position) func(position *Position) bool {
	return func(targetPosition *Position) bool {
		return position.x == targetPosition.x && position.y == targetPosition.y
	}
}

func createRenderFrameFunc(cv *canvas.Canvas, state *GameState) func() {
	return func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle(BACKGROUND_COLOR)
		cv.FillRect(0, 0, float64(w), float64(h))
		cv.SetFillStyle(SNAKE_COLOR)
		for _, square := range state.snake {
			cv.FillRect(float64(square.x*20), float64(square.y*20), float64(18), float64(18))
		}
		cv.SetFillStyle(APPLE_COLOR)
		cv.FillRect(float64(state.apple.x*20), float64(state.apple.y*20), float64(18), float64(18))
	}
}

func createKeyDownFunc(state *GameState) func(int, rune, string) {
	return func(_ int, _ rune, name string) {
		switch name {
		case ARROW_LEFT:
			state.nextMove = newPosition(-1, 0)
		case ARROW_RIGHT:
			state.nextMove = newPosition(1, 0)
		case ARROW_UP:
			state.nextMove = newPosition(0, -1)
		case ARROW_DOWN:
			state.nextMove = newPosition(0, 1)
		}
		if strings.HasPrefix(name, ARROW_PREFIX) {
			state.running = true
		}
	}
}

func newPosition(x int16, y int16) Position {
	return Position{
		x: x,
		y: y,
	}
}

func movePosition(position *Position, movePosition *Position) Position {
	return Position{
		x: processBound(position.x + movePosition.x),
		y: processBound(position.y + movePosition.y),
	}
}

func newDefaultGameState() *GameState {
	state := GameState{
		snake:    []Position{newPosition(10, 10)},
		move:     newPosition(0, 0),
		nextMove: newPosition(0, 0),
		points:   2,
		running:  false,
	}
	state.apple = generateApplePosition(&state)
	return &state
}

func hideConsole() {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcessId := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcessId {
			w32.ShowWindowAsync(console, w32.SW_HIDE)
		}
	}
}
