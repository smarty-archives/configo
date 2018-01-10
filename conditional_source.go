package configo

// ConditionalSource resolves values based on a condition supplied as a callback.
type ConditionalSource struct {
	condition func() bool
	inner     *DefaultSource
}

// NewConditionalSource creates a conditional source with the provided condition callback
// and key/value pairs.
func NewConditionalSource(condition func() bool, pairs ...DefaultPair) *ConditionalSource {
	return &ConditionalSource{
		condition: condition,
		inner:     NewDefaultSource(pairs...),
	}
}

// Strings returns the value of the corresponding key, or KeyNotFoundError if the condition is false.
func (this *ConditionalSource) Strings(key string) ([]string, error) {
	if !this.condition() {
		return nil, KeyNotFoundError
	}

	return this.inner.Strings(key)
}

func (this *ConditionalSource) Initialize() {}
