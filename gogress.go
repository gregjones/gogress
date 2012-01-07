package gogress

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	termWidth = 80
)

var (
	defaultWidgets = []Widget{NewPercentageWidget(),NewAnimatedMarker(), NewBarWidget()}
	DefaultWriter io.Writer = os.Stderr
)

// Progress bar contains the state (Current, Max) and a list of Widgets
type ProgressBar struct {
	Max int64
	Fd io.Writer
	Current int64
	Finished bool
	Widgets []Widget
	termWidth int
}

// Update with a value of zero
func (p *ProgressBar) Start() {
	p.Update(0)
}

// Update takes a number and iterates through the bar's widgets updating each
func (p *ProgressBar) Update(i int64) {
	p.Current = i
	if p.Current == p.Max {
		p.Finished = true
	}
	lines := []string{}
	width := p.termWidth
	for _, w := range(p.Widgets) {
		line := w.Update(p, width)
		width -= len(line)
		lines = append(lines, line)
	}
	lines = append(lines, "\r")
	result := strings.Join(lines, " ")
	p.Fd.Write([]byte(result))
}


// Progress returns a number between 0 and 1 (incl) indicating how far through we are
func (p *ProgressBar) Progress() float64 {
	if(p.Current >= p.Max) {
		return 1
	}
	return float64(p.Current) / float64(p.Max)
}

// NewProgressBar returns a ProgressBar with some sensible defaults. namely:
// Will write to DefaultWriter, with a copy of defaultWidgets
func NewProgressBar(max int64) *ProgressBar {
	widgets := make([]Widget, len(defaultWidgets))
	copy(widgets, defaultWidgets)
	return &ProgressBar{Max: max, Fd: DefaultWriter, termWidth: termWidth, Widgets: widgets}
}

// A Widget takes a ProgressBar instance and a remaining-width and returns a string indicating the new progress
type Widget interface {
	Update(p *ProgressBar, width int) (line string)
}

// A BarWidget will fill from left to right
type BarWidget struct {
	Marker string
	Left string
	Right string
	Fill string
}

// NewBarWidget returns a BarWidget with some sensible defaults
func NewBarWidget() *BarWidget {
	return &BarWidget{Marker: "#", Left: "|", Right: "|", Fill: " "}
}

func (w *BarWidget) Update(p *ProgressBar, width int) string {
	line := "%s%s%s"
	width -= len(w.Left) + len(w.Right)
	fill := strings.Repeat(w.Marker, int(p.Progress() * float64(width)))
	return fmt.Sprintf(line, w.Left, fill, w.Right)
}

// A Percentage widget will write the progress as a percentage
type PercentageWidget struct {}

func (w *PercentageWidget) Update(p *ProgressBar, width int) string {
	line := fmt.Sprintf("%3d%%", int(p.Progress() * 100))
	return line
}

// NewPercentageWidget returns a PercentageWidget instance
func NewPercentageWidget() *PercentageWidget {
	return &PercentageWidget{}
}

// AnimatedMarker will cycle through a list of characters writing each to give the effect of animation
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

// NewAnimatedMarker returns an AnimatedMarker instance setup using |/-\
func NewAnimatedMarker() *AnimatedMarker {
	return &AnimatedMarker{Markers: []byte{'|', '/', '-', '\\'}, Current: -1}
}