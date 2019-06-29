package main

import (
	"fmt"
	"github.com/zhbinary/heng/buffer"
	"github.com/zhbinary/heng/channel/embedded"
	"github.com/zhbinary/heng/handler"
	"github.com/zhbinary/heng/types"
)

func main() {
	byteBuf := buffer.NewHeapBytebuf(1024)
	for i := 0; i < 9; i++ {
		byteBuf.WriteUint8(uint8(i))
	}

	ch := embedded.NewChannel(&In1{})
	if !ch.WriteInbound(byteBuf) {
		fmt.Println("err")
	}

	if !ch.Finish() {
		fmt.Println("err")
	}
}

type In1 struct {
	*handler.ChannelInboundHandlerAdapter
}

func (this *In1) ChannelActive(ctx types.ChannelHandlerContext) {
	fmt.Println("In1 ChannelActive")
	ctx.FireChannelActive()
}

func (this *In1) ChannelRead(ctx types.ChannelHandlerContext, msg interface{}) () {
	fmt.Println("In1 ChannelRead", msg)
	ctx.FireChannelRead(msg)
}
