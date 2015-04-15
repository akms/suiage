package compress


import (
	"os"
	ui "github.com/gizak/termui"
)



func PostFileinfo(info os.FileInfo,g0 *ui.Gauge) {
	g0.Border.Label = info.Name()
}
