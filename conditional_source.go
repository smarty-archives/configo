package newton

type ConditionalSource struct {
	condition func() bool
	inner     *DefaultSource
}

func NewConditionalSource(condition func() bool) *ConditionalSource {
	return &ConditionalSource{
		condition: condition,
		inner:     NewDefaultSource(),
	}
}

func (this *ConditionalSource) Add(key string, values ...interface{}) {
	this.inner.Add(key, values...)
}

func (this *ConditionalSource) Strings(key string) ([]string, error) {
	if !this.condition() {
		return nil, KeyNotFoundError
	}

	return this.inner.Strings(key)
}
