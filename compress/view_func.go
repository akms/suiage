package compress

import (
	"github.com/gizak/termui"
	"time"
)

func View() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.UseTheme("helloworld")

	par1 := termui.NewPar("")
	par1.Height = 3
	par1.Width = 37
	par1.Y = 9
	par1.Border.Label = "Suiage Dir"

	par2 := termui.NewPar("")
	par2.Height = 3
	par2.Width = 37
	par2.Y = 9
	par2.Border.Label = "Suiage File"

	g0 := termui.NewGauge()
	g0.Percent = 0
	g0.Width = 100
	g0.Height = 3
	g0.Border.Label = "Suiage Gauge"
	g0.BarColor = termui.ColorRed
	g0.Border.FgColor = termui.ColorWhite
	g0.Border.LabelFgColor = termui.ColorCyan

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(4, 0, par1),
			termui.NewCol(4, 0, par2)),
		termui.NewRow(
			termui.NewCol(8, 0, g0)))
	
	termui.Body.Align()

	draw := func(t int) {
		g0.Percent = t
		termui.Render(termui.Body)
	}

	evt := termui.EventCh()
	i := 0
	for {
		select {
		case i := <-fin :
			par1.Text = i
			par2.Text = i
			termui.Render(termui.Body)
			return
		case str := <-ed:
			par1.Text = str
		case str := <-ef:
			par2.Text = str
		case e := <-evt:
			if e.Type == termui.EventKey && e.Ch == 'q' {
				return
			}			
		default:
			i++
			draw(i)
			if g0.Percent >= 100 {
				g0.Percent = 0
				termui.Body.Align()
			}
			termui.Render(termui.Body)
			time.Sleep(time.Second / 2)
		}
	}
}
