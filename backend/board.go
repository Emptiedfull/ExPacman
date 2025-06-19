package main

import (
	"context"
	"fmt"
	"sync"
)

var boardString = []string{
	"############################",
	"#............##............#",
	"#.####.#####.##.#####.####.#",
	"#.####.#####.##.#####.####.#",
	"#.####.#####.##.#####.####.#",
	"#..........................#",
	"#.####.##.########.##.####.#",
	"#.####.##.########.##.####.#",
	"#......##....##....##......#",
	"######.##### ## #####.######",
	"######.##### ## #####.######",
	"######.##          ##.######",
	"######.## ###  ### ##.######",
	"######.## #  b   # ##.######",
	" p     ##   c  d   ##       ",
	"######.## #   a  # ##.######",
	"######.## ###  ### ##.######",
	"######.##          ##.######",
	"######.## ######## ##.######",
	"######.## ######## ##.######",
	"#............##............#",
	"#.####.#####.##.#####.####.#",
	"#.####.#####.##.#####.####.#",
	"#o..##................##..o#",
	"###.##.##.########.##.##.###",
	"###.##.##.########.##.##.###",
	"#......##....##....##......#",
	"#.##########.##.##########.#",
	"#.##########.##.##########.#",
	"#..........................#",
	"############################",
}

type Enemy int

const (
	NoPlayer Enemy = iota
	enemyA
	enemyB
	enemyC
	enemyD
)

type background int

const (
	empty background = iota
	wall
	food
)

type Cell struct {
	background background
	pacman     bool
	enemy      Enemy
}

type Board struct {
	Width  int
	Height int
	Cells  [][]Cell
}

func newBoard(width, height int) *Board {
	board := &Board{
		Width:  width,
		Height: height,
		Cells:  make([][]Cell, height),
	}
	for i := range board.Cells {
		board.Cells[i] = make([]Cell, width)
		for j := range board.Cells[i] {
			board.Cells[i][j] = Cell{
				background: empty,
				pacman:     false,
				enemy:      NoPlayer,
			}

		}
	}

	return board

}

func (b Board) visualize() string {
	var Output string
	for _, row := range b.Cells {
		for _, cell := range row {
			var temp string
			switch cell.background {
			case wall:
				temp = "W "
			case food:
				temp = ". "
			case empty:
				temp = "  "
			}
			switch cell.enemy {
			case enemyA:
				temp = "a "
			case enemyB:
				temp = "b "
			case enemyC:
				temp = "c "
			case enemyD:
				temp = "d "
			}
			if cell.pacman {
				temp = "P "
			}
			Output += temp

		}
		Output += "\n"
	}
	return Output
}

func ParseBoardString(boardString []string) (*Board, [5][2]int) {
	board := newBoard(len(boardString[0]), len(boardString))
	var Positions [5][2]int
	for i, row := range boardString {
		for j, char := range row {
			switch char {
			case '#':
				board.Cells[i][j].background = wall
			case '.':
				board.Cells[i][j].background = food
			case 'p':
				Positions[0] = [2]int{i, j}
				board.Cells[i][j].pacman = true
			case 'a':
				Positions[1] = [2]int{i, j}
				board.Cells[i][j].enemy = enemyA
			case 'b':
				Positions[2] = [2]int{i, j}
				board.Cells[i][j].enemy = enemyB
			case 'c':
				Positions[3] = [2]int{i, j}
				board.Cells[i][j].enemy = enemyC
			case 'd':
				Positions[4] = [2]int{i, j}
				board.Cells[i][j].enemy = enemyD
			}
		}
	}
	return board, Positions

}

type move int

const (
	Up move = iota
	Down
	Left
	Right
	None
)

type GameState struct {
	Board           *Board
	MoveState       [5]move
	PlayerPositions [5][2]int
	Scores          [5]int
	mut             sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
}

func (g *GameState) move() {
	for i, m := range g.MoveState {
		if i == 0 {
			//pacman
			x, y := g.PlayerPositions[i][0], g.PlayerPositions[i][1]
			switch m {
			case Left:
				if y-1 < 0 {
					fmt.Println("wrap")
					g.PlayerPositions[i][1] = g.Board.Width - 1
					g.Board.Cells[x][g.Board.Width-1].pacman = true
					g.Board.Cells[x][y].pacman = false
					continue
				}
				if g.Board.Cells[x][y-1].background == wall {
					continue
				}

				g.PlayerPositions[i][1] = y - 1
				g.Board.Cells[x][y].pacman = false
				g.Board.Cells[x][y-1].pacman = true
			case Right:
				if y+1 >= g.Board.Width {
					g.PlayerPositions[i][1] = 0
					g.Board.Cells[x][0].pacman = true
					g.Board.Cells[x][y].pacman = false
					continue
				}
				if g.Board.Cells[x][y+1].background == wall {
					continue
				}
				g.Board.Cells[x][y].pacman = false

				g.PlayerPositions[i][1] = y + 1
				g.Board.Cells[x][y+1].pacman = true
			case Down:
				if x+1 >= g.Board.Height {
					g.PlayerPositions[i][0] = 0
					g.Board.Cells[0][y].pacman = true
					g.Board.Cells[x][y].pacman = false
					continue
				}
				if g.Board.Cells[x+1][y].background == wall {
					continue
				}
				g.Board.Cells[x][y].pacman = false

				g.PlayerPositions[i][0] = x + 1
				g.Board.Cells[x+1][y].pacman = true
			case Up:
				if x-1 < 0 {
					g.PlayerPositions[i][0] = g.Board.Height - 1
					g.Board.Cells[g.Board.Height-1][y].pacman = true
					g.Board.Cells[x][y].pacman = false
					continue
				}
				if g.Board.Cells[x-1][y].background == wall {
					continue
				}
				g.Board.Cells[x][y].pacman = false

				g.PlayerPositions[i][0] = x - 1
				g.Board.Cells[x-1][y].pacman = true

			}

		} else {
			//enemies

			x, y := g.PlayerPositions[i][0], g.PlayerPositions[i][1]
			switch m {
			case Left:
				if y-1 < 0 {
					fmt.Println("wrap")
					g.PlayerPositions[i][1] = g.Board.Width - 1
					g.Board.Cells[x][g.Board.Width-1].enemy = Enemy(i)

					g.Board.Cells[x][y].enemy = NoPlayer
					continue
				}
				if g.Board.Cells[x][y-1].background == wall || g.Board.Cells[x][y-1].enemy != NoPlayer {
					continue
				}

				g.PlayerPositions[i][1] = y - 1
				g.Board.Cells[x][y].enemy = NoPlayer
				g.Board.Cells[x][y-1].enemy = Enemy(i)

			case Right:
				if y+1 >= g.Board.Width {
					g.PlayerPositions[i][1] = 0
					g.Board.Cells[x][0].enemy = Enemy(i)
					g.Board.Cells[x][y].enemy = NoPlayer

					continue
				}
				if g.Board.Cells[x][y+1].background == wall || g.Board.Cells[x][y-1].enemy != NoPlayer {
					continue
				}
				g.Board.Cells[x][y].enemy = NoPlayer

				g.PlayerPositions[i][1] = y + 1
				g.Board.Cells[x][y+1].enemy = Enemy(i)

			case Down:

				if x+1 >= g.Board.Height {
					g.PlayerPositions[i][0] = 0
					g.Board.Cells[0][y].enemy = Enemy(i)
					g.Board.Cells[x][y].enemy = NoPlayer
					continue
				}
				if g.Board.Cells[x+1][y].background == wall || g.Board.Cells[x][y-1].enemy != NoPlayer {
					continue
				}
				g.Board.Cells[x][y].enemy = NoPlayer

				g.PlayerPositions[i][0] = x + 1
				g.Board.Cells[x+1][y].enemy = Enemy(i)
			case Up:
				if x-1 < 0 {
					g.PlayerPositions[i][0] = g.Board.Height - 1
					g.Board.Cells[g.Board.Height-1][y].enemy = Enemy(i)

					continue
				}
				if g.Board.Cells[x-1][y].background == wall || g.Board.Cells[x][y-1].enemy != NoPlayer {
					continue
				}
				g.Board.Cells[x][y].enemy = NoPlayer

				g.PlayerPositions[i][0] = x - 1
				g.Board.Cells[x-1][y].enemy = Enemy(i)

			}
		}
	}
}

func InitializeGameState() *GameState {
	board, positions := ParseBoardString(boardString)
	ctx, cancel := context.WithCancel(context.Background())
	return &GameState{
		Board:           board,
		MoveState:       [5]move{None, None, None, None, None},
		Scores:          [5]int{0, 200, 200, 200, 200},
		PlayerPositions: positions,
		mut:             sync.Mutex{},
		ctx:             ctx,
		cancel:          cancel,
	}

}

func (g *GameState) gametick() (over bool) {
	g.move()
	if g.Board.Cells[g.PlayerPositions[0][0]][g.PlayerPositions[0][1]].enemy != NoPlayer {
		fmt.Println("Game over!")
		return true
	}
	if g.Board.Cells[g.PlayerPositions[0][0]][g.PlayerPositions[0][1]].background == food {
		g.Board.Cells[g.PlayerPositions[0][0]][g.PlayerPositions[0][1]].background = empty
		g.Scores[0]++
	}
	for i, score := range g.Scores[1:] {
		if score > 0 {
			g.Scores[i+1] = score - 1
		}
	}

	return false

}
