package webui

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/oxtoacart/bpool"

	getwellsoon "github.com/krtffl/gws"
	"github.com/krtffl/gws/internal/cookie"
	"github.com/krtffl/gws/internal/domain"
	"github.com/krtffl/gws/internal/logger"
)

type Memories struct {
	GWSs []*Content
}

type Content struct {
	From    string    `json:"from"`
	Message string    `json:"message"`
	Memory  string    `json:"memory"`
	Date    time.Time `json:"date"`
	IsLast  bool      `json:"isLast"`
	Offset  int       `json:"offset"`
}

type Handler struct {
	svc       *Service
	template  *template.Template
	bpool     *bpool.BufferPool
	cookie    *cookie.Service
	challenge []string
}

func NewHandler(
	svc *Service,
	bpool *bpool.BufferPool,
	cookie *cookie.Service,
	challenge []string,
) *Handler {
	tmpls, err := template.New("").ParseFS(getwellsoon.Public, "public/templates/*.html")
	if err != nil {
		logger.Fatal("[WebuiHandler - Content] - Failed to parse templates. %v", err)
	}

	return &Handler{template: tmpls, svc: svc, bpool: bpool, cookie: cookie, challenge: challenge}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	logger.Info("[WebuiHandler - Content - Index] Incoming request")

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	if err := h.template.ExecuteTemplate(buf, "index.html", nil); err != nil {
		logger.Error("[WebuiHandler - Content - Index] Couldn't execute template. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf.WriteTo(w)
	return
}

func (h *Handler) Form(w http.ResponseWriter, r *http.Request) {
	logger.Info("[WebuiHandler - Content - Form] Incoming request")

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	if err := h.template.ExecuteTemplate(buf, "form.html", nil); err != nil {
		logger.Error("[WebuiHandler - Content - Form] Couldn't execute template. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf.WriteTo(w)
	return
}

func (h *Handler) Challenge(w http.ResponseWriter, r *http.Request) {
	logger.Info("[WebuiHandler - Content - GetChallenge] Incoming request")

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	if err := h.template.ExecuteTemplate(buf, "challenge.html", nil); err != nil {
		logger.Error("[WebuiHandler - Content - GetChallenge] Couldn't execute template. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf.WriteTo(w)
	return
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	logger.Info("[WebuiHandler - Content - Upload] Incoming request")

	from := r.FormValue("from")
	message := r.FormValue("message")

	gws := &domain.GWS{}
	gws.From = from
	gws.Message = message
	gws.Id = uuid.NewString()

	memory, _, err := r.FormFile("memory")
	switch {
	case err == http.ErrMissingFile:
		logger.Info("[WebuiHandler - Content - Upload] GWS without memory.")
	case err != nil:
		logger.Error("[WebuiHandler - Content - Upload] Error fetching memory. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	default:
		m, err := io.ReadAll(memory)
		if err != nil {
			logger.Error("[WebuiHandler - Content - Upload] Couldn't read memory. %v", err)
			h.template.ExecuteTemplate(w, "error.html", nil)
			return
		}
		gws.Memory = m
	}

	uploadedGWS, err := h.svc.Create(gws)
	if err != nil {
		logger.Error("[WebuiHandler - Content - Upload] Couldn't create gws. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	m := ""
	if len(uploadedGWS.Memory) > 0 {
		m = base64.RawStdEncoding.EncodeToString(uploadedGWS.Memory)
	}

	if err := h.template.ExecuteTemplate(buf, "uploaded.html", Content{
		From:    uploadedGWS.From,
		Message: uploadedGWS.Message,
		Date:    uploadedGWS.CreatedAt,
		Memory:  m,
	}); err != nil {
		logger.Error("[WebuiHandler - Content - Upload] Couldn't execute template. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf.WriteTo(w)
	return
}

func (h *Handler) Solve(w http.ResponseWriter, r *http.Request) {
	logger.Info("[WebuiHandler - Content - SolveChallenge] Incoming request")

	challenge := r.FormValue("challenge")

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	if !hasSolved(challenge, h.challenge) {
		logger.Info("[WebuiHandler - Content - SolveChallenge] Failed challenge: %s", challenge)
		if err := h.template.ExecuteTemplate(buf, "failed.html", nil); err != nil {
			logger.Error(
				"[WebuiHandler - Content - SolveChallenge] Couldn't execute template. %v",
				err,
			)
			h.template.ExecuteTemplate(w, "failed.html", nil)
			return
		}

		buf.WriteTo(w)
		return
	}

	logger.Info("[WebuiHandler - Content - SolveChallenge] Solved challenge: %s", challenge)
	c, err := h.cookie.CreateCookie(challenge)
	if err != nil {
		logger.Error(
			"[WebuiHandler - Content - SolveChallenge] Couldn't create cookie. %v",
			err,
		)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "/secure/memories", http.StatusFound)
	return
}

func (h *Handler) Memories(w http.ResponseWriter, r *http.Request) {
	logger.Info("[WebuiHandler - Content - Memories] Incoming request")

	challenge, ok := r.Context().Value("challenge").(string)
	if !ok {
		logger.Error(
			"[WebuiHandler - Content - Memories] Couldn't retrieve challenge from context",
		)
		http.Redirect(w, r, "/challenge", http.StatusFound)
		return
	}

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	if !hasSolved(challenge, h.challenge) {
		logger.Info("[WebuiHandler - Content - Memories] Failed challenge: %s", challenge)
		http.Redirect(w, r, "/challenge", http.StatusFound)
		return
	}

	offsetString := r.URL.Query().Get("offset")
	offset := 0

	if offsetString != "" {
		off, err := strconv.Atoi(offsetString)
		if err != nil {
			logger.Error(
				"[WebuiHandler - Content - Memories] Couldn't parse offset. %v",
				err,
			)
			h.template.ExecuteTemplate(w, "error.html", nil)
			return
		}
		offset = off
	}

	gwss, err := h.svc.List(5, uint(offset))
	if err != nil {
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	var memories []*Content
	for i, gws := range gwss {
		c := &Content{
			From:    gws.From,
			Message: gws.Message,
			Date:    gws.CreatedAt,
		}

		m := ""
		if len(gws.Memory) > 0 {
			m = base64.RawStdEncoding.EncodeToString(gws.Memory)
		}

		if i == 4 || i == len(gwss)-1 {
			c.IsLast = true
			c.Offset = offset + 5
		}

		c.Memory = m
		memories = append(memories, c)
	}

	if err := h.template.ExecuteTemplate(buf, "memories.html", Memories{
		GWSs: memories,
	}); err != nil {
		logger.Error(
			"[WebuiHandler - Content - SolveChallenge] Couldn't execute template. %v",
			err,
		)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf.WriteTo(w)
	return
}

func hasSolved(send string, challenge []string) bool {
	for _, ch := range challenge {
		if send == ch {
			return true
		}
	}

	return false
}
