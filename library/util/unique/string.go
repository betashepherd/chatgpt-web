package unique

func Strings(input []string) []string {
	u := []string{}
	m := map[string]bool{}

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
