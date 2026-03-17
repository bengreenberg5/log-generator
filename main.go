package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"
)

const lineLen = 256 // bytes per log line including newline
const poolSize = 50

var paths = []string{
	"/api/v1/query",
	"/api/v1/push",
	"/api/v1/labels",
	"/api/v1/series",
	"/api/v1/tail",
	"/api/v2/alerts",
	"/healthz",
	"/ready",
	"/metrics",
	"/api/v1/rules",
}

var levels = []string{"INFO", "WARN", "DEBUG", "ERROR"}

var hosts = []string{
	"gen-01", "gen-02", "gen-03", "gen-04", "gen-05",
}

type logTemplate struct {
	level     string
	host      string
	rid       string
	latencyMs int
	bytesIn   int
	path      string
}

func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func randHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generatePool() [poolSize]logTemplate {
	var pool [poolSize]logTemplate
	for i := range pool {
		pool[i] = logTemplate{
			level:     levels[randInt(len(levels))],
			host:      hosts[randInt(len(hosts))],
			rid:       randHex(8),
			latencyMs: randInt(5000),
			bytesIn:   randInt(65536),
			path:      paths[randInt(len(paths))],
		}
	}
	return pool
}

// generateRandomPool pre-generates random lines (each lineLen bytes including newline)
func generateRandomPool() [poolSize][]byte {
	var pool [poolSize][]byte
	for i := range pool {
		buf := make([]byte, lineLen)
		rand.Read(buf)
		// Make it printable hex, but we need exactly lineLen bytes.
		// Generate (lineLen-1)/2 random bytes, hex encode, take lineLen-1 chars, append newline.
		raw := make([]byte, lineLen) // more than enough
		rand.Read(raw)
		hexStr := hex.EncodeToString(raw)
		copy(buf, hexStr[:lineLen-1])
		buf[lineLen-1] = '\n'
		pool[i] = buf
	}
	return pool
}

func main() {
	rateKiB := 5
	mode := "structured"

	usage := func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [rate_KiB_per_sec] [structured|random]\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		v, err := strconv.Atoi(os.Args[1])
		if err != nil {
			usage()
		}
		rateKiB = v
	}
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "structured", "random":
			mode = os.Args[2]
		default:
			usage()
		}
	}

	bytesPerSec := rateKiB * 1024
	linesPerSec := bytesPerSec / lineLen
	if linesPerSec < 1 {
		linesPerSec = 1
	}

	fmt.Fprintf(os.Stderr, "Starting log generator at ~%dKiB/s (%d lines/s, %d bytes/line, mode=%s)\n", rateKiB, linesPerSec, lineLen, mode)

	// Batch writes for high throughput: write batchSize lines per tick
	batchSize := linesPerSec / 10
	if batchSize < 1 {
		batchSize = 1
	}
	ticksPerSec := linesPerSec / batchSize
	if ticksPerSec < 1 {
		ticksPerSec = 1
	}
	interval := time.Second / time.Duration(ticksPerSec)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	batchBuf := make([]byte, 0, batchSize*lineLen)

	if mode == "random" {
		randomPool := generateRandomPool()
		idx := 0
		for range ticker.C {
			batchBuf = batchBuf[:0]
			for b := 0; b < batchSize; b++ {
				batchBuf = append(batchBuf, randomPool[idx%poolSize]...)
				idx++
			}
			os.Stdout.Write(batchBuf)
		}
	} else {
		structPool := generatePool()
		idx := 0
		for range ticker.C {
			batchBuf = batchBuf[:0]
			for b := 0; b < batchSize; b++ {
				idx++
				t := structPool[idx%poolSize]
				ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

				prefix := fmt.Sprintf(
					`{"ts":"%s","level":"%s","i":%d,"host":"%s","rid":"%s","latency_ms":%d,"bytes_in":%d,"path":"%s","pad":"`,
					ts, t.level, idx, t.host, t.rid, t.latencyMs, t.bytesIn, t.path,
				)

				padLen := lineLen - len(prefix) - 3
				if padLen < 0 {
					padLen = 0
				}

				line := make([]byte, lineLen)
				copy(line, prefix)
				for j := len(prefix); j < len(prefix)+padLen; j++ {
					line[j] = 'x'
				}
				copy(line[len(prefix)+padLen:], "\"}\n")

				batchBuf = append(batchBuf, line...)
			}
			os.Stdout.Write(batchBuf)
		}
	}
}
