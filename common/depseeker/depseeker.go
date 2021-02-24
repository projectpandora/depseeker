package depseeker

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/cornelk/hashmap"
	"github.com/projectdiscovery/gologger"
)

// Depseeker represent an instance of the library client
type Depseeker struct {
	options Options
}

// New depseeker instance
func New(options *Options) (*Depseeker, error) {
	depseeker := &Depseeker{}
	return depseeker, nil
}

// waitFor blocks until eventName is received.
// Examples of events you can wait for:
//     init, DOMContentLoaded, firstPaint,
//     firstContentfulPaint, firstImagePaint,
//     firstMeaningfulPaintCandidate,
//     load, networkAlmostIdle, firstMeaningfulPaint, networkIdle
//
// This is not super reliable, I've already found incidental cases where
// networkIdle was sent before load. It's probably smart to see how
// puppeteer implements this exactly.
func waitFor(ctx context.Context, eventName string) error {
	ch := make(chan struct{})
	cctx, cancel := context.WithCancel(ctx)
	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			if e.Name == eventName {
				cancel()
				close(ch)
			}
		}
	})
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func navigateAndWaitFor(url string, eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		_, _, _, err := page.Navigate(url).Do(ctx)
		if err != nil {
			return err
		}

		return waitFor(ctx, eventName)
	}
}

// Run crawl a website and check if there any exposed package
func (d Depseeker) Run(ctx context.Context, url string) ([]Dependency, error) {
	hm := hashmap.HashMap{}
	returnDependencies := []Dependency{}
	// mutex
	mutex := sync.Mutex{}

	// create chrome instance
	options := []chromedp.ExecAllocatorOption{}
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)
	options = append(options, chromedp.UserAgent(d.options.UserAgent))
	options = append(options, chromedp.DisableGPU)
	options = append(options, chromedp.Flag("ignore-certificate-errors", true)) // RIP shittyproxy.go
	options = append(options, chromedp.WindowSize(1920, 1080))
	if d.options.HTTPProxy != "" {
		options = append(options, chromedp.ProxyServer(d.options.HTTPProxy))
	}

	// create context
	chromeCtx, xcancel := chromedp.NewExecAllocator(ctx, options...)
	defer xcancel()

	// start chrome
	// remove the 2nd param if you don't need debug information logged
	ctxt, cancel := chromedp.NewContext(chromeCtx)
	defer cancel()

	chromedp.ListenTarget(
		ctxt,
		func(ev interface{}) {
			if ev, ok := ev.(*network.EventResponseReceived); ok {
				if ev.Type != "Script" {
					return
				}
				go func() {
					// get response body
					c := chromedp.FromContext(ctxt)
					rbp := network.GetResponseBody(ev.RequestID)
					body, err := rbp.Do(cdp.WithExecutor(ctxt, c.Target))
					if err != nil {
						gologger.Error().Msgf("Encountered error: %v", err)
					}
					// grep for packages
					var re = regexp.MustCompile(`(?m)"([a-z\-\_\@\.\/]+)"\s*:\s*"([0-9\^\.\~\*x]+)"`)
					for _, match := range re.FindAllStringSubmatch(string(body), -1) {
						if len(match) >= 3 {
							packageName := strings.TrimSpace(match[1])
							packageVersion := strings.TrimSpace(match[2])
							if packageName != "" {
								isAllow := true
								// check for blacklist
								for _, check := range BlacklistPackageName {
									if check == packageName {
										isAllow = false
										break
									}
								}
								// check package version
								if packageVersion == "" {
									isAllow = false
								} else {
									// https://github.com/Masterminds/semver/blob/master/version.go#L42
									var reVersion = regexp.MustCompile(`v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
										`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
										`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`)
									if reVersion.Match([]byte(packageVersion)) == false {
										isAllow = false
									} else {
										if _, err := strconv.Atoi(packageVersion); err == nil {
											isAllow = false
										}
									}
								}
								if isAllow == true {
									_, loaded := hm.GetOrInsert(packageName, "")
									if loaded == false {
										// fmt.Println(packageName, ev.Response.URL)
										newDependency := Dependency{
											Name:    packageName,
											Version: packageVersion,
										}

										// check if package is existed in npm
										resp, err := http.Get("http://registry.npmjs.com/" + packageName)
										if err == nil {
											body, err := ioutil.ReadAll(resp.Body)
											if err == nil {
												if strings.TrimSpace(string(body)) == "{\"error\":\"Not found\"}" {
													newDependency.IsPrivate = true
												}
											} else {
												gologger.Error().Msgf("Error: %v", err)
											}
										} else {
											gologger.Error().Msgf("Error: %v", err)
										}

										// add to result
										mutex.Lock()
										returnDependencies = append(returnDependencies, newDependency)
										mutex.Unlock()
									}
								}
							}
						}
					}
				}()

			}
		},
	)

	err := chromedp.Run(ctxt, chromedp.Tasks{
		navigateAndWaitFor(url, "networkIdle"),
		chromedp.Sleep(time.Duration(15 * time.Second)),
	})
	if err != nil {
		gologger.Error().Msgf("Error: %v", err)
	}
	return returnDependencies, nil
}
