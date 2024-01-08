package main

type (
	Sitemap struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	}

	SitemapIndex struct {
		Sitemaps []Sitemap `xml:"sitemap"`
	}

	TheUrls struct {
		Locations    []string `xml:"url>loc"`
		LastModified string   `xml:"lastmod"`
	}
)
