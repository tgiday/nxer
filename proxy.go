package nxer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// NewRedirect return a new proxy handler
func NewRedirect(m map[string]string) func(w http.ResponseWriter, r *http.Request) *httputil.ReverseProxy {
	return func(w http.ResponseWriter, r *http.Request) *httputil.ReverseProxy {
		hst := r.Host
		s := strings.Split(hst, ".")
		subd := s[0]
		//u must be replaced by the url of docker container of sub domain ?
		u, _ := url.Parse(m[subd])
		proxy := httputil.NewSingleHostReverseProxy(u)
		r.URL.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = u.Host
		return proxy
	}
}

// Getdomains return list of subdomains, from the enviroment var DOMAIN
func Getdomains(domainName string) []string {
	subdomains := os.Getenv("DOMAIN")
	subs := strings.Split(subdomains, ",")
	s := []string{}
	for _, v := range subs {
		s = append(s, v+"."+domainName)
	}
	return s
}

// Getservicesmap return map of subdomains to services, from enviroment vars DOMAIN and SERVICE
func Getservicesmap() map[string]string {
	var m = map[string]string{}
	subdomains := os.Getenv("DOMAIN")
	//ser := os.Getenv("SERVICE")
	d := strings.Split(subdomains, ",")
	//c := strings.Split(ser, ",")
	for _, v := range d {
		x := "http://" + v + "runing"
		m[v] = x
	}
	return m
}
