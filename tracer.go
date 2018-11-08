package tracer

import (
	"fmt"
	"github.com/raff/godet"
	"log"
	"time"
)

const (
	ResultSuccess = "success"
	ResultTimeout = "timeout"
	ResultUndefined = "undefined"
)

type Tracer struct {
	BrowserConfig BrowserConfig
	ProxyInfo ProxyInfo
	Browsers  map[string]*Browser
	RemotePort int
}

type TraceConfig struct {
	Proxy string
	Url string
	UserAgent string
	TotalTimeout int
	NavigationFinished func(params godet.Params) bool
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
	browser, err := t.GetBrowser(config.Proxy, config.UserAgent)

	if err != nil {
		return nil, err
	}

	browser.CallbackEvent(godet.EventClosed, func(params godet.Params) {
		t.Browsers[config.Proxy] = nil
		log.Println("Stopped remote for ", config.Proxy)
	})

	tab, err := browser.NewTab(config.Url)

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

	defer func() {
		result.Duration = time.Since(start)
	}()

	timeout := false

	browser.CallbackEvent("Page.navigationRequested", func(params godet.Params) {
		url := params.String("url")

		result.Trace = append(result.Trace, url)

		log.Println("navigation request for ", params.String("url"))

		var response godet.NavigationResponse

		switch {
		case timeout:
			result.Result = ResultTimeout
			response = godet.NavigationCancelAndIgnore

		case config.NavigationFinished(params):
			result.Result = ResultSuccess
			response = godet.NavigationCancelAndIgnore

		default:
			response = godet.NavigationProceed
		}

		browser.ProcessNavigation(params.Int("navigationId"), response)
	})

	browser.CallbackEvent("Emulation.virtualTimeBudgetExpired", func(params godet.Params) {
		timeout = true
	})

	_, err = browser.Navigate(config.Url)

	if err != nil {
		log.Println("Could not navigate to ", config.Url, err)

		return nil, err
	}

	err = browser.CloseTab(tab)

	if err != nil {
		log.Println("Unable to close tab", err)
	}

	return result, nil
}

func (t *Tracer) GetBrowser(proxyKey string, userAgent string) (*Browser, error) {
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
