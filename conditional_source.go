package configo

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

func (this *ConditionalSource) Add(key string, values ...interface{}) *ConditionalSource {
	this.inner.Add(key, values...)
	return this
}

func (this *ConditionalSource) Strings(key string) ([]string, error) {
	if !this.condition() {
		return nil, KeyNotFoundError
	}

	return this.inner.Strings(key)
}

func (this *ConditionalSource) Initialize() {}
