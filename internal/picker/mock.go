package picker

// MockPicker is a test double that returns a pre-configured selection.
// It is exported so cmd tests (and any other package) can inject it.
type MockPicker struct {
	Selection string
	Err       error
}

// Pick implements Picker.
func (m *MockPicker) Pick(_ string, _ []string) (string, error) {
	return m.Selection, m.Err
}
