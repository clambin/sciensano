package sciensano

import (
	"testing"
	"time"
)

const benchmarkSortCount = 10000

func BenchmarkCopyMaybeSort_NoSort(b *testing.B) {
	var input APIVaccinationsResponses
	timestamp := time.Now()
	for i := 0; i < benchmarkSortCount; i++ {
		input = append(input, &APIVaccinationsResponse{
			TimeStamp: TimeStamp{Time: timestamp},
		})
		timestamp = timestamp.Add(time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copyMaybeSort(input)
	}
}

func BenchmarkCopyMaybeSort_Sort(b *testing.B) {
	var input APIVaccinationsResponses
	timestamp := time.Now()
	for i := 0; i < benchmarkSortCount; i++ {
		input = append(input, &APIVaccinationsResponse{
			TimeStamp: TimeStamp{Time: timestamp},
		})
		timestamp = timestamp.Add(-time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copyMaybeSort(input)
	}
}
