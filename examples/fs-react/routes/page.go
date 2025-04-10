package index

type Page struct{}

func (Page) Load(c any) (map[string][]string, error) {
	return map[string][]string{"names": {"one", "two", "three"}}, nil
}
