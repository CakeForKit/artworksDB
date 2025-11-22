package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stats struct {
	Min   float64
	Max   float64
	Avg   float64
	Count int
	Sum   float64
}

func (s *Stats) PrintStats() {
	fmt.Printf("📊 Статистика по второй колонке:\n")
	fmt.Printf("   Количество значений: %d\n", s.Count)
	fmt.Printf("   Минимальное: %.4f\n", s.Min)
	fmt.Printf("   Максимальное: %.4f\n", s.Max)
	fmt.Printf("   Среднее: %.4f\n", s.Avg)
	fmt.Printf("   Сумма: %.4f\n", s.Sum)
}

// CalculateStatsFromFile вычисляет статистику по второй колонке файла
func CalculateStatsFromFile(filename string) (*Stats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to open file: %v", err)
	}
	defer file.Close()

	stats := &Stats{
		Min: 1<<63 - 1,
		Max: -1 << 63,
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		columns := strings.Fields(line)
		if len(columns) < 2 {
			return nil, fmt.Errorf("❌ Invalid format on line %d: expected 2 columns, got %d", lineNumber, len(columns))
		}

		value, err := strconv.ParseFloat(columns[1], 64)
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to parse number on line %d: %v", lineNumber, err)
		}

		stats.Count++
		stats.Sum += value

		if value < stats.Min {
			stats.Min = value
		}
		if value > stats.Max {
			stats.Max = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("❌ Error reading file: %v", err)
	}

	if stats.Count == 0 {
		return nil, fmt.Errorf("❌ No valid data found in file")
	}

	stats.Avg = stats.Sum / float64(stats.Count)

	return stats, nil
}

func main() {
	// Пример использования
	filename := "./metrics_data/prometheus/container_cpu_usage_seconds_total.txt"

	// Базовая статистика
	stats, err := CalculateStatsFromFile(filename)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	stats.PrintStats()
}
