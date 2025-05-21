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
		return 0, fmt.Errorf("rse > 5%%")
	}
	return avg, nil
}

func DrawGraph(pathDir string, start int, stop int, step int) error {
	if start > stop || step <= 0 {
		return fmt.Errorf("DrawGraph: error start stop step params")
	}

	cntRowsVals := make([]int, 0)
	notIndexVals := make([]float64, 0)
	indexVals := make([]float64, 0)

	fmt.Printf("cnt rows | not Index | Index\n")
	for i := start; i < stop; i += step {
		filenameNotIndex := filepath.Join(cnfg.GetProjectRoot(), fmt.Sprintf("%s/data/%d_notIndex.txt", pathDir, i))
		tmNotIndex, err := readAvg(filenameNotIndex)
		if err != nil {
			return fmt.Errorf("DrawGraph: %v", err)
		}
		filenameIndex := filepath.Join(cnfg.GetProjectRoot(), fmt.Sprintf("%s/data/%d_Index.txt", pathDir, i))
		tmIndex, err := readAvg(filenameIndex)
		if err != nil {
			return fmt.Errorf("DrawGraph: %v", err)
		}
		cntRowsVals = append(cntRowsVals, i)
		notIndexVals = append(notIndexVals, tmNotIndex)
		indexVals = append(indexVals, tmIndex)
		fmt.Printf("%d | %.3f | %.3f \n", i, tmNotIndex, tmIndex)
	}

	// Создаем группированные столбцы
	w := vg.Points(15) // Ширина столбцов

	// Преобразуем данные в формат для гистограммы
	notIndexPoints := make(plotter.Values, len(notIndexVals))
	indexPoints := make(plotter.Values, len(indexVals))
	for i := range notIndexVals {
		notIndexPoints[i] = notIndexVals[i]
		indexPoints[i] = indexVals[i]
	}

	// Создаем столбцы для notIndexVals
	notIndexBars, err := plotter.NewBarChart(notIndexPoints, w)
	if err != nil {
		panic(err)
	}
	notIndexBars.Color = color.RGBA{R: 200, G: 200, B: 200, A: 255} // Оранжевый цвет
	notIndexBars.Offset = -w / 2                                    // Смещение влево

	// Создаем столбцы для indexVals
	indexBars, err := plotter.NewBarChart(indexPoints, w)
	if err != nil {
		panic(err)
	}
	indexBars.Color = color.RGBA{R: 128, G: 128, B: 128, A: 255} // Синий цвет
	indexBars.Offset = -w / 2                                    // Смещение влево

	// График
	p := plot.New()
	p.Title.Text = "Зависимость времени выполнения запроса от количества записей в таблице"
	p.Title.TextStyle.XAlign = draw.XCenter
	p.Title.TextStyle.YAlign = draw.YTop
	p.Title.Padding = vg.Points(10)
	p.X.Label.Text = "Количество записей в таблице artwork_event"
	p.Y.Label.Text = "Время, мс"
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.Add("Без индекса", notIndexBars)
	p.Legend.Add("С индексом", indexBars)
	// p.Add(plotter.NewGrid())
	fontText := font.Font{
		Size:     20,            // Размер шрифта в пунктах (1/72 дюйма)
		Typeface: "Times-Roman", // Название шрифта
		Variant:  "Sans",
	}
	p.Title.TextStyle.Font = fontText
	p.X.Label.TextStyle.Font = fontText
	p.Y.Label.TextStyle.Font = fontText
	p.X.Tick.Label.Font = fontText
	p.Y.Tick.Label.Font = fontText
	p.Legend.TextStyle.Font = fontText
	p.Add(notIndexBars, indexBars)

	// Настраиваем метки по оси X
	labels := make([]string, len(cntRowsVals))
	for i, val := range cntRowsVals {
		if val%10000 == 0 {
			labels[i] = fmt.Sprint(val)
		} else {
			labels[i] = ""
		}

	}
	p.NominalX(labels...)

	saveFile := filepath.Join(cnfg.GetProjectRoot(), fmt.Sprintf("%s/histogram.png", pathDir))
	if err := p.Save(12*vg.Inch, 7*vg.Inch, saveFile); err != nil {
		return fmt.Errorf("DrawGraph: %v", err)
	}
	return nil
}
