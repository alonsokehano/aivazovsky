package main

import (
	"fmt"
	"math"
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
	b.New(func(i, j, k int) Neuron {
		return Neuron{}
	})
}

/*
 Vertex rendering
*/
func (b *Block) Vertices(vertices []float32) {
	var index int
	for i := 0; i < b.x; i++ {
		for j := 0; j < b.y; j++ {
			for k := 0; k < b.z; k++ {
				vertices[index] = float32(i)
				vertices[index+1] = float32(j)
				vertices[index+2] = float32(k)
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
				if b.neurons[i][j][k].isActive() {
					colors[index] = 1.0
					colors[index+1] = 0.0
					colors[index+2] = 0.0
				} else if (b.neurons[i][j][k]).isRelaxing() {
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
					b.neurons[i][j][k].Status = 1
				}
			}
		}
	}
}

func (block *Block) Process() {
	var neuron *Neuron
	var posA, posB, posC int
	var dx, dy, sigma float64
	r := block.config.synapses_sens_radius
	d := block.config.synapses_sens_radius*2 + 1

	var s float32

	/* Create pattern of activity */
	p := make([][][]float32, d)
	for i := 0; i < d; i++ {
		p[i] = make([][]float32, d)
		for j := 0; j < d; j++ {
			p[i][j] = make([]float32, d)
			for k := 0; k < d; k++ {
				dx = math.Pow(float64(i-r), 2.)
				dy = math.Pow(float64(j-r), 2.)
				sigma = math.Pow(5., 2.)
				p[i][j][k] = float32((1 / (2 * math.Pi * sigma)) * math.Exp(-1/(2*sigma)*(dx+dy)))
			}
			s += p[i][j][0]
		}
	}

	fmt.Println(s)

	/* Run through all neurons and calculate synapses activity */
	for i := 0; i < block.x; i++ {
		for j := 0; j < block.y; j++ {
			for k := 0; k < block.z; k++ {
				neuron = &block.neurons[i][j][k]
				neuron.Activity = 0
				neuron.Relax = 0

				if neuron.isIdle() {
					/*
					 If neuron is in 'idle' state, then calculate
					 synaps activity and update his new value
					*/
					for a := 0; a < d; a++ {
						posA = i - r + a
						if posA >= 0 && posA < block.x {
							for b := 0; b < d; b++ {
								posB = j - r + b
								if posB >= 0 && posB < block.y {
									for c := 0; c < d; c++ {
										posC = k - r + c
										if posC >= 0 && posC < block.z {
											if block.neurons[posA][posB][posC].isActive() {
												neuron.Activity += p[a][b][c]
											} else if block.neurons[posA][posB][posC].isRelaxing() {
												neuron.Relax += p[a][b][c]
											}
										}
									}
								}
							}
						}
					}
				} else if neuron.isActive() {
					/*
						In case if neuron is already in 'active' state
						just decrement his new value
					*/
					// neuron.newvalue = neuron.value - block.config.spiking_speed
				} else if neuron.isRelaxing() {
					/*
						In case if neuron is in 'relaxing' state
						just decrement his new value according to relaxation speed
					*/
					// neuron.newvalue = neuron.value - block.config.relaxation_speed
				}
			}
		}
	}

	for i := 0; i < block.x; i++ {
		for j := 0; j < block.y; j++ {
			for k := 0; k < block.z; k++ {
				if block.neurons[i][j][k].isIdle() {
					if block.neurons[i][j][k].Activity > 0.008 && block.neurons[i][j][k].Relax < block.neurons[i][j][k].Activity && rand.Float32() < 0.07 {
						fmt.Println(block.neurons[i][j][k].Relax)
						block.neurons[i][j][k].Status = 1
					}
				} else if block.neurons[i][j][k].isActive() {
					block.neurons[i][j][k].Status = 2
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

/* Matrix */

func (b *Block) New(f func(a, b, c int) Neuron) *Block {
	b.neurons = make([][][]Neuron, b.x)
	for i := 0; i < b.x; i++ {
		b.neurons[i] = make([][]Neuron, b.y)
		for j := 0; j < b.y; j++ {
			b.neurons[i][j] = make([]Neuron, b.z)
			for k := 0; k < b.z; k++ {
				b.neurons[i][j][k] = f(i, j, k)
			}
		}
	}
	return b
}

func (b *Block) Each(f func(a, b, c int, value Neuron)) *Block {
	for i := 0; i < b.x; i++ {
		for j := 0; j < b.y; j++ {
			for k := 0; k < b.z; k++ {
				f(i, j, k, b.neurons[i][j][k])
			}
		}
	}
	return b
}
