package favpath

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	urlparse "net/url"
	"regexp"

	gq "github.com/PuerkitoBio/goquery"
)

var defaultUserAgent = "favpath"

type Finder struct {
	request *http.Request
	client  *http.Client
}

func NewFinder() *Finder {
	var request http.Request
	request.Method = "GET"
	request.Header = http.Header{}
	request.Header.Set("User-Agent", defaultUserAgent)

	jar, _ := cookiejar.New(nil)

	return &Finder{
		request: &request,
		client:  &http.Client{Jar: jar},
	}
}

func (f *Finder) Header(key, value string) *Finder {
	f.request.Header.Set(key, value)
	return f
}

func (f *Finder) Find(url string) (string, error) {
	favurl, err := f.FindFromDoc(url)
	if err != nil {
		return f.defaultPath(url), err
	}

	return favurl, nil
}

func (f *Finder) defaultPath(url string) string {
	p, _ := urlparse.Parse(url)
	p.Path = "/favicon.ico"

	return p.String()
}

// link[rel='icon'],link[rel='Icon'],
// link[rel='shortcut icon'],link[rel='Shortcut Icon']
func (f *Finder) FindFromDoc(url string) (string, error) {
	doc, err := f.Doc(url)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?i)^(shortcut )?icon$`)
	sel := doc.Find("link").FilterFunction(func(i int, s *gq.Selection) bool {
		val, _ := s.Attr("rel")
		return re.MatchString(val)
	})

	href, ok := sel.Attr("href")
	if !ok {
		return "", errors.New("Favicon does not found")
	}

	u, _ := urlparse.Parse(href)
	if u.Scheme != "" {
		return href, nil
	}

	p, _ := urlparse.Parse(url)
	p.Path = u.Path
	return p.String(), nil
}

func (f *Finder) Doc(urlStr string) (*gq.Document, error) {
	url, _ := urlparse.Parse(urlStr)
	f.request.URL = url

	resp, err := f.client.Do(f.request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return gq.NewDocumentFromResponse(resp)
}
