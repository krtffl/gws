package webui

import (
	"io"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/oxtoacart/bpool"

	"github.com/krtffl/get-well-soon/internal/domain"
	"github.com/krtffl/get-well-soon/internal/logger"
)

type Handler struct {
	svc      *Service
	template *template.Template
	bpool    *bpool.BufferPool
}

func NewHandler(svc *Service, bpool *bpool.BufferPool) *Handler {
	tmpls, err := template.New("").ParseGlob("public/templates/*.html")
	if err != nil {
		logger.Fatal("[WebuiHandler - Content] - Failed to parse templates. %v", err)
	}

	return &Handler{template: tmpls, svc: svc, bpool: bpool}
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
	logger.Info("[WebuiHandler - Content - Challenge] Incoming request")

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	if err := h.template.ExecuteTemplate(buf, "challenge.html", nil); err != nil {
		logger.Error("[WebuiHandler - Content - Challenge] Couldn't execute template. %v", err)
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

	_, err = h.svc.Create(gws)
	if err != nil {
		logger.Error("[WebuiHandler - Content - Upload] Couldn't create gws. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf := h.bpool.Get()
	defer h.bpool.Put(buf)

	if err := h.template.ExecuteTemplate(buf, "uploaded.html", nil); err != nil {
		logger.Error("[WebuiHandler - Content - Upload] Couldn't execute template. %v", err)
		h.template.ExecuteTemplate(w, "error.html", nil)
		return
	}

	buf.WriteTo(w)
	return
}
