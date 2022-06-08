package ballast

import (
    "context"
    "errors"
    "github.com/shirou/gopsutil/mem"
    "runtime"
    "sync"
)

var (
    isExits = 0
    setLock = &sync.RWMutex{}
)

//
// SetBallast
//  @Description: 方法会默认尝试将当前系统的最大内存设置为稳流器，以此来避免单次goroutine小内存频繁触发GC导致不必要的GC占用
//  @param maxSize
//  @return error
//
func SetBallast(ctx context.Context, maxSize ...int) error {
    setLock.Lock()
    defer setLock.Unlock()
    if isExits != 0 {
        return errors.New("already setBallast")
    }
    // getMaxMem
    osMaxSize, err := mem.VirtualMemoryWithContext(ctx)
    if err != nil {
        return err
    }
    maxSzAdvice := osMaxSize.Total >> 2
    if len(maxSize) != 0 && uint64(maxSize[0]) < maxSzAdvice {
        maxSzAdvice = uint64(maxSize[0])
    }
    ballast := make([]byte, maxSzAdvice)
    runtime.KeepAlive(ballast)
    isExits = 1
    return nil
}
