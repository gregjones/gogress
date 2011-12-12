package gogress

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ProgressBar struct {
	Max int64
	Fd io.Writer
	Current int64
	Finished bool
	Widgets []Widget
	termWidth int
}

func (p *ProgressBar) Start() {
	p.Update(0)
	p.Widgets = []Widget{NewPercentageWidget(),NewAnimatedMarker(), NewBarWidget()}
}

func (p ProgressBar) Update(i int64) {
	p.Current = i
	if p.Current == p.Max {
		p.Finished = true
	}
	lines := []string{}
	width := p.termWidth
	for _, w := range(p.Widgets) {
		line := w.Update(&p, width)
		width -= len(line)
		lines = append(lines, line)
	}
	lines = append(lines, "\r")
	result := strings.Join(lines, " ")
	p.Fd.Write([]byte(result))
}

func (p *ProgressBar) Percentage() int64 {
	return int64(float64(p.Current) / float64(p.Max) * 100)
}

func (p *ProgressBar) Proportion() float64 {
	return float64(p.Current) / float64(p.Max)
}

func NewProgressBar(max int64) *ProgressBar {
	return &ProgressBar{Max: max, Fd: os.Stderr, termWidth: 80}
}

type Widget interface {
	Update(p *ProgressBar, width int) (line string)
}

type BarWidget struct {
	Marker string
	Left string
	Right string
	Fill string
}

func NewBarWidget() *BarWidget {
	return &BarWidget{Marker: "#", Left: "|", Right: "|", Fill: " "}
}

func (w *BarWidget) Update(p *ProgressBar, width int) string {
	line := "%s%s%s"
	width -= len(w.Left) + len(w.Right)
	fill := strings.Repeat(w.Marker, int(p.Proportion() * float64(width)))
	return fmt.Sprintf(line, w.Left, fill, w.Right)
}


type PercentageWidget struct {}

func (w *PercentageWidget) Update(p *ProgressBar, width int) string {
	line := fmt.Sprintf("%3d%%", p.Percentage())
	return line
}

func NewPercentageWidget() *PercentageWidget {
	return &PercentageWidget{}
}


type AnimatedMarker struct {
	Markers []byte
	Current int
}

func (w *AnimatedMarker) Update(p *ProgressBar, width int) string {
	if p.Finished {
		return string(w.Markers[0])
	}
	w.Current = (w.Current + 1)	% len(w.Markers)
	return string(w.Markers[w.Current])
}

func NewAnimatedMarker() *AnimatedMarker {
	return &AnimatedMarker{Markers: []byte{'|', '/', '-', '\\'}, Current: -1}
}