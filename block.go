package main

import (
	"math/rand"
)

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

func (b *Block) Render(vertices []float32) {
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

func (b *Block) CreatePattern(x, y, z, r int, probability float32) {
	for i := maxInt(0, x-r); i < minInt(b.x, x+r); i++ {
		for j := maxInt(0, y-r); j < minInt(b.y, y+r); j++ {
			for k := maxInt(0, z-r); k < minInt(b.z, z+r); k++ {
				if rand.Float32() <= probability {
					b.neurons[i][j][k].value = 1.0
				}
			}
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
