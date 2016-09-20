package main

type Matrix struct {
	x, y, z int
	values  [][][]interface{}
}

func (m *Matrix) New(f func(a, b, c int) interface{}) *Matrix {
	m.values = make([][][]interface{}, m.x)
	for i := 0; i < m.x; i++ {
		m.values[i] = make([][]interface{}, m.y)
		for j := 0; j < m.y; j++ {
			m.values[i][j] = make([]interface{}, m.z)
			for k := 0; k < m.z; k++ {
				m.values[i][j][k] = f(i, j, k)
			}
		}
	}
	return m
}

func (m *Matrix) Each(f func(a, b, c int, value interface{})) *Matrix {
	for i := 0; i < m.x; i++ {
		for j := 0; j < m.y; j++ {
			for k := 0; k < m.z; k++ {
				f(i, j, k, m.values[i][j][k])
			}
		}
	}
	return m
}
