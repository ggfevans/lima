package layout

// Dimensions holds computed dimensions for each panel.
type Dimensions struct {
	Width  int
	Height int

	// Panel widths (0 means collapsed)
	SidebarWidth  int
	ConvListWidth int
	ThreadWidth   int

	// Content height (total minus header and status bar)
	ContentHeight int

	// Layout mode
	Mode LayoutMode
}

type LayoutMode int

const (
	// ThreePanel: sidebar + conversation list + thread (>= 90 cols)
	ThreePanel LayoutMode = iota
	// TwoPanel: conversation list + thread (60-89 cols)
	TwoPanel
	// SinglePanel: only the focused panel (< 60 cols)
	SinglePanel
)

const (
	headerHeight    = 1
	statusBarHeight = 1
	minSidebarWidth = 12
	sidebarWidth    = 14
)

// Calculate computes panel dimensions from terminal size.
func Calculate(width, height int) Dimensions {
	d := Dimensions{
		Width:         width,
		Height:        height,
		ContentHeight: max(height-headerHeight-statusBarHeight, 1),
	}

	switch {
	case width >= 60:
		d.Mode = TwoPanel
		d.SidebarWidth = 0
		d.ConvListWidth = width * 30 / 100
		if d.ConvListWidth < 20 {
			d.ConvListWidth = 20
		}
		d.ThreadWidth = width - d.ConvListWidth
	default:
		d.Mode = SinglePanel
		d.SidebarWidth = 0
		d.ConvListWidth = 0
		d.ThreadWidth = width
	}

	return d
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
