package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/url"

	ws "github.com/NetSepio/webscrape"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {

	websiteURL := "https://github.com/NetSepio"
	domain, err := url.Parse(websiteURL)
	if err != nil {
		log.Fatal(err)
	}
	htmlFile := domain.Host

	ws.CheckDomain(websiteURL)

	hash, _ := ws.AddFileToIpfs(htmlFile)
	fmt.Printf("https://ipfs.infura.io/ipfs/" + hash)

	ws.GetObjectFromIpfs(hash, "output.txt")

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(websiteURL, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
