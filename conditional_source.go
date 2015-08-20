package configo

// ConditionalSource resolves values based on a condition supplied as a callback.
type ConditionalSource struct {
	condition func() bool
	inner     *DefaultSource
}

// NewConditionalSource creates a conditional source with the provided condition callback.
func NewConditionalSource(condition func() bool) *ConditionalSource {
	return &ConditionalSource{
		condition: condition,
		inner:     NewDefaultSource(),
	}
}

// Add registers a key/values pairing for retrieval, so long as the supplied condition is true.
func (this *ConditionalSource) Add(key string, values ...interface{}) *ConditionalSource {
	this.inner.Add(key, values...)
	return this
}

// Strings returns the value of the corresponding key, or KeyNotFoundError if the condition is false.
func (this *ConditionalSource) Strings(key string) ([]string, error) {
	if !this.condition() {
		return nil, KeyNotFoundError
	}

	return this.inner.Strings(key)
}

func (this *ConditionalSource) Initialize() {}
