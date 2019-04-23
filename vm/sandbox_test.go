package vm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Baldomo/v8-go"
	"testing"
)

var iso *v8.Isolate

func init() {
	iso = v8.NewIsolate()
}

func TestDefaultSandbox(t *testing.T) {
	ctx := iso.NewContext()
	DefaultSandbox(ctx, false)

	b, _ :=ctx.Global().MarshalJSON()

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	fmt.Println(out.String())

	res, _ := ctx.Eval("jwplayer.ping()", "fangs_test.js")
	fmt.Println(res)
}

func TestDefault(t *testing.T) {
	ctx := iso.NewContext()

	// A typical use of Create is to return values from callbacks:
	var nextId int
	getNextIdCallback := func(in v8.CallbackArgs) (*v8.Value, error) {
		nextId++
		return ctx.Create(nextId) // Return the created corresponding v8.Value or an error.
	}

	// Because Create will use reflection to map a Go value to a JS object, it
	// can also be used to easily bind a complex object into the JS VM.
	resetIdsCallback := func(in v8.CallbackArgs) (*v8.Value, error) {
		nextId = 0
		return nil, nil
	}
	myIdAPI, _ := ctx.Create(map[string]interface{}{
		"next":  getNextIdCallback,
		"reset": resetIdsCallback,
		// Can also include other stuff:
		"my_api_version": "v1.2",
	})

	// now let's use those two callbacks and the api value:
	_ = ctx.Global().Set("ids", myIdAPI)
	var res *v8.Value
	res, _ = ctx.Eval(`ids.my_api_version`, `test.js`)
	fmt.Println(`ids.my_api_version =`, res)
	res, _ = ctx.Eval(`ids.next()`, `test.js`)
	fmt.Println(`ids.next() =`, res)
	res, _ = ctx.Eval(`ids.next()`, `test.js`)
	fmt.Println(`ids.next() =`, res)
	res, _ = ctx.Eval(`ids.reset(); ids.next()`, `test.js`)
	fmt.Println(`ids.reset()`)
	fmt.Println(`ids.next() =`, res)
}
