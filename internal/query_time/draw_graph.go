package querytime

import (
	"bufio"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func readAvg(filename string) (float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("DrawGraph: %v", err)
	}
	defer file.Close()

	var sumTimeMeasures float64 = 0
	var cntTimeMeasures int = 0
	timeMeasures := make([]float64, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tm, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			return 0, fmt.Errorf("DrawGraph: err convert to float")
		}
		timeMeasures = append(timeMeasures, tm)
		sumTimeMeasures += tm
		cntTimeMeasures++
	}
	if cntTimeMeasures == 0 {
		return 0, fmt.Errorf("readAvg: empty file")
	}
	avg := sumTimeMeasures / float64(cntTimeMeasures)

	var sum2 float64 = 0
	for _, val := range timeMeasures {
		sum2 += (val - avg) * (val - avg)
	}
	D := float64(sum2) / float64(len(timeMeasures)-1)
	stdErr := math.Sqrt(D) / math.Sqrt(float64(len(timeMeasures)))
	rse := stdErr / avg
	if rse > 5 {
		return 0, fmt.Errorf("rse > 5%")
	}
	return avg, nil
}

func DrawGraph(start int, stop int, step int) error {
	if start > stop || step <= 0 {
		return fmt.Errorf("DrawGraph: error start stop step params")
	}

	cntRowsVals := make([]int, 0)
	notIndexVals := make([]float64, 0)
	indexVals := make([]float64, 0)

	for i := start; i < stop; i += step {
		filenameNotIndex := filepath.Join(cnfg.GetProjectRoot(), fmt.Sprintf("/measure_results/%d_notIndex.txt", i))
		tmNotIndex, err := readAvg(filenameNotIndex)
		if err != nil {
			return fmt.Errorf("DrawGraph: %v", err)
		}
		filenameIndex := filepath.Join(cnfg.GetProjectRoot(), fmt.Sprintf("/measure_results/%d_Index.txt", i))
		tmIndex, err := readAvg(filenameIndex)
		if err != nil {
			return fmt.Errorf("DrawGraph: %v", err)
		}
		cntRowsVals = append(cntRowsVals, i)
		notIndexVals = append(notIndexVals, tmNotIndex)
		indexVals = append(indexVals, tmIndex)
	}

	pnotIndex := make(plotter.XYs, len(cntRowsVals))
	pIndex := make(plotter.XYs, len(cntRowsVals))
	for i := range cntRowsVals {
		pnotIndex[i].X = float64(cntRowsVals[i])
		pnotIndex[i].Y = notIndexVals[i]
		pIndex[i].X = float64(cntRowsVals[i])
		pIndex[i].Y = indexVals[i]
	}

	p := plot.New()
	p.Title.Text = "Зависимость времени выполнения запроса от количества записей в таблице"
	p.Title.TextStyle.XAlign = draw.XCenter
	p.Title.TextStyle.YAlign = draw.YTop
	p.Title.Padding = vg.Points(10)
	p.X.Label.Text = "Количество записей в таблицу artworks_event"
	p.Y.Label.Text = "Время, мс"
	p.Legend.Top = true
	p.Legend.Left = true
	p.Add(plotter.NewGrid())

	fontText := font.Font{
		Size:     14,            // Размер шрифта в пунктах (1/72 дюйма)
		Typeface: "Times-Roman", // Название шрифта
		Variant:  "Sans",
	}
	p.Title.TextStyle.Font = fontText
	p.X.Label.TextStyle.Font = fontText
	p.Y.Label.TextStyle.Font = fontText
	p.X.Tick.Label.Font = fontText
	p.Y.Tick.Label.Font = fontText

	widthLines := vg.Points(1)
	scatterRadius := vg.Points(4)
	colorNotIndex := color.RGBA{R: 244, G: 67, B: 54, A: 255}
	colorIndex := color.RGBA{R: 76, G: 175, B: 80, A: 255}

	lineNotIndex, err := plotter.NewLine(pnotIndex)
	if err != nil {
		return fmt.Errorf("DrawGraph: %v", err)
	}
	lineNotIndex.LineStyle.Width = widthLines
	lineNotIndex.LineStyle.Color = colorNotIndex
	scatterNotIndex, err := plotter.NewScatter(pnotIndex)
	if err != nil {
		return fmt.Errorf("DrawGraph: %v", err)
	}
	scatterNotIndex.GlyphStyle = draw.GlyphStyle{
		Color:  colorNotIndex,
		Shape:  draw.CircleGlyph{},
		Radius: scatterRadius,
	}

	lineIndex, err := plotter.NewLine(pIndex)
	if err != nil {
		return fmt.Errorf("DrawGraph: %v", err)
	}
	lineIndex.LineStyle.Width = widthLines
	lineIndex.LineStyle.Color = colorIndex
	scatterIndex, err := plotter.NewScatter(pIndex)
	if err != nil {
		return fmt.Errorf("DrawGraph: %v", err)
	}
	scatterIndex.GlyphStyle = draw.GlyphStyle{
		Color:  colorIndex,
		Shape:  draw.PyramidGlyph{},
		Radius: scatterRadius,
	}

	p.Add(lineNotIndex, scatterNotIndex)
	p.Add(lineIndex, scatterIndex)
	p.Legend.Add("Без индекса", lineNotIndex, scatterNotIndex)
	p.Legend.Add("С индексом", lineIndex, scatterIndex)

	// // -----
	// err = plotutil.AddLines(p,
	// 	"Без индекса", pnotIndex,
	// 	"С индексом", pIndex,
	// )
	// if err != nil {
	// 	return fmt.Errorf("DrawGraph: %v", err)
	// }

	saveFile := filepath.Join(cnfg.GetProjectRoot(), "/measure_results/graph.svg")
	if err := p.Save(8*vg.Inch, 5*vg.Inch, saveFile); err != nil {
		return fmt.Errorf("DrawGraph: %v", err)
	}
	return nil
}
