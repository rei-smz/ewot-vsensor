package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
)

var (
	port     = flag.String("server.port", "8080", "service port")
	repo     = flag.String("server.repository", "", "SPARQL UPDATE endpoint URL (e.g. http://localhost:3030/ds/update)")
	descDir  = flag.String("server.description", "./", "directory with TTL files")
	dataPath = "/house/temperature/data"
	apPath   = "/house/temperature/data/access_point"

	rdfType, _  = rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	wotThing, _ = rdf.NewIRI("http://iot.linkeddata.es/def/wot#Thing")
	thingRe     = regexp.MustCompile(`<([^>]+)>\s+a\s+<http://iot.linkeddata.es/def/wot#Thing>`)
)

type sensorState struct {
	mu          sync.RWMutex
	Temperature float64 `json:"temperature"`
	Timestamp   int64   `json:"timestamp"`
}

func (s *sensorState) update() {
	s.mu.Lock()
	s.Temperature = -10.0 + rand.Float64()*40.0
	s.Timestamp = time.Now().UnixMilli()
	s.mu.Unlock()
}
func (s *sensorState) snapshot() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]any{
		"temperature": s.Temperature,
		"timestamp":   s.Timestamp,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	addr := ":" + *port

	if err := register(*port, *descDir, *repo); err != nil {
		log.Printf("Registration failed: %v", err)
		return
	} else {
		log.Printf("Registered device at %s", *repo)
	}

	var st sensorState
	st.update()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			st.update()
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc(dataPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(st.snapshot())
	})
	mux.HandleFunc(apPath, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, dataPath, http.StatusTemporaryRedirect)
	})

	log.Printf("Sensor running at http://127.0.0.1%s%s", addr, dataPath)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func register(port, dir, endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("no repository endpoint provided")
	}

	thingFile := filepath.Join(dir, fmt.Sprintf("thing-%s-temperature.ttl", port))
	descFile := filepath.Join(dir, fmt.Sprintf("description-%s-temperature.ttl", port))

	thingTTL, err := os.Open(thingFile)
	if err != nil {
		return fmt.Errorf("read thing file: %w", err)
	}
	defer thingTTL.Close()
	descTTL, err := os.Open(descFile)
	if err != nil {
		return fmt.Errorf("read description file: %w", err)
	}
	defer descTTL.Close()

	tris1, err := parseTTL(thingTTL)
	if err != nil {
		return fmt.Errorf("parse thing: %w", err)
	}
	tris2, err := parseTTL(descTTL)
	if err != nil {
		return fmt.Errorf("parse desc: %w", err)
	}
	all := append(tris1, tris2...)

	graphIRI, err := findThingIRI(tris1)
	if err != nil {
		return fmt.Errorf("find thing iri: %w", err)
	}

	var buf bytes.Buffer
	for _, t := range all {
		buf.WriteString(t.Serialize(rdf.NTriples))
	}

	query := fmt.Sprintf(`INSERT DATA { GRAPH <%s> { %s } }`, graphIRI, buf.String())

	repo, err := sparql.NewRepo(endpoint)
	if err != nil {
		return err
	}

	if err := repo.Update(query); err != nil {
		return fmt.Errorf("SPARQL update: %w", err)
	}
	return nil
}

func parseTTL(f *os.File) ([]rdf.Triple, error) {
	decoder := rdf.NewTripleDecoder(f, rdf.Turtle)
	triples, err := decoder.DecodeAll()
	if err != nil {
		return nil, err
	}
	return triples, nil
}

func findThingIRI(tris []rdf.Triple) (string, error) {
	for _, t := range tris {

		if t.Pred.String() == rdfType.String() && t.Obj.String() == wotThing.String() {
			return t.Subj.String(), nil
		}
	}

	return "", fmt.Errorf("no wot:Thing found")
}
