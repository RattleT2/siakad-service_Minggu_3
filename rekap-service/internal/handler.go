package internal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.GET("/rekap/jurusan", h.PerJurusan)
	r.GET("/rekap/top-ipk", h.TopIPK)
	r.GET("/rekap/mahasiswa/:nim", h.RingkasanMahasiswa)
}

func (h *Handler) PerJurusan(c *gin.Context) {
	data, err := h.svc.PerJurusan(c.Request.Context())
	if err != nil {
		handleRekapError(c, err)
		return
	}
	Success(c, http.StatusOK, data)
}

func (h *Handler) TopIPK(c *gin.Context) {
	n, err := strconv.Atoi(c.DefaultQuery("n", "3"))
	if err != nil || n <= 0 {
		Error(c, http.StatusBadRequest, nil)
		return
	}

	data, err := h.svc.TopIPK(c.Request.Context(), n)
	if err != nil {
		handleRekapError(c, err)
		return
	}
	Success(c, http.StatusOK, data)
}

func (h *Handler) RingkasanMahasiswa(c *gin.Context) {
	nim := c.Param("nim")
	data, err := h.svc.RingkasanMahasiswa(c.Request.Context(), nim)
	if err != nil {
		handleRekapError(c, err)
		return
	}
	Success(c, http.StatusOK, data)
}

func handleRekapError(c *gin.Context, err error) {
	status := MapErrorToHTTP(err)
	Error(c, status, err)
}
