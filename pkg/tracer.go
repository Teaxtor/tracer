package pkg

import (
	"fmt"
	"github.com/raff/godet"
	"log"
	"time"
)

const (
	ResultSuccess = "success"
	ResultUndefined = "undefined"
)

type Tracer struct {
	BrowserConfig BrowserConfig
	ProxyInfo ProxyInfo
	Browsers  map[string]*Browser
	RemotePort int
}

type TraceConfig struct {
	Proxy   string `json:"proxy"`
	Url     string `json:"url"`
}

type TraceResult struct {
	Trace []string `json:"trace"`
	Duration time.Duration `json:"duration"`
	Result string `json:"result"`
}

func New(config BrowserConfig, proxy ProxyInfo, remotePort int) (*Tracer) {
	return &Tracer{
		BrowserConfig: config,
		ProxyInfo: proxy,
		Browsers:  make(map[string]*Browser, 0),
		RemotePort: remotePort,
	}
}

func (t *Tracer) Stop() {
	for country, browser := range t.Browsers {
		fmt.Println("stopping browser for country ", country)
		browser.Close()
	}
}

func (t *Tracer) Trace(config TraceConfig) (*TraceResult, error) {
	browser, err := t.GetBrowser(config.Proxy)

	if err != nil {
		return nil, err
	}

	browser.CallbackEvent(godet.EventClosed, func(params godet.Params) {
		t.Browsers[config.Proxy] = nil
		log.Println("Stopped remote for ", config.Proxy)
	})

	tab, err := browser.NewTab("about:blank")

	if err != nil {
		log.Println("Could not create tab", err)

		return nil, err
	}

	err = browser.ActivateTab(tab)

	if err != nil {
		log.Println("Could not activate tab", err)

		return nil, err
	}

	browser.EnableEvents()

	result := &TraceResult{
		Result: ResultUndefined,
		Trace:  make([]string, 0),
	}

	start := time.Now()
	_, err = browser.Navigate(config.Url)
	result.Duration = time.Since(start)

	if err != nil {
		log.Println("Could not navigate to ", config.Url, err)

		return nil, err
	}

	result.Result = ResultSuccess
	_, navEntries, err := browser.GetNavigationHistory()

	if err != nil {
		log.Println("Unable to get navigation history")
	} else {
		for _, entry := range navEntries {
			result.Trace = append(result.Trace, entry.URL)
		}
	}

	err = browser.CloseTab(tab)

	if err != nil {
		log.Println("Unable to close tab", err)
	}

	return result, nil
}

func (t *Tracer) GetBrowser(proxyKey string) (*Browser, error) {
	if browser, ok := t.Browsers[proxyKey]; ok {
		return browser, nil
	}

	remotePort := t.RemotePort
	t.RemotePort++

	browser, err := NewBrowser(t.BrowserConfig, t.ProxyInfo, proxyKey, remotePort)

	if err != nil {
		return nil, err
	}

	t.Browsers[proxyKey] = browser

	return t.Browsers[proxyKey], nil
}
