package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"qin-culture-site/internal/domain"
	"qin-culture-site/internal/service"
)

type Server struct {
	service *service.Service
	tpl     *template.Template
	maxBody int64
}

func NewServer(svc *service.Service, maxBody int64) (*Server, error) {
	if svc == nil {
		return nil, fmt.Errorf("service is required")
	}
	if maxBody <= 0 {
		maxBody = 1 << 20
	}
	tpl, err := template.New("home").Funcs(template.FuncMap{"safe": domain.SafeText, "join": service.JoinLabels}).Parse(homeTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse home template: %w", err)
	}
	return &Server{service: svc, tpl: tpl, maxBody: maxBody}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/api/songs", s.handleSongs)
	mux.HandleFunc("/api/study", s.handleStudy)
	mux.HandleFunc("/api/report", s.handleReport)
	mux.HandleFunc("/api/experience", s.handleExperience)
	mux.HandleFunc("/healthz", s.handleHealth)
	return loggingMiddleware(mux)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	model, err := s.service.Browse(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := struct {
		Title string
		Nav   []service.NavigationItem
		Home  service.HomeModel
	}{Title: "听见山水：古琴艺术专题", Nav: service.Navigation(), Home: model}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tpl.Execute(w, page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "只读接口")
		return
	}
	writeJSON(w, http.StatusOK, s.service.Catalog().Pieces())
}

func (s *Server) handleStudy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "只读接口")
		return
	}
	model, err := s.service.Study(r.Context(), service.StudyRequest{PieceID: r.URL.Query().Get("piece"), Query: r.URL.Query().Get("q"), Difficulty: r.URL.Query().Get("difficulty")})
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model)
}

func (s *Server) handleExperience(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSONError(w, http.StatusMethodNotAllowed, "需要 POST")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBody)
	input, err := decodeExperience(r)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := s.service.SubmitExperience(r.Context(), input)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": item.ID, "status": item.Status, "message": "已收到你的听琴体验"})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "只读接口")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func decodeExperience(r *http.Request) (domain.ExperienceSubmission, error) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		var input domain.ExperienceSubmission
		decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
		if err := decoder.Decode(&input); err != nil {
			return domain.ExperienceSubmission{}, fmt.Errorf("invalid JSON: %w", err)
		}
		return input, nil
	}
	if err := r.ParseForm(); err != nil {
		return domain.ExperienceSubmission{}, fmt.Errorf("invalid form: %w", err)
	}
	return domain.ExperienceSubmission{Name: r.FormValue("name"), Contact: r.FormValue("contact"), Interest: r.FormValue("interest"), Message: r.FormValue("message")}, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("X-Qin-Route", strconv.Itoa(len(r.URL.Path)))
		}
		next.ServeHTTP(w, r)
	})
}
