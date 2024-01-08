package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	Policy struct {
		userAgent  string
		disallowed []string
		delay      time.Duration
		sitemap    string
	}

	Policies map[string]*Policy
)

func (pol *Policy) Disallowed(link *url.URL) bool {
	for _, url := range pol.disallowed {
		url = strings.TrimSpace(url)

		match := strings.Contains(link.Path, url)

		if match {
			log.Println("This url was disallowed: ", link.String())
			return true
		}
	}

	return false
}

func (pol *Policy) Delay(u *url.URL) time.Duration {
	return pol.delay
}

func (pol *Policy) Sitemap(u *url.URL) string {
	return pol.sitemap
}

func downloadPolicy(u url.URL) *Policy {
	pol := Policy{}
	u.Path = "robots.txt"
	rsp, err := http.Get(u.String())
	if err != nil || rsp.StatusCode != 200 {
		return &pol
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
		return &pol
	}

	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		comp := strings.Split(line, ": ")

		if len(comp) > 1 {
			pat := strings.ReplaceAll(comp[1], "*", "")
			switch strings.ToLower(comp[0]) {
			case "user-agent":
				pol.userAgent = pat
			case "disallow":
				pol.disallowed = append(pol.disallowed, pat)
			case "crawl-delay":
				i, err := strconv.Atoi(comp[1])
				if err != nil {
					pol.delay = 0
					log.Fatal(err)
				} else {
					pol.delay = time.Duration(i) * 1000
				}
			case "sitemap":
				pol.sitemap = comp[1]
			}
		}
	}

	return &pol
}

func (pols *Policies) getPolicy(u *url.URL) *Policy {
	hn := u.Hostname()
	if _, exists := (*pols)[hn]; !exists {
		(*pols)[hn] = downloadPolicy(*u)
	}
	return (*pols)[hn]
}

func (pols *Policies) Disallowed(u string) bool {
	if l, err := url.Parse(u); err == nil {
		if p := pols.getPolicy(l); p != nil {
			return p.Disallowed(l)
		}
	} else {
		log.Println(err)
	}

	return false
}

func (pols *Policies) Delay(u string) time.Duration {
	if l, err := url.Parse(u); err == nil {
		if p := pols.getPolicy(l); p != nil {
			return p.Delay(l)
		}
	}
	return 0
}

func (pols *Policies) Sitemap(u string) string {

	xml := ""

	if l, err := url.Parse(u); err == nil {
		if p := pols.getPolicy(l); p != nil {
			xml = p.Sitemap(l)
		}
	} else {
		log.Println(err)
	}

	return xml
}

func NewPolicies() *Policies {
	p := make(Policies)
	return &p
}
