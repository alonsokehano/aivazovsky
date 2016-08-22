package main

import "math/rand"

type BlockConfig struct {
	synapses_sens_radius int
}

type Neuron struct {
	x, y, z int
	value   float32
	weights [][][]float32
}

func (n *Neuron) initialize(config BlockConfig) {
	r := config.synapses_sens_radius
	n.weights = make([][][]float32, r)
	for i := 0; i < r; i++ {
		n.weights[i] = make([][]float32, r)
		for j := 0; j < r; j++ {
			n.weights[i][j] = make([]float32, r)
			for k := 0; k < r; k++ {
				n.weights[i][j][k] = rand.Float32()
			}
		}
	}
}

type Block struct {
	x, y, z int
	neurons [][][]Neuron
	config  BlockConfig
}

func (b Block) NewBlock(x, y, z int) Block {
	neurons := make([][][]Neuron, x)
	for i := 0; i < x; i++ {
		neurons[i] = make([][]Neuron, y)
		for j := 0; j < y; j++ {
			neurons[i][j] = make([]Neuron, z)
			for k := 0; k < z; k++ {
				neurons[i][j][k] = Neuron{x: i, y: j, z: k}
				neurons[i][j][k].initialize(b.config)
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
