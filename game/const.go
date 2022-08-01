package game

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

	winX = "xxx"
	winO = "ooo"

	// emptyTurnResult = TurnResult{}
)
