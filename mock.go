package slot

// Mock slots
type Mock struct {
	Data  string            // Slot data
	Named map[string]string // Named slots
}

var _ Slots = (*Mock)(nil)

func (m *Mock) ReadString() (string, error) {
	return m.Data, nil
}

func (m *Mock) Slot(slot string) NamedSlot {
	return &mockNamed{m.Named[slot]}
}

type mockNamed struct {
	data string
}

var _ NamedSlot = (*mockNamed)(nil)

func (m *mockNamed) ReadString() (string, error) {
	return m.data, nil
}

func (m *mockNamed) WriteString(data string) error {
	m.data += data
	return nil
}
