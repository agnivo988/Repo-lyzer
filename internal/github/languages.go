package github

import gocache "github.com/patrickmn/go-cache"

func (c *Client) GetLanguages(owner, repo string) (map[string]int, error) {
	cacheKey := "languages:" + owner + "/" + repo
	if cached, found := c.cache.Get(cacheKey); found {
		return copyLanguagesMap(cached.(map[string]int)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyLanguagesMap(cached.(map[string]int)), nil
		}

		var langs map[string]int
		if err := c.get("https://api.github.com/repos/"+owner+"/"+repo+"/languages", &langs); err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, langs, gocache.DefaultExpiration)
		return copyLanguagesMap(langs), nil
	})
	if err != nil {
		return nil, err
	}
	return v.(map[string]int), nil
}
