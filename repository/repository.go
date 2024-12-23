package repository

type Repository interface {
	Save(str string, value any, exp int) error
	Get(str string) any
}
