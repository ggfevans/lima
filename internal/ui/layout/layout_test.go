package layout

import "testing"

func TestTwoPanelWide(t *testing.T) {
	d := Calculate(120, 40)

	if d.Mode != TwoPanel {
		t.Errorf("expected TwoPanel mode, got %d", d.Mode)
	}
	if d.SidebarWidth != 0 {
		t.Errorf("expected SidebarWidth=0, got %d", d.SidebarWidth)
	}
	if d.ContentHeight != 38 {
		t.Errorf("expected ContentHeight=38, got %d", d.ContentHeight)
	}
	if d.ConvListWidth+d.ThreadWidth != 120 {
		t.Errorf("expected widths to sum to 120, got %d+%d=%d",
			d.ConvListWidth, d.ThreadWidth,
			d.ConvListWidth+d.ThreadWidth)
	}
}

func TestTwoPanelNarrow(t *testing.T) {
	d := Calculate(70, 30)

	if d.Mode != TwoPanel {
		t.Errorf("expected TwoPanel mode, got %d", d.Mode)
	}
	if d.SidebarWidth != 0 {
		t.Errorf("expected SidebarWidth=0, got %d", d.SidebarWidth)
	}
	if d.ContentHeight != 28 {
		t.Errorf("expected ContentHeight=28, got %d", d.ContentHeight)
	}
	if d.ConvListWidth+d.ThreadWidth != 70 {
		t.Errorf("expected ConvListWidth+ThreadWidth=70, got %d+%d=%d",
			d.ConvListWidth, d.ThreadWidth,
			d.ConvListWidth+d.ThreadWidth)
	}
}

func TestSinglePanelMode(t *testing.T) {
	d := Calculate(50, 20)

	if d.Mode != SinglePanel {
		t.Errorf("expected SinglePanel mode, got %d", d.Mode)
	}
	if d.SidebarWidth != 0 {
		t.Errorf("expected SidebarWidth=0, got %d", d.SidebarWidth)
	}
	if d.ConvListWidth != 0 {
		t.Errorf("expected ConvListWidth=0, got %d", d.ConvListWidth)
	}
	if d.ThreadWidth != 50 {
		t.Errorf("expected ThreadWidth=50, got %d", d.ThreadWidth)
	}
	if d.ContentHeight != 18 {
		t.Errorf("expected ContentHeight=18, got %d", d.ContentHeight)
	}
}

func TestContentHeightAllModes(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"TwoPanelWide", 120, 40},
		{"TwoPanelNarrow", 70, 30},
		{"SinglePanel", 50, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Calculate(tt.width, tt.height)
			expected := tt.height - 2
			if d.ContentHeight != expected {
				t.Errorf("expected ContentHeight=%d (height-2), got %d", expected, d.ContentHeight)
			}
		})
	}
}

func TestVerySmallTerminal(t *testing.T) {
	d := Calculate(10, 5)

	if d.ContentHeight < 1 {
		t.Errorf("expected ContentHeight >= 1, got %d", d.ContentHeight)
	}
}

func TestBoundaryModeTransitions(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected LayoutMode
	}{
		{"width=90 is TwoPanel", 90, TwoPanel},
		{"width=60 is TwoPanel", 60, TwoPanel},
		{"width=59 is SinglePanel", 59, SinglePanel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Calculate(tt.width, 30)
			if d.Mode != tt.expected {
				t.Errorf("expected mode %d for width=%d, got %d", tt.expected, tt.width, d.Mode)
			}
		})
	}
}

func TestWidthSumsTwoPanel(t *testing.T) {
	widths := []int{60, 70, 80, 90, 100, 120, 200}
	for _, w := range widths {
		d := Calculate(w, 30)
		if d.Mode != TwoPanel {
			t.Fatalf("expected TwoPanel for width=%d, got %d", w, d.Mode)
		}
		total := d.ConvListWidth + d.ThreadWidth
		if total != w {
			t.Errorf("TwoPanel width=%d: ConvListWidth(%d)+ThreadWidth(%d)=%d, want %d",
				w, d.ConvListWidth, d.ThreadWidth, total, w)
		}
	}
}
