package newton

type Source interface {
	Name() string
	Strings() ([]string, error)
}
