package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	Debug   bool
	Headers []string
}

// Test for headless clients via User-agent and return a raw response rather than HTML
var headlessClients = []string{"curl", "Wget", "Go", "ddclient"}
var headlessMatcher mux.MatcherFunc = func(r *http.Request, m *mux.RouteMatch) bool {
	for _, c := range headlessClients {
		if strings.Contains(r.UserAgent(), c) {
			return true
		}
	}
	return false
}

func New(headers []string) *Server {
	return &Server{Headers: headers}
}

func ipFromRequest(headers []string, r *http.Request) (net.IP, error) {
	remoteIP := ""
	for _, h := range headers {
		remoteIP = r.Header.Get(h) // try our trusted headers

		if http.CanonicalHeaderKey(h) == "X-Forwarded-For" {
			/* X-Forwarded-For can be a CSV list of proxies, eg:
			 * <client>, <proxy1>, <proxy2>. We just want the client
			 */
			s := strings.Index(h, ",")
			if s == -1 {
				remoteIP = h
			} else {
				remoteIP = h[:s]
			}
		}
		if remoteIP != "" {
			break
		}
	}
	if remoteIP == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, err
		}
		remoteIP = host
	}
	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return nil, fmt.Errorf("Error parsing IP from: %s", remoteIP)
	}
	return ip, nil
}

func (s *Server) response(r *http.Request) (net.IP, error) {
	ip, err := ipFromRequest(s.Headers, r)
	if err != nil {
		return nil, err
	}
	return ip, nil
}

func (s *Server) HeadlessHandler(rw http.ResponseWriter, r *http.Request) {
	ip, err := s.response(r)
	if err != nil {
		http.Error(rw, "Unable to parse IP", http.StatusBadRequest)
	}
	fmt.Fprintln(rw, ip.String())
}

func (s *Server) JSONHandler(rw http.ResponseWriter, r *http.Request) {
	ip, err := s.response(r)
	// TODO: Response as struct, marshal struct
	data, err := json.Marshal(map[string]string{"ip": ip.String()})
	if err != nil {
		http.Error(rw, "Unable to parse IP", http.StatusBadRequest)
	}
	fmt.Fprintln(rw, string(data))
}

// TODO: return HTML
func (s *Server) HTMLHandler(rw http.ResponseWriter, r *http.Request) {
	ip, err := s.response(r)
	if err != nil {
		http.Error(rw, "Unable to parse IP", http.StatusBadRequest)
	}
	fmt.Fprintln(rw, "Your IP Address: "+ip.String())
}

func (s *Server) loggingMiddleware(nxt http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Header)
	})
}

func (s *Server) Handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", s.JSONHandler).Methods("GET").Headers("Accept", "application/json")
	r.HandleFunc("/", s.HeadlessHandler).Methods("GET").MatcherFunc(headlessMatcher)
	r.HandleFunc("/", s.HTMLHandler).Methods("GET")
	r.Use(s.loggingMiddleware)

	return r
}

func (s *Server) ListenAndServe(port string) error {
	return http.ListenAndServe(port, s.Handler())
}
