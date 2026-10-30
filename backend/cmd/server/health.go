package main

import "net/http"

// health handles GET /api/health.
//
//declscope:shared // main.go routes to it
func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	if _, err := h.client.HealthWithContext(r.Context()); err != nil {
		writeErrorResponse(w, http.StatusServiceUnavailable, "meilisearch is not healthy")
		return
	}
	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}
