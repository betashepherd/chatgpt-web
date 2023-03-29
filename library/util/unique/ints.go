package unique

func Ints(input []int) []int {
	u := []int{}
	m := map[int]bool{}

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func Uints(input []uint) []uint {
	u := []uint{}
	m := map[uint]bool{}

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

