package gogress

import (
	. "launchpad.net/gocheck"
	"testing"
)

type writer struct {
	Count  int
	Writes [][]byte
}

func (w *writer) Write(p []byte) (int, error) {
	w.Count++
	w.Writes = append(w.Writes, p)
	return len(p), nil
}

func Test(t *testing.T) {
	DefaultWriter = new(writer)
	TestingT(t)
}

type S struct{}

var _ = Suite(&S{})

func (s *S) TestDefault(c *C) {
	p := NewProgressBar(100)
	p.Start()
	var i int64
	for i = 0; i < 100; i++ {
		p.Update(i + 1)
	}
}

func (s *S) TestPercentage(c *C) {
	w := NewPercentageWidget()
	b := new(writer)
	p := &ProgressBar{Max: 100, Fd: b, Widgets: []Widget{w}}
	p.Start()
	c.Assert(b.Count, Equals, 1)

	c.Assert(b.Writes[0], Equals, []byte{' ', ' ', '0', '%', ' ', '\r'})

	p.Update(1)
	c.Assert(b.Count, Equals, 2)
	c.Assert(b.Writes[1], Equals, []byte{' ', ' ', '1', '%', ' ', '\r'})
	var i int64
	for i = 0; i < 15; i++ {
		p.Update(i + 2)
	}
	c.Assert(b.Count, Equals, 17)
	c.Assert(b.Writes[16], Equals, []byte{' ', '1', '6', '%', ' ', '\r'})

	p.Update(100)
	c.Assert(b.Count, Equals, 18)
	c.Assert(b.Writes[17], Equals, []byte{'1', '0', '0', '%', ' ', '\r'})

	p.Update(200)
	c.Assert(b.Writes[18], Equals, []byte{'1', '0', '0', '%', ' ', '\r'})
}
