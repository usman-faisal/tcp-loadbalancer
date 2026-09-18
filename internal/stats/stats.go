package stats

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"usman-faisal/tcp-loadbalancer/internal/types"
)

func fetchStats(backend string) *types.Stats {
	statsRoute := "/stats"
	url := backend + statsRoute

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error occurred while fetching stats %e", err)
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Server returned status: %d\n", res.StatusCode)
		return nil
	}

	var stats types.Stats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	return &stats
}

func calcScore(stats *types.Stats) float64 {
	if stats == nil {
		return math.MaxFloat64
	}

	memPercent := float64(0)
	if stats.MemTotalBytes > 0 {
		memPercent = float64(stats.MemUsedBytes) / float64(stats.MemTotalBytes) * 100
	}

	score := (stats.CPUPercent * 0.5) +
		(memPercent * 0.3) +
		(float64(stats.ActiveConns) * 0.2)

	return score
}
