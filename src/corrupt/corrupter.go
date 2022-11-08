// package corrupt contains functions that scramble byte data for the purpose of data corruption
package corrupt

import (
	"math/rand"
	"time"
)

type Corrupter struct {
	data     []byte
	r        *rand.Rand
	strength int // 1/strength chance for replace/defect.
}

func New(data []byte) *Corrupter {
	return &Corrupter{
		data:     data,
		strength: 1e4,
		r:        rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (c *Corrupter) SetRand(r *rand.Rand) *Corrupter {
	c.r = r
	return c
}

func (c *Corrupter) SetStrength(strength int) *Corrupter {
	c.strength = strength
	return c
}

// Transpose rearranges parts of data
func (c *Corrupter) Transpose() *Corrupter {
	pieces := make([][]byte, 0)
	last := 0
	for i := range c.data {
		if c.r.Intn(c.strength) == 0 {
			pieces = append(pieces, c.data[last:i])
			last = i
		}
	}
	pieces = append(pieces, c.data[last:])

	c.r.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})

	acc := 0
	for i, piece := range pieces {
		for j := range piece {
			c.data[acc] = pieces[i][j]
			acc++
		}
	}

	return c
}

// Replace rewrites random parts of data
func (c *Corrupter) Replace() *Corrupter {
	for i := range c.data {
		if c.r.Intn(c.strength) == 0 {
			c.data[i] = byte(c.r.Intn(256))
		}
	}

	return c
}

// Delete deletes random parts of data
func (c *Corrupter) Delete() *Corrupter {
	j := 0
	for _, b := range c.data {
		if c.r.Intn(c.strength) != 0 {
			c.data[j] = b
			j++
		}
	}
	c.data = c.data[:j]

	return c
}

// Delete replaces random parts of data with empty strings
func (c *Corrupter) Defect() *Corrupter {
	for i := range c.data {
		if c.r.Intn(c.strength) == 0 {
			c.data[i] = byte(0)
		}
	}

	return c
}

func (c *Corrupter) Data() []byte {
	return c.data
}
