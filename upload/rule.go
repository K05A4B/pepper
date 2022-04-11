package upload

type Rule struct {
	MaxSize   int64
	MinSize   int64
	MaxNumber int
	Mime      Mime
}
