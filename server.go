package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

type initMiddleware struct {
	next xhandler.HandlerC
	ctx  context.Context
}

func (h initMiddleware) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = h.ctx
	h.next.ServeHTTPC(ctx, w, r)
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key", "global context")
	c := xhandler.Chain{}
	c.UseC(func(next xhandler.HandlerC) xhandler.HandlerC {
		return initMiddleware{
			next: next,
			ctx:  ctx,
		}
	})
	c.UseC(xhandler.CloseHandler)
	c.UseC(xhandler.TimeoutHandler(2 * time.Second))

	xh := xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		value := ctx.Value("key").(string)
		w.Write([]byte("Hello " + value))
	})
	http.Handle("/test", c.Handler(xh))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
