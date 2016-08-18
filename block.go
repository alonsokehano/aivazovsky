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

func (b Block) Render(vertices []float32) {
	var index int
	for i := 0; i < b.x; i++ {
		for j := 0; j < b.y; j++ {
			for k := 0; k < b.z; k++ {
				vertices[index] = float32(b.neurons[i][j][k].x)
				vertices[index+1] = float32(b.neurons[i][j][k].y)
				vertices[index+2] = float32(b.neurons[i][j][k].z)
				index += 3
			}
		}
	}
}
