package transitions

type Transition struct {
	From        string
	To          string
	Destination string

	enter *Enter
}
