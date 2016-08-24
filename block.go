package main

import (
	"fmt"
	"math/rand"
)

type BlockConfig struct {
	synapses_sens_radius int
}

type Neuron struct {
	x, y, z         int
	value, newvalue float32
	weights         [][][]float32
}

func (n *Neuron) initialize(config BlockConfig) {
	d := config.synapses_sens_radius*2 + 1
	n.weights = make([][][]float32, d)
	for i := 0; i < d; i++ {
		n.weights[i] = make([][]float32, d)
		for j := 0; j < d; j++ {
			n.weights[i][j] = make([]float32, d)
			for k := 0; k < d; k++ {
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

func (b *Block) Initialize() {
	b.neurons = make([][][]Neuron, b.x)
	for i := 0; i < b.x; i++ {
		b.neurons[i] = make([][]Neuron, b.y)
		for j := 0; j < b.y; j++ {
			b.neurons[i][j] = make([]Neuron, b.z)
			for k := 0; k < b.z; k++ {
				b.neurons[i][j][k] = Neuron{x: i, y: j, z: k}
				b.neurons[i][j][k].initialize(b.config)
			}
		}
	}
}

func (b *Block) Vertices(vertices []float32) {
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

func (b *Block) Colors(colors []float32) {
	var index int
	for i := 0; i < b.x; i++ {
		for j := 0; j < b.y; j++ {
			for k := 0; k < b.z; k++ {
				if b.neurons[i][j][k].value >= 1 {
					colors[index] = 1.0
				} else {
					colors[index] = 0.0
				}
				colors[index+1] = 0.
				colors[index+2] = 0.
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

func (block *Block) Process() {
	var sum float32
	var posA, posB, posC int
	r := block.config.synapses_sens_radius
	d := block.config.synapses_sens_radius*2 + 1
	for i := 0; i < block.x; i++ {
		for j := 0; j < block.y; j++ {
			for k := 0; k < block.z; k++ {
				for a := 0; a < d; a++ {
					posA = i - r + a
					if posA >= 0 && posA < block.x && posA != i {
						for b := 0; b < d; b++ {
							posB = j - r + b
							if posB >= 0 && posB < block.y && posB != j {
								for c := 0; c < d; c++ {
									posC = k - r + c
									if posC >= 0 && posC < block.z && posC != k {
										sum += block.neurons[i][j][k].weights[a][b][c] * block.neurons[posA][posB][posC].value
									}
								}
							}
						}
					}
				}
				block.neurons[i][j][k].newvalue = sum
			}
		}
	}

	for i := 0; i < block.x; i++ {
		for j := 0; j < block.y; j++ {
			for k := 0; k < block.z; k++ {
				block.neurons[i][j][k].value = block.neurons[i][j][k].newvalue
			}
		}
	}
}

func (b *Block) Run(c chan int) {
	for i := 0; i < 100000; i++ {
		fmt.Println("step", i)
		b.Process()
		c <- i
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
