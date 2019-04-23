package vm

import (
	"encoding/base64"
	"github.com/Baldomo/Fangs/logger"
	"github.com/Baldomo/Fangs/vm/shims"
	"github.com/Baldomo/v8-go"
	"math"
)

var defaultIsolate *v8.Isolate

func init() {
	defaultIsolate = v8.NewIsolate()
}

// Run given JS code in the default isolate, inside a fresh context
// Consider having a default context containing a default sandbox to avoid having to create a new one
func Run(js string, useDefaultSandbox bool, includeBrowserShims bool, jwPlayerCallback v8.Callback, shims ...shims.Shim) error {
	ctx := defaultIsolate.NewContext()

	if useDefaultSandbox {
		DefaultSandbox(ctx, includeBrowserShims, jwPlayerCallback)
	}

	for _, shim := range shims {
		// Also shim the window object by default
		shim(ctx, true)
	}

	_, err := ctx.Eval(js, "fangs_sandbox.js")
	if err != nil {
		logger.Error("Error running JS code", "err", err.Error(), "js", js[:100], "useDefaultSandbox", useDefaultSandbox, "includeBrowserShims", includeBrowserShims)
		return err
	}

	ctx.Terminate()
	return nil
}

var atob = func(in v8.CallbackArgs) (*v8.Value, error) {
	decoded, err := base64.StdEncoding.DecodeString(in.Arg(0).String())
	if err != nil {
		return nil, err
	}
	return in.Context.Create(string(decoded))
}

var btoa = func(in v8.CallbackArgs) (*v8.Value, error) {
	encoded := base64.StdEncoding.EncodeToString(
		[]byte(in.Arg(0).String()),
	)
	return in.Context.Create(encoded)
}

var sin = func(in v8.CallbackArgs) (*v8.Value, error) {
	num := in.Arg(0).Float64()
	return in.Context.Create(math.Sin(num))
}

type navigator struct {
	UserAgent string `json:"userAgent"`
}

var window = struct {
	Atob      v8.Callback `json:"atob"`
	Btoa      v8.Callback `json:"btoa"`
	Sin       v8.Callback `json:"sin"`
	Navigator navigator   `json:"navigator"`
}{
	Atob:      atob,
	Btoa:      btoa,
	Sin:       sin,
	Navigator: navigator{
		UserAgent: "",
	},
}

// Sets up the default sandbox objects and shims in a given v8.Context
func DefaultSandbox(ctx *v8.Context, includeBrowserShims bool, jwPlayerCallback ...v8.Callback) {
	// Create window object
	windowObj, _ := ctx.Create(window)

	// Add the window object to the context
	ctx.Global().Set("window", windowObj)

	// Include standard shims and also shim window.<shim>
	shims.Clappr(ctx, true)
	shims.JQuery(ctx, true)
	shims.JWPlayer(ctx, true, jwPlayerCallback...)

	if includeBrowserShims {
		// Shim the document object with a special proxy
		shims.NativeProxy(ctx, "Document", false)
		shims.NativeProxy(ctx, "document", false)
	}
}
