package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ProcessingTimes struct {
	timers sync.Map
}

func NewTimer() *ProcessingTimes {
	return &ProcessingTimes{}
}

type Stopwatch struct {
	Start    time.Time
	Duration *time.Duration
}

func (p *ProcessingTimes) Start(name string) *Stopwatch {
	stopwatch := &Stopwatch{
		Start: time.Now().UTC(),
	}
	p.timers.Store(name, stopwatch)
	return stopwatch
}

func (s *Stopwatch) Stop() {
	duration := time.Since(s.Start)
	s.Duration = &duration
}

type Field struct {
	Name  string
	Value interface{}
}

func NewField(name string, value interface{}) Field {
	return Field{
		Name:  name,
		Value: value,
	}
}

func (p *ProcessingTimes) ToLogField() Field {
	fields := map[string]int64{}
	p.timers.Range(func(key, value interface{}) bool {
		stopwatch := value.(*Stopwatch)
		if stopwatch.Duration == nil {
			return true
		}

		fields[key.(string)] = int64(float64(*stopwatch.Duration) / float64(time.Millisecond))
		return true
	})

	return NewField("processing-times", fields)
}

func main() {
	ctx := context.Background()
	timer := NewTimer()
	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		ex := timer.Start("brasilApi")
		brasilApi, err := ExternalIntegration(ctx, http.MethodGet, "https://brasilapi.com.br/api/cep/v1/01153000")
		ex.Stop()
		if err != nil {
			fmt.Println("brasil api error: ", err)
		}
		fmt.Println("brasilApi", brasilApi)
		wait.Done()
	}()
	go func() {
		ex2 := timer.Start("viaCep")
		viaCep, err := ExternalIntegration(ctx, http.MethodPost, "http://viacep.com.br/ws/01153000/json/")
		ex2.Stop()
		if err != nil {
			fmt.Println("via cep error: ", err)
		}
		fmt.Println("viaCep", viaCep)
		wait.Done()
	}()
	wait.Wait()
	fmt.Println("processing Times: ", timer.ToLogField())
}

func ExternalIntegration(ctx context.Context, method string, url string) (string, error) {
	var address string
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		fmt.Println("new request error: ", err)
		return address, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("http client error: ", err)
		return address, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read body error: ", err)
		return address, err
	}

	// err = json.Unmarshal(body, &address)
	// if err != nil {
	// 	fmt.Println("unmarshal error: ", err)
	// 	return address, err
	// }
	address = string(body)

	return address, nil
}
