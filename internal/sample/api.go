package sample

import (
	"net/http"

	"github.com/pgillich/micro-server/pkg/logger"
)

type ApiServer struct {
	service *HttpService
}

func (s *ApiServer) GetHello(w http.ResponseWriter, r *http.Request) {
	_, log := logger.FromContext(r.Context())
	log.Info("GetHello")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("World"))
}
