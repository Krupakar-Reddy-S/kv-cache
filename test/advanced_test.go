package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"runtime"
	"github.com/shirou/gopsutil/v3/cpu"
)

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

type TestStats struct {
	totalRequests      int64
	successfulPuts     int64
	successfulGets     int64
	failedPuts        int64
	failedGets        int64
	cacheMisses       int64 // in microseconds
	totalLatency      int64 // in microseconds
	maxLatency        int64 // in microseconds
	putLatencyTotal   int64
	getLatencyTotal   int64
	putCount          int64
	getCount          int64
	maxMemoryUsed     uint64 // in MB
	currentMemoryMB   uint64
	cpuUsage          float64
}

const (
	numWorkers    = 50      // Number of concurrent workers
	testDuration  = 300     // Test duration in seconds
	reportInterval = 5      // Status report interval in seconds
	maxKeyLength   = 200    // Maximum key length to test
	maxValueLength = 200    // Maximum value length to test
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Add a thread-safe key store
type KeyStore struct {
	sync.RWMutex
	keys []string
}

func (ks *KeyStore) Add(key string) {
	ks.Lock()
	ks.keys = append(ks.keys, key)
	ks.Unlock()
}

func (ks *KeyStore) GetRandom() string {
	ks.RLock()
	defer ks.RUnlock()
	if len(ks.keys) == 0 {
		return ""
	}
	return ks.keys[rand.Intn(len(ks.keys))]
}

func getSystemMetrics() (memoryMB uint64, cpuPercent float64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	memoryMB = m.Alloc / 1024 / 1024

	// Get CPU usage
	var cpuStats cpu.TimesStat
	if stats, err := cpu.Times(false); err == nil && len(stats) > 0 {
		cpuStats = stats[0]
		cpuPercent = cpuStats.User + cpuStats.System
	}
	
	return memoryMB, cpuPercent
}

func TestAdvancedLoadTest(t *testing.T) {
	baseURL := "http://localhost:7171"
	maxRetries := 5
	var err error

	t.Log("Checking if service is running...")
	for i := 0; i < maxRetries; i++ {
		_, err = http.Get(baseURL + "/get?key=test")
		if err == nil {
			break
		}
		t.Logf("Service not ready, retrying in 2 seconds (attempt %d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		t.Fatalf("Service is not running on %s", baseURL)
	}

	t.Log("Service is running! Starting load test...")
	t.Logf("Test configuration:")
	t.Logf("- Number of concurrent workers: %d", numWorkers)
	t.Logf("- Test duration: %d seconds", testDuration)
	t.Logf("- Report interval: %d seconds", reportInterval)
	t.Logf("- Maximum key length: %d", maxKeyLength)
	t.Logf("- Maximum value length: %d", maxValueLength)
	
	rand.Seed(time.Now().UnixNano())
	runLoadTest(t, baseURL)
}

func runLoadTest(t *testing.T, baseURL string) {
	stats := &TestStats{}
	keyMap := sync.Map{}
	keyStore := &KeyStore{keys: make([]string, 0, 1000)}
	start := time.Now()
	var wg sync.WaitGroup

	// Create channels for coordination
	statsDone := make(chan bool)
	workersDone := make(chan bool)

	// Start statistics reporter
	go reportStats(t, stats, start, statsDone)

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(t, i, baseURL, stats, &keyMap, keyStore, &wg, workersDone)
	}

	// Wait for test duration
	time.Sleep(time.Duration(testDuration) * time.Second)
	
	// Signal all workers to stop
	close(workersDone)
	
	// Wait for all workers to finish
	wg.Wait()
	
	// Stop the stats reporter
	statsDone <- true

	// Final report
	printFinalReport(t, stats, start)
}

func worker(t *testing.T, id int, baseURL string, stats *TestStats, keyMap *sync.Map, keyStore *KeyStore, wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	for {
		select {
		case <-done:
			return
		default:
			// Randomly choose between PUT and GET (60% PUT, 40% GET)
			if rand.Float32() < 0.6 {
				doPut(t, client, baseURL, stats, keyMap, keyStore)
			} else {
				doGet(t, client, baseURL, stats, keyMap, keyStore)
			}
		}
	}
}

func doPut(t *testing.T, client *http.Client, baseURL string, stats *TestStats, keyMap *sync.Map, keyStore *KeyStore) {
	key := randomString(rand.Intn(maxKeyLength) + 1)
	value := randomString(rand.Intn(maxValueLength) + 1)
	req := PutRequest{Key: key, Value: value}
	jsonData, _ := json.Marshal(req)

	start := time.Now()
	resp, err := client.Post(baseURL+"/put", "application/json", bytes.NewBuffer(jsonData))
	latency := time.Since(start).Microseconds()

	atomic.AddInt64(&stats.totalRequests, 1)
	atomic.AddInt64(&stats.putCount, 1)
	atomic.AddInt64(&stats.putLatencyTotal, latency)
	atomic.AddInt64(&stats.totalLatency, latency)
	
	if latency > atomic.LoadInt64(&stats.maxLatency) {
		atomic.StoreInt64(&stats.maxLatency, latency)
	}

	if err != nil {
		atomic.AddInt64(&stats.failedPuts, 1)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		atomic.AddInt64(&stats.successfulPuts, 1)
		keyMap.Store(key, value)
		keyStore.Add(key)  // Add successful key to our store
	} else {
		atomic.AddInt64(&stats.failedPuts, 1)
	}
}

func doGet(t *testing.T, client *http.Client, baseURL string, stats *TestStats, keyMap *sync.Map, keyStore *KeyStore) {
	// Try to get a stored key 80% of the time
	var key string
	if rand.Float32() < 0.8 {
		key = keyStore.GetRandom()
	}

	// If no stored keys yet or 20% of the time, use a random key
	if key == "" {
		key = randomString(rand.Intn(maxKeyLength) + 1)
	}

	start := time.Now()
	resp, err := client.Get(fmt.Sprintf("%s/get?key=%s", baseURL, key))
	latency := time.Since(start).Microseconds()

	atomic.AddInt64(&stats.totalRequests, 1)
	atomic.AddInt64(&stats.getCount, 1)
	atomic.AddInt64(&stats.getLatencyTotal, latency)
	atomic.AddInt64(&stats.totalLatency, latency)

	if latency > atomic.LoadInt64(&stats.maxLatency) {
		atomic.StoreInt64(&stats.maxLatency, latency)
	}

	if err != nil {
		atomic.AddInt64(&stats.failedGets, 1)
		return
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		atomic.AddInt64(&stats.failedGets, 1)
		return
	}

	if resp.StatusCode == http.StatusOK {
		atomic.AddInt64(&stats.successfulGets, 1)
		// Verify if the value matches (for known keys)
		if expectedValue, ok := keyMap.Load(key); ok && response.Value != expectedValue.(string) {
			atomic.AddInt64(&stats.cacheMisses, 1)
		}
	} else if resp.StatusCode == http.StatusNotFound {
		if _, exists := keyMap.Load(key); exists {
			atomic.AddInt64(&stats.cacheMisses, 1)
		}
	} else {
		atomic.AddInt64(&stats.failedGets, 1)
	}
}

func reportStats(t *testing.T, stats *TestStats, start time.Time, done chan bool) {
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			memMB, cpuPercent := getSystemMetrics()
			atomic.StoreUint64(&stats.currentMemoryMB, memMB)
			if memMB > atomic.LoadUint64(&stats.maxMemoryUsed) {
				atomic.StoreUint64(&stats.maxMemoryUsed, memMB)
			}
			stats.cpuUsage = cpuPercent
			printStats(t, stats, start)
		case <-done:
			return
		}
	}
}

func printStats(t *testing.T, stats *TestStats, start time.Time) {
	elapsed := time.Since(start).Seconds()
	puts := atomic.LoadInt64(&stats.putCount)
	gets := atomic.LoadInt64(&stats.getCount)
	
	t.Logf("\n=== Performance Report (%.1f seconds elapsed) ===", elapsed)
	t.Logf("Total Requests: %d (%.1f req/sec)", 
		atomic.LoadInt64(&stats.totalRequests),
		float64(atomic.LoadInt64(&stats.totalRequests))/elapsed)
	t.Logf("Successful PUTs: %d, Failed: %d", 
		atomic.LoadInt64(&stats.successfulPuts),
		atomic.LoadInt64(&stats.failedPuts))
	t.Logf("Successful GETs: %d, Failed: %d",
		atomic.LoadInt64(&stats.successfulGets),
		atomic.LoadInt64(&stats.failedGets))
	t.Logf("Cache Misses: %d", atomic.LoadInt64(&stats.cacheMisses))
	
	if puts > 0 {
		t.Logf("Average PUT latency: %.2f ms", 
			float64(atomic.LoadInt64(&stats.putLatencyTotal))/float64(puts)/1000)
	}
	if gets > 0 {
		t.Logf("Average GET latency: %.2f ms",
			float64(atomic.LoadInt64(&stats.getLatencyTotal))/float64(gets)/1000)
	}
	t.Logf("Max latency: %.2f ms", float64(atomic.LoadInt64(&stats.maxLatency))/1000)
	t.Logf("Current Memory Usage: %d MB", atomic.LoadUint64(&stats.currentMemoryMB))
	t.Logf("Max Memory Usage: %d MB", atomic.LoadUint64(&stats.maxMemoryUsed))
	t.Logf("CPU Usage: %.2f%%", stats.cpuUsage)
}

func printFinalReport(t *testing.T, stats *TestStats, start time.Time) {
	elapsed := time.Since(start).Seconds()
	totalReqs := atomic.LoadInt64(&stats.totalRequests)
	
	t.Logf("\n=== Final Test Report ===")
	t.Logf("Test Duration: %.1f seconds", elapsed)
	t.Logf("Total Requests: %d", totalReqs)
	t.Logf("Average Throughput: %.1f requests/second", float64(totalReqs)/elapsed)
	t.Logf("Successful Operations: %d (PUTs: %d, GETs: %d)",
		atomic.LoadInt64(&stats.successfulPuts)+atomic.LoadInt64(&stats.successfulGets),
		atomic.LoadInt64(&stats.successfulPuts),
		atomic.LoadInt64(&stats.successfulGets))
	t.Logf("Failed Operations: %d (PUTs: %d, GETs: %d)",
		atomic.LoadInt64(&stats.failedPuts)+atomic.LoadInt64(&stats.failedGets),
		atomic.LoadInt64(&stats.failedPuts),
		atomic.LoadInt64(&stats.failedGets))
	t.Logf("Cache Misses: %d", atomic.LoadInt64(&stats.cacheMisses))
	t.Logf("Average Latency: %.2f ms", 
		float64(atomic.LoadInt64(&stats.totalLatency))/float64(totalReqs)/1000)
	t.Logf("Max Latency: %.2f ms", 
		float64(atomic.LoadInt64(&stats.maxLatency))/1000)
	t.Logf("Peak Memory Usage: %d MB", atomic.LoadUint64(&stats.maxMemoryUsed))
	t.Logf("Final CPU Usage: %.2f%%", stats.cpuUsage)
} 