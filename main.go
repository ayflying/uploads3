package main

import (
	_ "uploads3/internal/packed"

	_ "uploads3/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"uploads3/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
