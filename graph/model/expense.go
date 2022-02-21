package model

type Expense struct {
	Id          int
	Date        string
	Description string
	Amount      float64
	Categories  []Category
	Comment     string
}
