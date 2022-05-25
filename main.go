package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	x    mark   = "x"
	o    mark   = "o"
	winX string = string(x) + string(x) + string(x)
	winO string = string(o) + string(o) + string(o)
)

var (
	ErrInvalidCellNumber = errors.New(`invalid cell number: must be in range 1 - 9`)
	ErrInvalidUserMark   = errors.New(`invalid mark, only "x" / "o" is allowed`)
)

var (
	// to win a game certain indexes should be matched with below ones
	winComboRow1      = []int{0, 1, 2}
	winComboRow2      = []int{3, 4, 5}
	winComboRow3      = []int{6, 7, 8}
	winComboCol1      = []int{0, 3, 6}
	winComboCol2      = []int{1, 4, 7}
	winComboCol3      = []int{2, 5, 8}
	winComboDiagonal1 = []int{0, 4, 8}
	winComboDiagonal2 = []int{2, 4, 6}

	winCombos = [][]int{winComboRow1, winComboRow2, winComboRow3, winComboCol1, winComboCol2, winComboCol3, winComboDiagonal1, winComboDiagonal2}
)

func init() {
	fmt.Println("starting tic-tac-toe game")
}

func main() {
	// readTest()

	fmt.Println(`select your mark: "x", "o"`)

	userMark := readUserMark()

	fmt.Println("gamerules: type a number from 1 to 9 to place selected mark")

	runGameCycle(userMark)
}

func runGameCycle(userMark mark) {
	// todo: change error processing?

	grid := newGrid()
	grid.draw()

	for {
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			log.Fatal("[error] reading player's number from stdin: ", err)
		}

		switch char {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			cellNum, err := strconv.Atoi(string(char))
			if err != nil {
				fmt.Println("[error] " + ErrInvalidCellNumber.Error())
			}

			if err = grid.addMark(cellNum, userMark); err != nil {
				fmt.Println("[error] " + err.Error())
			}

			grid.draw()

			if isWinCombo(grid.data, userMark) {
				fmt.Printf("mark [%s] is a winner! (rewrite me plz)\n", userMark)
				fmt.Println("thanks for playing!")
				os.Exit(0)
			}
		case 'q':
			exitApp()
		default:
			fmt.Println("[error] invalid char. allowed 1 - 9")
			continue
		}
	}
}

type mark string

func (m mark) String() string {
	return string(m)
}

type marksData [9]mark

var emptyMarksData = marksData{"-", "-", "-", "-", "-", "-", "-", "-", "-"}

func newGrid() *grid {
	return &grid{data: emptyMarksData}
}

type grid struct {
	data marksData
}

func (g grid) draw() {
	b := strings.Builder{}
	defer b.Reset()

	for i, mark := range g.data {
		if (i+3)%3 == 0 {
			b.WriteString("\n")
		}

		b.WriteString(mark.String())
	}

	fmt.Println(b.String())
}

func (g *grid) addMark(cellNum int, m mark) error {
	if cellNum < 1 && cellNum > 9 {
		return ErrInvalidCellNumber
	}

	g.data[cellNum-1] = m

	return nil
}

func readUserMark() mark {
	for {
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("[error] reading rune from stdin: ", err)
			exitApp()
		}

		switch char {
		case 'x':
			return x
		case 'X':
			return x
		case 'o':
			return o
		case 'O':
			return o
		case '0':
			return o
		case 'q':
			exitApp()
		}

		fmt.Println("[error] ", ErrInvalidUserMark)
	}
}

func isWinCombo(in marksData, m mark) bool {
	if len(in) != 9 {
		panic("input must consist of 9 marks")
	}

	for _, winIndexes := range winCombos {
		var temp string

		for _, winIndex := range winIndexes {
			temp += string(in[winIndex])
		}

		if len(temp) != 3 {
			continue
		}

		if m == x && temp == winX {
			return true
		}

		if m == o && temp == winO {
			return true
		}
	}

	return false
}

// comboMatcher/comboChecker

func linuxClear() {
	// clear["linux"] = func() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
	// }
}

func exitApp() {
	// some clear/app close func?
	fmt.Println("[exiting]")
	os.Exit(0)
}

// ======================== TEMP ========================

func cliTest() {
	app := &cli.App{
		Name:  "greet",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {

		log.Fatal(err)
	}
}
