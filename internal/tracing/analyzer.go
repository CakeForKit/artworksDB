package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// TraceAnalyzer - утилита для извлечения и анализа трейсов из Jaeger
type TraceAnalyzer struct {
	jaegerURL  string
	httpClient *http.Client
}

// Trace - упрощенная структура трейса для анализа
type Trace struct {
	TraceID   string        `json:"traceID"`
	Spans     []Span        `json:"spans"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
}

// Span - упрощенная структура спана
type Span struct {
	TraceID       string            `json:"traceID"`
	SpanID        string            `json:"spanID"`
	ParentSpanID  string            `json:"parentSpanID,omitempty"`
	OperationName string            `json:"operationName"`
	StartTime     time.Time         `json:"startTime"`
	Duration      time.Duration     `json:"duration"`
	Tags          map[string]string `json:"tags"`
	ServiceName   string            `json:"serviceName"`
}

// JaegerAPIResponse - структура ответа Jaeger API
type JaegerAPIResponse struct {
	Data []*JaegerTrace `json:"data"`
}

// JaegerTrace - полная структура трейса из Jaeger API
type JaegerTrace struct {
	TraceID   string                    `json:"traceID"`
	Spans     []*JaegerSpan             `json:"spans"`
	Processes map[string]*JaegerProcess `json:"processes"`
}

// JaegerSpan - полная структура спана из Jaeger API
type JaegerSpan struct {
	TraceID       string             `json:"traceID"`
	SpanID        string             `json:"spanID"`
	OperationName string             `json:"operationName"`
	StartTime     int64              `json:"startTime"` // microseconds
	Duration      int64              `json:"duration"`  // microseconds
	Tags          []*JaegerTag       `json:"tags"`
	References    []*JaegerReference `json:"references"`
	ProcessID     string             `json:"processID"`
}

// JaegerProcess - информация о процессе из Jaeger
type JaegerProcess struct {
	ServiceName string       `json:"serviceName"`
	Tags        []*JaegerTag `json:"tags"`
}

// JaegerTag - тег спана/процесса
type JaegerTag struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type,omitempty"`
}

// JaegerReference - ссылки между спанами
type JaegerReference struct {
	RefType string `json:"refType"`
	TraceID string `json:"traceID"`
	SpanID  string `json:"spanID"`
}

// NewTraceAnalyzer создает новый анализатор трейсов
func NewTraceAnalyzer(jaegerURL string) *TraceAnalyzer {
	return &TraceAnalyzer{
		jaegerURL: jaegerURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchTraces получает трейсы за указанный период
func (a *TraceAnalyzer) FetchTraces(ctx context.Context, serviceName string, lookback time.Duration) ([]*Trace, error) {
	return a.FetchTracesWithLimit(ctx, serviceName, lookback, 100)
}

// FetchTracesWithLimit получает трейсы с ограничением по количеству
func (a *TraceAnalyzer) FetchTracesWithLimit(
	ctx context.Context, serviceName string, lookback time.Duration, limit int,
) ([]*Trace, error) {
	end := time.Now()
	start := end.Add(-lookback)

	// Формируем URL для Jaeger API
	baseURL := fmt.Sprintf("%s/api/traces", a.jaegerURL)
	params := url.Values{}
	params.Add("service", serviceName)
	params.Add("start", fmt.Sprintf("%d", start.UnixMicro()))
	params.Add("end", fmt.Sprintf("%d", end.UnixMicro()))
	params.Add("limit", fmt.Sprintf("%d", limit))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jaeger API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse JaegerAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return a.convertJaegerTraces(apiResponse.Data), nil
}

// FetchTraceByID получает конкретный трейс по ID
func (a *TraceAnalyzer) FetchTraceByID(ctx context.Context, traceID string) (*Trace, error) {
	url := fmt.Sprintf("%s/api/traces/%s", a.jaegerURL, traceID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jaeger API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse JaegerAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(apiResponse.Data) == 0 {
		return nil, fmt.Errorf("trace not found: %s", traceID)
	}

	traces := a.convertJaegerTraces(apiResponse.Data)
	if len(traces) == 0 {
		return nil, fmt.Errorf("failed to convert trace: %s", traceID)
	}

	return traces[0], nil
}

func (a *TraceAnalyzer) convertSpans(jaegerTrace *JaegerTrace, trace *Trace) {
	for _, jaegerSpan := range jaegerTrace.Spans {
		span := Span{
			TraceID:       jaegerSpan.TraceID,
			SpanID:        jaegerSpan.SpanID,
			OperationName: jaegerSpan.OperationName,
			StartTime:     time.Unix(0, jaegerSpan.StartTime*1000), // convert microseconds to nanoseconds
			Duration:      time.Duration(jaegerSpan.Duration) * time.Microsecond,
			Tags:          make(map[string]string),
		}

		// Получаем service name из процесса
		if process, exists := jaegerTrace.Processes[jaegerSpan.ProcessID]; exists {
			span.ServiceName = process.ServiceName
		}

		// Обрабатываем теги
		for _, tag := range jaegerSpan.Tags {
			if strValue, ok := a.tagValueToString(tag.Value); ok {
				span.Tags[tag.Key] = strValue
			}
		}

		// Обрабатываем родительские ссылки
		for _, ref := range jaegerSpan.References {
			if ref.RefType == "CHILD_OF" {
				span.ParentSpanID = ref.SpanID
				break
			}
		}

		trace.Spans = append(trace.Spans, span)
	}
}

// convertJaegerTraces конвертирует трейсы из формата Jaeger API во внутренний формат
func (a *TraceAnalyzer) convertJaegerTraces(jaegerTraces []*JaegerTrace) []*Trace {
	var traces []*Trace

	for _, jaegerTrace := range jaegerTraces {
		trace := &Trace{
			TraceID: jaegerTrace.TraceID,
			Spans:   []Span{},
		}

		// Конвертируем спаны
		a.convertSpans(jaegerTrace, trace)
		/*
			for _, jaegerSpan := range jaegerTrace.Spans {
				span := Span{
					TraceID:       jaegerSpan.TraceID,
					SpanID:        jaegerSpan.SpanID,
					OperationName: jaegerSpan.OperationName,
					StartTime:     time.Unix(0, jaegerSpan.StartTime*1000), // convert microseconds to nanoseconds
					Duration:      time.Duration(jaegerSpan.Duration) * time.Microsecond,
					Tags:          make(map[string]string),
				}

				// Получаем service name из процесса
				if process, exists := jaegerTrace.Processes[jaegerSpan.ProcessID]; exists {
					span.ServiceName = process.ServiceName
				}

				// Обрабатываем теги
				for _, tag := range jaegerSpan.Tags {
					if strValue, ok := a.tagValueToString(tag.Value); ok {
						span.Tags[tag.Key] = strValue
					}
				}

				// Обрабатываем родительские ссылки
				for _, ref := range jaegerSpan.References {
					if ref.RefType == "CHILD_OF" {
						span.ParentSpanID = ref.SpanID
						break
					}
				}

				trace.Spans = append(trace.Spans, span)
			}
		*/

		// Вычисляем общее время трейса
		if len(trace.Spans) > 0 {
			trace.StartTime = trace.Spans[0].StartTime
			var maxEndTime time.Time
			for _, span := range trace.Spans {
				endTime := span.StartTime.Add(span.Duration)
				if endTime.After(maxEndTime) {
					maxEndTime = endTime
				}
			}
			trace.Duration = maxEndTime.Sub(trace.StartTime)
		}

		traces = append(traces, trace)
	}

	return traces
}

// tagValueToString конвертирует значение тега в строку
func (a *TraceAnalyzer) tagValueToString(value interface{}) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case float64:
		return fmt.Sprintf("%g", v), true
	case bool:
		return fmt.Sprintf("%t", v), true
	case int:
		return fmt.Sprintf("%d", v), true
	case int64:
		return fmt.Sprintf("%d", v), true
	default:
		return fmt.Sprintf("%v", v), false
	}
}

// FindSpansByOperation находит все спаны с указанным именем операции
func (a *TraceAnalyzer) FindSpansByOperation(traces []*Trace, operationName string) []Span {
	var result []Span
	for _, trace := range traces {
		for _, span := range trace.Spans {
			if span.OperationName == operationName {
				result = append(result, span)
			}
		}
	}
	return result
}

// FindSpansByService находит все спаны для указанного сервиса
func (a *TraceAnalyzer) FindSpansByService(traces []*Trace, serviceName string) []Span {
	var result []Span
	for _, trace := range traces {
		for _, span := range trace.Spans {
			if span.ServiceName == serviceName {
				result = append(result, span)
			}
		}
	}
	return result
}

// GetTraceTree возвращает корневые спаны трейса (те, у которых нет родителя)
func (a *TraceAnalyzer) GetTraceTree(trace *Trace) []Span {
	var rootSpans []Span
	spanMap := make(map[string]Span)

	// Создаем мапу для быстрого поиска
	for _, span := range trace.Spans {
		spanMap[span.SpanID] = span
	}

	// Находим спаны без родителя
	for _, span := range trace.Spans {
		if span.ParentSpanID == "" || spanMap[span.ParentSpanID].SpanID == "" {
			rootSpans = append(rootSpans, span)
		}
	}

	return rootSpans
}

// PrintTraceSummary выводит краткую информацию о трейсах
func (a *TraceAnalyzer) PrintTraceSummary(traces []*Trace) {
	fmt.Printf("Found %d traces:\n", len(traces))
	for i, trace := range traces {
		fmt.Printf("Trace %d: ID=%s, Spans=%d, Duration=%v\n",
			i+1, trace.TraceID, len(trace.Spans), trace.Duration)

		rootSpans := a.GetTraceTree(trace)
		for _, span := range rootSpans {
			a.printSpanTree(span, trace.Spans, 1)
		}
		fmt.Println()
	}
}

// printSpanTree рекурсивно выводит дерево спанов
func (a *TraceAnalyzer) printSpanTree(current Span, allSpans []Span, indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	fmt.Printf("%s- %s (Service: %s, Duration: %v)\n",
		prefix, current.OperationName, current.ServiceName, current.Duration)

	// Находим дочерние спаны
	for _, span := range allSpans {
		if span.ParentSpanID == current.SpanID {
			a.printSpanTree(span, allSpans, indent+1)
		}
	}
}

// TestConnection проверяет соединение с Jaeger
func (a *TraceAnalyzer) TestConnection(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/services", a.jaegerURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jaeger API returned status %d", resp.StatusCode)
	}

	return nil
}

type ServicesResponse struct {
	Data   []string      `json:"data"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
	Errors []interface{} `json:"errors,omitempty"`
}

// GetServices получает список сервисов из Jaeger
func (a *TraceAnalyzer) GetServices(ctx context.Context) (*ServicesResponse, error) {
	url := fmt.Sprintf("%s/api/services", a.jaegerURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jaeger API returned status %d", resp.StatusCode)
	}

	var response ServicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("decode response: %w, body: %s", err, string(body))
	}

	return &response, nil
}

// -----------------------------------------

func (a *TraceAnalyzer) FetchTracesByTimeRange(
	ctx context.Context,
	serviceName string,
	startTime, endTime time.Time,
) ([]*Trace, error) {
	return a.FetchTracesByTimeRangeWithLimit(ctx, serviceName, startTime, endTime, 100)
}

// FetchTracesByTimeRangeWithLimit получает трейсы за временной отрезок с ограничением по количеству
func (a *TraceAnalyzer) FetchTracesByTimeRangeWithLimit(
	ctx context.Context,
	serviceName string,
	startTime, endTime time.Time,
	limit int,
) ([]*Trace, error) {
	// Валидация временного диапазона
	if startTime.After(endTime) {
		return nil, fmt.Errorf("start time cannot be after end time")
	}

	if endTime.After(time.Now()) {
		endTime = time.Now()
	}

	// Формируем URL для Jaeger API
	baseURL := fmt.Sprintf("%s/api/traces", a.jaegerURL)
	params := url.Values{}
	params.Add("service", serviceName)
	params.Add("start", fmt.Sprintf("%d", startTime.UnixMicro()))
	params.Add("end", fmt.Sprintf("%d", endTime.UnixMicro()))
	params.Add("limit", fmt.Sprintf("%d", limit))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jaeger API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse JaegerAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return a.convertJaegerTraces(apiResponse.Data), nil
}

// FetchTracesByTimeRangeMultiService получает трейсы для нескольких сервисов за временной отрезок
func (a *TraceAnalyzer) FetchTracesByTimeRangeMultiService(
	ctx context.Context,
	serviceNames []string,
	startTime, endTime time.Time,
	limit int,
) ([]*Trace, error) {
	var allTraces []*Trace

	for _, serviceName := range serviceNames {
		traces, err := a.FetchTracesByTimeRangeWithLimit(ctx, serviceName, startTime, endTime, limit)
		if err != nil {
			return nil, fmt.Errorf("fetch traces for service %s: %w", serviceName, err)
		}
		allTraces = append(allTraces, traces...)
	}

	return allTraces, nil
}

// AnalyzeTimeRange выполняет анализ трейсов за временной отрезок
func (a *TraceAnalyzer) AnalyzeTimeRange(
	ctx context.Context,
	serviceName string,
	startTime, endTime time.Time,
) (*TimeRangeAnalysis, error) {
	traces, err := a.FetchTracesByTimeRange(ctx, serviceName, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return a.AnalyzeTraces(traces), nil
}

// TimeRangeAnalysis - результат анализа временного отрезка
type TimeRangeAnalysis struct {
	TimeRange        TimeRangeStats   `json:"timeRange"`
	ServiceStats     ServiceStats     `json:"serviceStats"`
	OperationStats   []OperationStats `json:"operationStats"`
	PerformanceStats PerformanceStats `json:"performanceStats"`
	ErrorStats       ErrorStats       `json:"errorStats"`
}

// TimeRangeStats - статистика по временному отрезку
type TimeRangeStats struct {
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	Duration        string    `json:"duration"`
	TotalTraces     int       `json:"totalTraces"`
	TotalSpans      int       `json:"totalSpans"`
	TracesPerMinute float64   `json:"tracesPerMinute"`
	SpansPerMinute  float64   `json:"spansPerMinute"`
}

// ServiceStats - статистика по сервисам
type ServiceStats struct {
	Services         []string                 `json:"services"`
	ServiceCounts    map[string]int           `json:"serviceCounts"`
	ServiceDurations map[string]time.Duration `json:"serviceDurations"`
}

// OperationStats - статистика по операциям
type OperationStats struct {
	OperationName string        `json:"operationName"`
	ServiceName   string        `json:"serviceName"`
	Count         int           `json:"count"`
	AvgDuration   time.Duration `json:"avgDuration"`
	MinDuration   time.Duration `json:"minDuration"`
	MaxDuration   time.Duration `json:"maxDuration"`
	P95Duration   time.Duration `json:"p95Duration"`
	TotalDuration time.Duration `json:"totalDuration"`
}

// PerformanceStats - статистика производительности
type PerformanceStats struct {
	AvgTraceDuration  time.Duration    `json:"avgTraceDuration"`
	MinTraceDuration  time.Duration    `json:"minTraceDuration"`
	MaxTraceDuration  time.Duration    `json:"maxTraceDuration"`
	AvgSpanDuration   time.Duration    `json:"avgSpanDuration"`
	SlowestOperations []OperationStats `json:"slowestOperations"`
}

// ErrorStats - статистика ошибок
type ErrorStats struct {
	TotalErrors      int            `json:"totalErrors"`
	ErrorRate        float64        `json:"errorRate"`
	ErrorByService   map[string]int `json:"errorByService"`
	ErrorByOperation map[string]int `json:"errorByOperation"`
}

// AnalyzeTraces анализирует коллекцию трейсов
func (a *TraceAnalyzer) AnalyzeTraces(traces []*Trace) *TimeRangeAnalysis {
	if len(traces) == 0 {
		return &TimeRangeAnalysis{}
	}

	analysis := &TimeRangeAnalysis{
		ServiceStats: ServiceStats{
			ServiceCounts:    make(map[string]int),
			ServiceDurations: make(map[string]time.Duration),
		},
		ErrorStats: ErrorStats{
			ErrorByService:   make(map[string]int),
			ErrorByOperation: make(map[string]int),
		},
	}

	// Собираем временные границы
	analysis.collectTimeRangeStats(traces)

	// Анализируем сервисы и операции
	analysis.collectServiceAndOperationStats(traces)

	// Анализируем производительность
	analysis.collectPerformanceStats(traces)

	// Анализируем ошибки
	analysis.collectErrorStats(traces)

	return analysis
}

func (a *TimeRangeAnalysis) collectTimeRangeStats(traces []*Trace) {
	var startTime, endTime time.Time
	totalSpans := 0

	for i, trace := range traces {
		totalSpans += len(trace.Spans)

		if i == 0 || trace.StartTime.Before(startTime) {
			startTime = trace.StartTime
		}

		traceEnd := trace.StartTime.Add(trace.Duration)
		if i == 0 || traceEnd.After(endTime) {
			endTime = traceEnd
		}
	}

	timeRangeDuration := endTime.Sub(startTime)
	minutes := timeRangeDuration.Minutes()
	if minutes == 0 {
		minutes = 1
	}

	a.TimeRange = TimeRangeStats{
		StartTime:       startTime,
		EndTime:         endTime,
		Duration:        timeRangeDuration.String(),
		TotalTraces:     len(traces),
		TotalSpans:      totalSpans,
		TracesPerMinute: float64(len(traces)) / minutes,
		SpansPerMinute:  float64(totalSpans) / minutes,
	}
}

func (a *TimeRangeAnalysis) collectServiceAndOperationStats(traces []*Trace) {
	operationData := make(map[string]*operationAccumulator)
	serviceSet := make(map[string]bool)

	for _, trace := range traces {
		for _, span := range trace.Spans {
			// Собираем статистику по сервисам
			serviceSet[span.ServiceName] = true
			a.ServiceStats.ServiceCounts[span.ServiceName]++
			a.ServiceStats.ServiceDurations[span.ServiceName] += span.Duration

			// Собираем статистику по операциям
			key := span.ServiceName + ":" + span.OperationName
			if op, exists := operationData[key]; exists {
				op.count++
				op.totalDuration += span.Duration
				if span.Duration < op.minDuration {
					op.minDuration = span.Duration
				}
				if span.Duration > op.maxDuration {
					op.maxDuration = span.Duration
				}
				op.durations = append(op.durations, span.Duration)
			} else {
				operationData[key] = &operationAccumulator{
					serviceName:   span.ServiceName,
					operationName: span.OperationName,
					count:         1,
					totalDuration: span.Duration,
					minDuration:   span.Duration,
					maxDuration:   span.Duration,
					durations:     []time.Duration{span.Duration},
				}
			}
		}
	}

	// Преобразуем набор сервисов в слайс
	for service := range serviceSet {
		a.ServiceStats.Services = append(a.ServiceStats.Services, service)
	}

	// Преобразуем данные операций в статистику
	for _, op := range operationData {
		avgDuration := time.Duration(int64(op.totalDuration) / int64(op.count))
		p95 := a.calculatePercentile(op.durations, 95)

		a.OperationStats = append(a.OperationStats, OperationStats{
			OperationName: op.operationName,
			ServiceName:   op.serviceName,
			Count:         op.count,
			AvgDuration:   avgDuration,
			MinDuration:   op.minDuration,
			MaxDuration:   op.maxDuration,
			P95Duration:   p95,
			TotalDuration: op.totalDuration,
		})
	}
}

func (a *TimeRangeAnalysis) collectPerformanceStats(traces []*Trace) {
	var totalTraceDuration, totalSpanDuration time.Duration
	var minTraceDuration, maxTraceDuration time.Duration
	totalSpans := 0

	for i, trace := range traces {
		totalTraceDuration += trace.Duration
		totalSpans += len(trace.Spans)

		for _, span := range trace.Spans {
			totalSpanDuration += span.Duration
		}

		if i == 0 || trace.Duration < minTraceDuration {
			minTraceDuration = trace.Duration
		}
		if i == 0 || trace.Duration > maxTraceDuration {
			maxTraceDuration = trace.Duration
		}
	}

	a.PerformanceStats = PerformanceStats{
		AvgTraceDuration: time.Duration(int64(totalTraceDuration) / int64(len(traces))),
		MinTraceDuration: minTraceDuration,
		MaxTraceDuration: maxTraceDuration,
		AvgSpanDuration:  time.Duration(int64(totalSpanDuration) / int64(totalSpans)),
	}

	// Находим самые медленные операции
	a.PerformanceStats.SlowestOperations = a.findSlowestOperations(10)
}

func (a *TimeRangeAnalysis) collectErrorStats(traces []*Trace) {
	totalSpans := 0
	errorSpans := 0

	for _, trace := range traces {
		for _, span := range trace.Spans {
			totalSpans++

			// Проверяем теги на наличие ошибок
			if isErrorSpan(span) {
				errorSpans++
				a.ErrorStats.ErrorByService[span.ServiceName]++
				a.ErrorStats.ErrorByOperation[span.OperationName]++
			}
		}
	}

	a.ErrorStats.TotalErrors = errorSpans
	if totalSpans > 0 {
		a.ErrorStats.ErrorRate = float64(errorSpans) / float64(totalSpans) * 100
	}
}

func isErrorSpan(span Span) bool {
	// Проверяем различные признаки ошибок
	if status, exists := span.Tags["http.status_code"]; exists {
		if status >= "400" && status <= "599" {
			return true
		}
	}

	if errorTag, exists := span.Tags["error"]; exists {
		return errorTag == "true"
	}

	if span.Tags["span.status"] == "ERROR" {
		return true
	}

	return false
}

func (a *TimeRangeAnalysis) findSlowestOperations(limit int) []OperationStats {
	// Сортируем операции по убыванию среднего времени
	operations := make([]OperationStats, len(a.OperationStats))
	copy(operations, a.OperationStats)

	for i := 0; i < len(operations)-1; i++ {
		for j := i + 1; j < len(operations); j++ {
			if operations[j].AvgDuration > operations[i].AvgDuration {
				operations[i], operations[j] = operations[j], operations[i]
			}
		}
	}

	if len(operations) > limit {
		return operations[:limit]
	}
	return operations
}

func (a *TimeRangeAnalysis) calculatePercentile(durations []time.Duration, percentile int) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Сортируем durations
	for i := 0; i < len(durations)-1; i++ {
		for j := i + 1; j < len(durations); j++ {
			if durations[j] < durations[i] {
				durations[i], durations[j] = durations[j], durations[i]
			}
		}
	}

	index := (percentile * len(durations)) / 100
	if index >= len(durations) {
		index = len(durations) - 1
	}

	return durations[index]
}

type operationAccumulator struct {
	serviceName   string
	operationName string
	count         int
	totalDuration time.Duration
	minDuration   time.Duration
	maxDuration   time.Duration
	durations     []time.Duration
}

// PrintTimeRangeAnalysis выводит анализ временного отрезка в читаемом формате
func (a *TraceAnalyzer) PrintTimeRangeAnalysis(analysis *TimeRangeAnalysis) {
	fmt.Printf("=== Time Range Analysis ===\n")
	fmt.Printf("Time Range: %s - %s (%s)\n",
		analysis.TimeRange.StartTime.Format("2006-01-02 15:04:05"),
		analysis.TimeRange.EndTime.Format("2006-01-02 15:04:05"),
		analysis.TimeRange.Duration)
	fmt.Printf("Total Traces: %d, Total Spans: %d\n",
		analysis.TimeRange.TotalTraces, analysis.TimeRange.TotalSpans)
	fmt.Printf("Traces/Minute: %.2f, Spans/Minute: %.2f\n\n",
		analysis.TimeRange.TracesPerMinute, analysis.TimeRange.SpansPerMinute)

	fmt.Printf("=== Services ===\n")
	for _, service := range analysis.ServiceStats.Services {
		count := analysis.ServiceStats.ServiceCounts[service]
		avgDuration := analysis.ServiceStats.ServiceDurations[service] / time.Duration(count)
		fmt.Printf("- %s: %d spans, avg duration: %v\n", service, count, avgDuration)
	}
	fmt.Println()

	fmt.Printf("=== Performance ===\n")
	fmt.Printf("Avg Trace Duration: %v\n", analysis.PerformanceStats.AvgTraceDuration)
	fmt.Printf("Trace Duration Range: %v - %v\n",
		analysis.PerformanceStats.MinTraceDuration, analysis.PerformanceStats.MaxTraceDuration)
	fmt.Printf("Avg Span Duration: %v\n\n", analysis.PerformanceStats.AvgSpanDuration)

	fmt.Printf("=== Slowest Operations (Top 5) ===\n")
	for i, op := range analysis.PerformanceStats.SlowestOperations {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s.%s: %v (p95: %v)\n",
			i+1, op.ServiceName, op.OperationName, op.AvgDuration, op.P95Duration)
	}
	fmt.Println()

	fmt.Printf("=== Errors ===\n")
	fmt.Printf("Total Errors: %d, Error Rate: %.2f%%\n",
		analysis.ErrorStats.TotalErrors, analysis.ErrorStats.ErrorRate)
	for service, count := range analysis.ErrorStats.ErrorByService {
		fmt.Printf("- %s: %d errors\n", service, count)
	}
}
