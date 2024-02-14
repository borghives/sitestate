package data

import (
	"log"
	"sync"

	"github.com/borghives/sitepages"
)

var SITE = "site.json"

type SiteCache struct {
	page   map[string]sitepages.SitePage
	stanza map[string]sitepages.Stanza
}

var (
	cache     *SiteCache
	cacheOnce sync.Once
)

func (sc *SiteCache) GetPage(id string) sitepages.SitePage {
	if sc == nil {
		return sitepages.SitePage{}
	}
	return sc.page[id]
}

func (sc *SiteCache) GetStanza(id string) sitepages.Stanza {
	if sc == nil {
		return sitepages.Stanza{}
	}
	return sc.stanza[id]
}

func LoadSiteCache() {
	cacheOnce.Do(func() {
		// TODO: load from disk
		site := sitepages.LoadSitePages(SITE)
		cache = &SiteCache{
			page:   make(map[string]sitepages.SitePage),
			stanza: make(map[string]sitepages.Stanza),
		}

		for _, page := range site {
			cache.page[page.ID.Hex()] = page
			for _, stanza := range page.StanzaData {
				cache.stanza[stanza.ID.Hex()] = stanza
			}
		}

	})
}

func GetCache() *SiteCache {
	if cache == nil {
		log.Printf("Cache is not initialized")
	}
	return cache
}
