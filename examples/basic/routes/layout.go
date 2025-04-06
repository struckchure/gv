package index

type Layout struct{}

func (Layout) Before(c any) error {
	return nil
}

func (Layout) After(c any) error {
	return nil
}
