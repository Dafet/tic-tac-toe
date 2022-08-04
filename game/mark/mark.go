package mark

type Mark string

func (m Mark) Str() string {
	return string(m)
}

const (
	X Mark = "x"
	O Mark = "o"

	None Mark = "-"
)
