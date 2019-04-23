package shims

import (
	"encoding/json"
	"fmt"
	"github.com/Baldomo/v8-go"
	"testing"
)

var iso *v8.Isolate

func init() {
	iso = v8.NewIsolate()
}

func TestClappr(t *testing.T) {
	// TODO
	t.Parallel()

	t.Run("standard shim", func(t *testing.T) {
		ctx := iso.NewContext()

		Clappr(ctx, false)

		obj, err := ctx.Global().Get("Clappr")
		if err != nil {
			t.Error(err)
		}

		objJson, err := obj.MarshalJSON()
		if err != nil {
			t.Error(err)
		}

		fmt.Println(string(objJson))
	})

	t.Run("window shimmed", func(t *testing.T) {
		ctx := iso.NewContext()

		Clappr(ctx, true)

		window, err := ctx.Global().Get("window")
		if err != nil {
			t.Error(err)
		}

		obj, err := window.Get("clappr")
		if err != nil {
			t.Error(err)
		}

		objJson, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			t.Error(err)
		}

		t.Log(string(objJson))
	})
}
