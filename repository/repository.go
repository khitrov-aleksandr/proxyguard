package repository

type Repository interface {
	Save(str string, value any, exp int) error
	Get(str string) any
	Incr(str string) int64
	Expr(str string, exp int) bool
}
