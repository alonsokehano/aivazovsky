package main

type Matrix struct {
	X, Y, Z int
	Values  [][][]interface{}
}

func (m *Matrix) New(f func(a, b, c int) interface{}) *Matrix {
	m.Values = make([][][]interface{}, m.X)
	for i := 0; i < m.X; i++ {
		m.Values[i] = make([][]interface{}, m.Y)
		for j := 0; j < m.Y; j++ {
			m.Values[i][j] = make([]interface{}, m.Z)
			for k := 0; k < m.Z; k++ {
				m.Values[i][j][k] = f(i, j, k)
			}
		}
	}
	return m
}

func (m *Matrix) Each(f func(a, b, c int, value interface{})) *Matrix {
	for i := 0; i < m.X; i++ {
		for j := 0; j < m.Y; j++ {
			for k := 0; k < m.Z; k++ {
				f(i, j, k, m.Values[i][j][k])
			}
		}
	}
	return m
}
