package pkg

import (
	"github.com/raff/godet"
	"log"
	"os/exec"
	"strconv"
	"time"
)

type BrowserConfig struct {
	ScreenWidth         int
	ScreenHeight        int
	UseMobile           bool
	UserAgent           string
	Timeout             int
	RemoteConnects      int
	WaitBetweenConnects time.Duration
	Headless 			bool
}

type Browser struct {
	*godet.RemoteDebugger
}

func NewBrowser(config BrowserConfig, info ProxyInfo, proxyKey string, remotePort int) (*Browser, error) {
	proxySetting := info.GenerateProxy(proxyKey)

	binary := "chromium"
	options := []string{"--remote-debugging-port=" + strconv.Itoa(remotePort), "--hide-scrollbars", "--disable-extensions", "--disable-gpu"}
	target := "about:blank"

	if len(proxySetting) > 0 {
		options = append(options, proxySetting)
	}

	if config.Headless {
		options = append(options, "--headless")
	}

	options = append(options, target)

	log.Println("starting browser for ", proxyKey)
	cmd := exec.Command(binary, options...)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var remote *godet.RemoteDebugger
	var err error

	for i := 0 ; i < config.RemoteConnects; i++ {
		if i > 0 {
			time.Sleep(config.WaitBetweenConnects)
		}

		remote, err = godet.Connect("localhost:" + strconv.Itoa(remotePort), true)

		if err == nil {
			break
		}

		log.Println("Unable to connect to remote ", remotePort, " ", err)
	}

	if err != nil {
		log.Fatal("Failed to launch remote ", remotePort, " ", err)
	}

	remote.SetUserAgent(config.UserAgent)
	remote.SetVisibleSize(config.ScreenWidth, config.ScreenHeight)
	remote.SetDeviceMetricsOverride(config.ScreenWidth, config.ScreenHeight, 3, config.UseMobile, false)

	return &Browser{remote}, nil
}

func (b *Browser) EnableEvents() {
	b.RuntimeEvents(true)
	b.NetworkEvents(true)
	b.PageEvents(true)
	b.EmulationEvents(true)
	b.SetControlNavigations(true)
}