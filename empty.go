package slot

// Empty slots
func Empty() Slots {
	return emptySlots{}
}

type emptySlots struct {
}

var _ Slots = (emptySlots{})

func (emptySlots) ReadString() (string, error) {
	return "", nil
}

func (emptySlots) Slot(slot string) NamedSlot {
	return &emptySlot{
		data: "",
	}
}

type emptySlot struct {
	data string
}

var _ NamedSlot = (*emptySlot)(nil)

func (e *emptySlot) ReadString() (string, error) {
	return e.data, nil
}

func (e *emptySlot) WriteString(data string) error {
	e.data += data
	return nil
}
