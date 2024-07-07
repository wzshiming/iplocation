package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/netip"

	"github.com/wzshiming/iplocation"
	"github.com/wzshiming/iplocation/source"
)

type Server struct {
	source map[string]source.Source

	indexData json.RawMessage
}

type paths struct {
	Paths []string `json:"paths"`
}

func NewServer(source map[string]source.Source) *Server {
	indexData, _ := json.Marshal(paths{
		Paths: []string{
			"/api/",
			"/api/ip/{ip}/asn",
			"/api/ip/{ip}/country",
			"/api/ip/{ip}/city",
		},
	})

	return &Server{
		source:    source,
		indexData: indexData,
	}
}

func (s *Server) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/", s.Index)
	mux.HandleFunc("GET /api/ip/{ip}/asn", s.ASN)
	mux.HandleFunc("GET /api/ip/{ip}/country", s.Country)
	mux.HandleFunc("GET /api/ip/{ip}/city", s.City)
	return mux
}

func (s *Server) Index(rw http.ResponseWriter, r *http.Request) {
	responseJSON(rw, http.StatusOK, s.indexData)
}

func (s *Server) ASN(rw http.ResponseWriter, r *http.Request) {
	s.infoGet(rw, r, "asn")
}

func (s *Server) Country(rw http.ResponseWriter, r *http.Request) {
	s.infoGet(rw, r, "country")
}

func (s *Server) City(rw http.ResponseWriter, r *http.Request) {
	s.infoGet(rw, r, "city")
}

func (s *Server) infoGet(rw http.ResponseWriter, r *http.Request, kind string) {
	qsource := r.URL.Query().Get("source")
	if qsource == "" {
		qsource = "geolite2"
	}
	qprovider := r.URL.Query().Get("provider")
	if qprovider == "" {
		qprovider = "sapics"
	}
	source := qprovider + "-" + qsource + "-" + kind

	src, ok := s.source[source]
	if !ok {
		responseJSON(rw, http.StatusNotFound, errorResponse{"not found"})
		return
	}

	ip := r.PathValue("ip")

	addr, err := netip.ParseAddr(ip)
	if err != nil {
		responseJSON(rw, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}

	data, err := src.Lookup(addr)
	if err != nil {
		if errors.Is(err, iplocation.ErrNotFound) {
			responseJSON(rw, http.StatusNotFound, errorResponse{err.Error()})
		} else {
			responseJSON(rw, http.StatusBadRequest, errorResponse{err.Error()})
		}
		return
	}

	responseJSON(rw, http.StatusOK, data)
}

func responseJSON(rw http.ResponseWriter, statusCode int, v any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(v)
}

type errorResponse struct {
	Error string `json:"error"`
}
