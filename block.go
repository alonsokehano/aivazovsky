package main

import (
	"math/rand"
)

type BlockConfig struct {
	/* Length of synapses */
	synapses_sens_radius int

	/* Synaps activity when neuron became active (spiking) */
	synapses_threshold float32

	/* Speed of decreasing of internal neuron value while spiking */
	spiking_speed float32

	/* Speed of decreasing of internal neuron value while relaxing */
	relaxation_speed float32

	/* Condition (internal neuron value) when relaxation should ends */
	relaxation_threshold float32
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
				b.neurons[i][j][k].Initialize(b.config)
			}
		}
	}
}

/*
 Vertex rendering
*/
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

/*
 Colors rendering
*/
func (b *Block) Colors(colors []float32) {
	var index int
	for i := 0; i < b.x; i++ {
		for j := 0; j < b.y; j++ {
			for k := 0; k < b.z; k++ {
				if b.neurons[i][j][k].IsActive() {
					colors[index] = 1.0
					colors[index+1] = 0.0
					colors[index+2] = 0.0
				} else if (b.neurons[i][j][k]).IsRelaxing() {
					colors[index] = 0.0
					colors[index+1] = 0.0
					colors[index+2] = 1.0
				} else {
					colors[index] = 0.0
					colors[index+1] = 0.0
					colors[index+2] = 0.0
				}
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
					b.neurons[i][j][k].SetValue(b.config.synapses_threshold, b.config)
				}
			}
		}
	}
}

func (block *Block) Process() {
	var sum float32
	var neuron Neuron
	var posA, posB, posC int
	r := block.config.synapses_sens_radius
	d := block.config.synapses_sens_radius*2 + 1

	/* Run through all neurons and calculate synapses activity */
	for i := 0; i < block.x; i++ {
		for j := 0; j < block.y; j++ {
			for k := 0; k < block.z; k++ {
				neuron = block.neurons[i][j][k]

				if neuron.IsIdle() {
					/*
					 If neuron is in 'idle' state, then calculate
					 synaps activity and update his new value
					*/
					for a := 0; a < d; a++ {
						posA = i - r + a
						if posA >= 0 && posA < block.x && posA != i {
							for b := 0; b < d; b++ {
								posB = j - r + b
								if posB >= 0 && posB < block.y && posB != j {
									for c := 0; c < d; c++ {
										posC = k - r + c
										if posC >= 0 && posC < block.z && posC != k {
											if block.neurons[posA][posB][posC].IsActive() {
												sum += neuron.weights[a][b][c]
											}
										}
									}
								}
							}
						}
					}
					neuron.newvalue = sum
					sum = 0
				} else if neuron.IsActive() {
					/*
						In case if neuron is already in 'active' state
						just decrement his new value
					*/
					neuron.newvalue -= block.config.spiking_speed
				} else if neuron.IsRelaxing() {
					/*
						In case if neuron is in 'relaxing' state
						just decrement his new value according to relaxation speed
					*/
					neuron.newvalue -= block.config.relaxation_speed
				}
			}
		}
	}

	for i := 0; i < block.x; i++ {
		for j := 0; j < block.y; j++ {
			for k := 0; k < block.z; k++ {
				block.neurons[i][j][k].SetValue(block.neurons[i][j][k].newvalue, block.config)
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
