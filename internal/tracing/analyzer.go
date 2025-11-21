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
