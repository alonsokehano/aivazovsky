package main

type Neuron struct {
	x, y, z int
	value   float32
}

type Block struct {
	x, y, z int
	neurons [][][]Neuron
}

func (b Block) NewBlock(x, y, z int) Block {
	neurons := make([][][]Neuron, x)
	for i := 0; i < x; i++ {
		neurons[i] = make([][]Neuron, y)
		for j := 0; j < y; j++ {
			neurons[i][j] = make([]Neuron, z)
			for k := 0; k < z; k++ {
				neurons[i][j][k] = Neuron{x: i, y: j, z: k}
			}
		}
	}
	return Block{x: x, y: y, z: z, neurons: neurons}
}
