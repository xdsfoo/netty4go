//Created by zhbinary on 2019-01-29.
//Email: zhbinary@gmail.com
package channel

import (
	"context"
	"errors"
	"github.com/zhbinary/heng/concurrent"
	"github.com/zhbinary/heng/types"
	"sync"
)

var ErrEventLoopClosed = errors.New("Event loop closed ")

type RoutinePerChannelEventLoop struct {
	tasks             chan types.RunnablePromise
	complete          chan error
	mutex             sync.RWMutex
	closed            bool
	parent            types.EventLoop
	ctx               context.Context
	cancel            context.CancelFunc
	terminationFuture types.Promise
	channel           types.Channel
}

func (this *RoutinePerChannelEventLoop) Execute(task types.Runnable) {
	panic("implement me")
}

func (this *RoutinePerChannelEventLoop) IsShutDown() bool {
	return this.closed
}

func (this *RoutinePerChannelEventLoop) ShutdownGracefully() types.ChannelFutrue {
	this.closed = true
	//this.complete <- nil
	close(this.tasks)
	// todo
	return nil
}

func (this *RoutinePerChannelEventLoop) Next() types.EventLoop {
	panic("implement me")
}

func (this *RoutinePerChannelEventLoop) Register(channel types.Channel) types.ChannelFutrue {
	promise := concurrent.NewDefaultChannelPromise(channel)
	channel.Unsafe().Register(this, promise)
	this.channel = channel
	return promise
}

func (this *RoutinePerChannelEventLoop) Register0(promise types.ChannelPromise) types.ChannelFutrue {
	panic("implement me")
}

func (this *RoutinePerChannelEventLoop) Register1(channel types.Channel, promise types.ChannelPromise) types.ChannelFutrue {
	panic("implement me")
}

func (this *RoutinePerChannelEventLoop) Submit(task types.Runnable) (future types.Future) {
	promise := concurrent.NewDefaultRunnablePromise(task)
	future = promise
	if this.closed {
		promise.SetFailure(ErrEventLoopClosed)
		return
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.tasks <- promise
	return
}

func (this *RoutinePerChannelEventLoop) Parent() types.EventLoopGroup {
	return this.parent
}

func (this *RoutinePerChannelEventLoop) run() {
	for {
		select {
		case promise, ok := <-this.tasks:
			if !ok {
				// Quit
				return
			}
			promise.Run()
			promise.SetSuccess()
			break
		case err := <-this.complete:
			if err != nil {
				//logrus.Errorf("Event loop exit err:%d\n", err)
			}
			return
		}
	}
}
