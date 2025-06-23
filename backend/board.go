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
	"######.## #      # ##.######",
	" p     ##          ##       ",
	"######.## #      # ##.######",
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

type EnemyLocations struct {
	EnemyA [2]int
	EnemyB [2]int
	EnemyC [2]int
	EnemyD [2]int
}

var enemyLoc = EnemyLocations{
	EnemyA: [2]int{15, 14},
	EnemyB: [2]int{14, 12},
	EnemyC: [2]int{14, 15},
	EnemyD: [2]int{12, 12},
}

var EnemyLocArr = [4][2]int{{15, 14}, {14, 12}, {14, 15}, {12, 12}}

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

func (b Board) visualize() []string {
	var Output []string
	for _, row := range b.Cells {
		var OutputRow string
		for _, cell := range row {
			var temp string
			switch cell.background {
			case wall:
				temp = "#"
			case food:
				temp = "."
			case empty:
				temp = " "
			}
			switch cell.enemy {
			case enemyA:
				temp = "a"
			case enemyB:
				temp = "b"
			case enemyC:
				temp = "c"
			case enemyD:
				temp = "d"
			}
			if cell.pacman {
				temp = "P"
			}
			OutputRow += temp

		}
		Output = append(Output, OutputRow)
	}
	return Output
}

func ParseBoardString(boardString []string, Users int) (*Board, [][2]int) {
	board := newBoard(len(boardString[0]), len(boardString))
	Positions := make([][2]int, Users)
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
				// case 'a':
				// 	Positions[1] = [2]int{i, j}
				// 	board.Cells[i][j].enemy = enemyA
				// case 'b':
				// 	Positions[2] = [2]int{i, j}
				// 	board.Cells[i][j].enemy = enemyB
				// case 'c':
				// 	Positions[3] = [2]int{i, j}
				// 	board.Cells[i][j].enemy = enemyC
				// case 'd':
				// 	Positions[4] = [2]int{i, j}
				// 	board.Cells[i][j].enemy = enemyD
			}
		}
	}

	for i := 1; i < Users; i++ {
		enemy := Enemy(i)
		Positions[i] = EnemyLocArr[i-1]
		fmt.Println("Enemy", i, "at position", Positions[i])
		board.Cells[Positions[i][0]][Positions[i][1]].enemy = enemy
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
	MoveState       []move
	PlayerPositions [][2]int
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

func InitializeGameState(Users int) *GameState {
	board, positions := ParseBoardString(boardString, Users)
	fmt.Println(positions)
	ctx, cancel := context.WithCancel(context.Background())
	return &GameState{
		Board:           board,
		MoveState:       make([]move, Users),
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
