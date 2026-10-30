// Package resources prepares game resources independently of HTTP requests.
package resources

import (
	"bunker-core/protocol/mpay/android"
	"bunker-core/protocol/mpay/windows"
	"context"
	"log"
	"time"
)

func Start(ctx context.Context) {
	go refresh(ctx, "android", android.RefreshPublishedPackages, 15*time.Minute, time.Minute)
	go refresh(ctx, "windows", windows.RefreshPublishedPackages, 15*time.Minute, time.Minute)
}

func refresh(ctx context.Context, platform string, prepare func() error, interval, retry time.Duration) {
	delay := retry
	for {
		if ctx.Err() != nil {
			return
		}
		next := interval
		if err := prepare(); err != nil {
			log.Printf("%s resource preparation: %v", platform, err)
			next = delay
			if delay < interval/2 {
				delay *= 2
			} else {
				delay = interval
			}
		} else {
			delay = retry
		}
		timer := time.NewTimer(next)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}
