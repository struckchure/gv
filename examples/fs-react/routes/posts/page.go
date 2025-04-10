package posts

type Page struct{}

func (Page) Load() (any, error) {
	return nil, nil
}
