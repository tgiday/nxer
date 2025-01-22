package nxer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewRedirect return a new proxy handler
func NewRedirect(m map[string]string) func(w http.ResponseWriter, r *http.Request) *httputil.ReverseProxy {
	return func(w http.ResponseWriter, r *http.Request) *httputil.ReverseProxy {
		hst := r.Host
		s := strings.Split(hst, ".")
		subd := s[0]
		//u must be replaced by the url of docker container of sub domain ?
		_, ok := m[subd]
		if !ok {
			u, _ := url.Parse(m["www"])
			proxy := httputil.NewSingleHostReverseProxy(u)
			r.URL.Host = u.Host
			r.URL.Scheme = u.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = u.Host
			return proxy
		}
		u, _ := url.Parse(m[subd])
		proxy := httputil.NewSingleHostReverseProxy(u)
		r.URL.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = u.Host
		return proxy
	}
}

func maper[S any, D any](s []S, d *[]D, m func(S) D) {
	if cap(*d) == 0 {
		(*d) = make([]D, len(s))
	}
	for i, e := range s {
		(*d)[i] = m(e)
	}
}

// Getdomains return list of subdomains, from the enviroment variable
func Getdomains(sub []string, domain string) []string {
	domains := []string{}
	domfunc := func(s string) string {
		return s + "." + domain
	}
	maper(sub, &domains, domfunc)
	domains = append(domains, domain)
	return domains
}

// Getservicesmap return map of subdomains to services, from enviroment variable and suffix string(eg: "www":"http://www" + sfx)
func Getservicesmap(sub []string, sfx string) map[string]string {
	serfunc := func(s string) string {
		return "http://" + s + sfx
	}
	services := []string{}
	maper(sub, &services, serfunc)
	m := map[string]string{}
	for i, v := range sub {
		m[v] = services[i]
	}

	return m
}
