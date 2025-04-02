package main

import (
	"fmt"
	"net/http"
)

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	/* json.NewEncoder(w).Encode(map[string]int32{
		"fileserve_hits": c.fileserveHits.Load(),
	}) */
	//hits := c.fileserveHits.Load()
	w.Write([]byte(fmt.Sprintf(`
	<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>

	</html>
	`, c.fileserveHits.Load())))
	/* _, err := fmt.Fprintf(w, "<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", hits)
	if err != nil {
		log.Printf("Error writing metrics: %v", err)
	} */
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	//fileserveHits = atomic.Int32.Add(1)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserveHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
