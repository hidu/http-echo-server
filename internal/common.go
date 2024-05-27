// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/27

package internal

import (
	"context"
	"time"
)

func ctxSleep(ctx context.Context, dur time.Duration) bool {
	tm := time.NewTimer(dur)
	select {
	case <-tm.C:
		return true
	case <-ctx.Done():
		return false
	}
}
