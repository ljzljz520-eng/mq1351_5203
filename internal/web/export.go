package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"qin-culture-site/internal/service"
)

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "只读接口")
		return
	}
	report, err := s.service.Report(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if r.URL.Query().Get("format") == "text" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(strings.Join(report.Lines(), "\n")))
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func reportLink(query string) string {
	return "/api/report?format=text&q=" + url.QueryEscape(query)
}

func reportSummary(report service.CultureReport) string {
	if report.SearchTerm == "" {
		return fmt.Sprintf("专题资料 %d 条", report.Stats.PieceCount+report.Stats.StoryCount)
	}
	return fmt.Sprintf("“%s”命中 %d 条", report.SearchTerm, report.SearchHits)
}
