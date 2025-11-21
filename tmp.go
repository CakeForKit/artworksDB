package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/tracing"
)

func main() {
	analyzer := tracing.NewTraceAnalyzer("http://localhost:16686")

	startStr := "2025-11-21T20:29:40Z"
	endStr := "2025-11-21T20:29:52Z"
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		panic(err)
	}

	analysis, err := analyzer.AnalyzeTimeRange(
		context.Background(),
		"artworks-timing",
		startTime,
		endTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Выводим результаты
	analyzer.PrintTimeRangeAnalysis(analysis)

	main1()
}

func main1() {
	// Создаем анализатор
	analyzer := tracing.NewTraceAnalyzer("http://localhost:16686")
	ctx := context.Background()

	fmt.Println("1. Testing Jaeger connection...")
	if err := analyzer.TestConnection(ctx); err != nil {
		log.Printf("Connection test failed: %v", err)
		return
	} else {
		fmt.Println("✓ Connection successful")
	}

	fmt.Println("\n2. Fetching services...")
	services, err := analyzer.GetServices(ctx)
	if err != nil {
		log.Printf("Failed to get services: %v", err)
	} else {
		fmt.Printf("Available services: %v\n", services)
	}

	serviceName := "artworks-timing"
	fmt.Printf("\n3. Trying to fetch traces for service: %s\n", serviceName)
	traces, err := analyzer.FetchTraces(ctx, serviceName, 2*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	// Выводим summary
	analyzer.PrintTraceSummary(traces)

	// Анализируем конкретные операции
	databaseSpans := analyzer.FindSpansByOperation(traces, "HTTP GET /api/v1/museum/artworks")
	fmt.Printf("Found %d database queries\n", len(databaseSpans))

	analyzeSpanAverages(analyzer, traces)
	// Анализируем производительность
	// for _, span := range databaseSpans {
	// 	if span.Duration > 100*time.Millisecond {
	// 		fmt.Printf("Slow query: %s took %v\n", span.OperationName, span.Duration)
	// 	}
	// }
}

// analyzeSpanAverages вычисляет и выводит среднюю статистику для каждого уникального спана
func analyzeSpanAverages(analyzer *tracing.TraceAnalyzer, traces []*tracing.Trace) {
	// Структура для хранения статистики по каждому operation name
	type SpanStats struct {
		OperationName string
		Count         int
		TotalDuration time.Duration
		MinDuration   time.Duration
		MaxDuration   time.Duration
		Durations     []time.Duration // для медианы
	}

	// Собираем статистику по всем спанам
	statsMap := make(map[string]*SpanStats)

	for _, trace := range traces {
		for _, span := range trace.Spans {
			if _, exists := statsMap[span.OperationName]; !exists {
				statsMap[span.OperationName] = &SpanStats{
					OperationName: span.OperationName,
					MinDuration:   span.Duration,
					MaxDuration:   span.Duration,
					Durations:     []time.Duration{},
				}
			}

			stats := statsMap[span.OperationName]
			stats.Count++
			stats.TotalDuration += span.Duration
			stats.Durations = append(stats.Durations, span.Duration)

			if span.Duration < stats.MinDuration {
				stats.MinDuration = span.Duration
			}
			if span.Duration > stats.MaxDuration {
				stats.MaxDuration = span.Duration
			}
		}
	}

	// Выводим статистику
	fmt.Println("=== Average Time per Span ===")
	fmt.Printf("%-40s %8s %8s %8s %8s %8s\n",
		"Operation Name", "Count", "Avg", "Min", "Max", "Median")
	fmt.Printf("%-40s %8s %8s %8s %8s %8s\n",
		"-------------", "-----", "---", "---", "---", "------")

	for opName, stats := range statsMap {
		avgDuration := stats.TotalDuration / time.Duration(stats.Count)
		medianDuration := calculateMedian(stats.Durations)

		fmt.Printf("%-40s %8d %8v %8v %8v %8v\n",
			truncateString(opName, 40),
			stats.Count,
			formatDuration(avgDuration),
			formatDuration(stats.MinDuration),
			formatDuration(stats.MaxDuration),
			formatDuration(medianDuration))
	}
}

// calculateMedian вычисляет медиану длительностей
func calculateMedian(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Создаем копию для сортировки
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)

	// Простая bubble sort (для небольших наборов данных)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	// Вычисляем медиану
	middle := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[middle-1] + sorted[middle]) / 2
	}
	return sorted[middle]
}

// formatDuration форматирует длительность для красивого вывода
func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.2fns", float64(d))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d)/float64(time.Microsecond))
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2fs", float64(d)/float64(time.Second))
}

// truncateString обрезает строку если она слишком длинная
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

/*
func mainEmail() {
	emailCnfg := cnfg.LoadEmailCnfg()
	s := emailserv.NewEmailService(
		emailCnfg.Host,
		emailCnfg.Port,
		emailCnfg.Username,
		emailCnfg.Password,
		emailCnfg.From,
	)
	fmt.Printf("mail: %s, password: %s\n", emailCnfg.From, emailCnfg.Password)
	err := s.SendEmail([]string{"tmpforread@mail.ru"}, "subject", "body2")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	emailReaderCnfg := cnfg.LoadEmailReaderCnfg()

	emailReader := emailreader.NewEmailReader(
		emailReaderCnfg.Host,     // IMAP хост
		emailReaderCnfg.Port,     // IMAP порт
		emailReaderCnfg.Username, // Email
		emailReaderCnfg.Password, // Пароль
	)

	emails, err := emailReader.FindEmailByCriteria(emailreader.SearchCriteria{
		From:    emailCnfg.From,
		Subject: "subject",
	})
	if err != nil {
		log.Fatalf("Ошибка чтения писем: %v", err)
	}
	emailReader.PrintEmail(emails)
}
*/
