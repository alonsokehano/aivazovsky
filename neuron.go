package main

import (
	"math/rand"
)

/*
	Neuron states:
	0 - idle
	1 - active (spiking)
	2 - relaxing
*/

type Neuron struct {
	x, y, z         int
	value, newvalue float32
	weights         [][][]float32
	state           int
}

func (n *Neuron) Initialize(config BlockConfig) {
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

func (n *Neuron) IsIdle() bool {
	return n.state == 0
}

func (n *Neuron) IsActive() bool {
	return n.state == 1
}

func (n *Neuron) IsRelaxing() bool {
	return n.state == 2
}

func (n *Neuron) SetValue(value float32, config BlockConfig) {
	n.value = value
	if value >= config.synapses_threshold {
		n.state = 1
	} else if n.IsActive() && value >= config.relaxation_threshold {
		n.state = 2
	} else if n.IsRelaxing() && value <= config.relaxation_threshold {
		n.state = 0
	}
}
