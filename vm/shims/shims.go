package shims

import (
	"fmt"
	"github.com/Baldomo/Fangs/logger"
	"github.com/Baldomo/Fangs/utils"
	"github.com/Baldomo/v8-go"
)

var noop = func(in v8.CallbackArgs) (*v8.Value, error) { return in.Context.Create(struct{}{}) }

type jsObject map[string]interface{}

// Type Shim represents a function which can modify a context (and its window object)
type Shim func(ctx *v8.Context, shimWindow bool)

var clappr = jsObject{
	"player": func(in v8.CallbackArgs) (*v8.Value, error) {
		return in.Context.Create(jsObject{
			"options": noop,
			"attachTo": noop,
			"listenTo": noop,
			"configure": noop,
			"play": noop,

			// Note: in javascriptEval, this is a function which handles a 'config' object, may be relevant
			"load": noop,
		})
	},
	"events": func(in v8.CallbackArgs) (*v8.Value, error) {
		return in.Context.Create(jsObject{
			"PLAYER_FULLSCREEN": "fullscreen",
			"on":                noop,
			"once":              noop,
			"off":               noop,
			"trigger":           noop,
			"stopListening":     noop,
		})
	},
	"$": func(in v8.CallbackArgs) (*v8.Value, error) {
		return in.Context.Create(jsObject{
			"text": noop,
		})
	},
}

// Overrides the Clappr, LevelSelector and TheaterMode objects in the context
func Clappr(ctx *v8.Context, shimWindow bool) {
	clapprObj, err := ctx.Create(clappr)
	if err != nil {
		logger.Error("Could not create clappr object")
		return
	}

	emptyObj, _ := ctx.Create(jsObject{})

	ctx.Global().Set("Clappr", clapprObj)
	ctx.Global().Set("LevelSelector", emptyObj)
	ctx.Global().Set("TheaterMode", emptyObj)

	if shimWindow {
		window, _ := ctx.Global().Get("window")
		window.Set("Clappr", clapprObj)
		window.Set("LevelSelector", emptyObj)
		window.Set("TheaterMode", emptyObj)
	}
}

var jquery = jsObject{
	"ready":     noop,
	"click":     noop,
	"hide":      noop,
	"show":      noop,
	"mouseup":   noop,
	"mousedown": noop,
	"hasClass":  noop,
	"attr":      noop,
	"post":      noop,
	"get":       noop,
	"append":    noop,
	"cookie": func(in v8.CallbackArgs) (*v8.Value, error) {
		return in.Context.Create(true)
	},
}

// Overrides the jQuery and $ objects in the context
func JQuery(ctx *v8.Context, shimWindow bool) {
	jqueryObj, err := ctx.Create(jquery)
	if err != nil {
		logger.Error("Could not create jQuery object")
		return
	}

	ctx.Global().Set("jQuery", jqueryObj)
	ctx.Global().Set("$", jqueryObj)

	if shimWindow {
		window, _ := ctx.Global().Get("window")
		window.Set("jQuery", jqueryObj)
		window.Set("$", jqueryObj)
	}
}

var jwplayer = jsObject{
	"on":         noop,
	"addButton":  noop,
	"onTime":     noop,
	"onComplete": noop,
}

// Overrides the jwplayer object in the context
func JWPlayer(ctx *v8.Context, shimWindow bool, setupCallback ...v8.Callback) {
	jwplayerObj, err := ctx.Create(jwplayer)
	if err != nil {
		logger.Error("Could not create jwplayer object")
		return
	}

	if len(setupCallback) != 0 {
		jwplayerObj.Set("setup", ctx.Bind("setup", setupCallback[0]))
	} else {
		jwplayerObj.Set("setup", ctx.Bind("setup", noop))
	}

	ctx.Global().Set("jwplayer", jwplayerObj)

	if shimWindow {
		window, _ := ctx.Global().Get("window")
		window.Set("jwplayer", jwplayerObj)
	}
}

const proxy = `let %s;
%s = new Proxy(function () {
	return %s;
}, {
	get: function (target, name, proxy) {
		if (target.hasOwnProperty(name)) {
			return target[name];
		}

		if ('hasOwnProperty' === name) {
			return function () {
				return true;
			}
		} else if ('toString' === name) {
			return ` + "`toString: ${name}`;" + `
		} else if ('length' === name) {
			return 0;
		} else if (name === '__proto__') {
			return proxy ? proxy : %s;
		}
	}
})
`

// Create a native shim object e.g. Document, which allows you call any functions/properties and allow them
// to resolve without throwing a "Cannot read property of undefined" error.
// Useful for bypassing tricky resolvers like Openload which check for browser objects.
// Set useDeepProxy to true if you need to shim complex objects with dynamic properties.
func NativeProxy(ctx *v8.Context, objectName string, useDeepProxy bool) {
	if useDeepProxy {
		ctx.Eval(fmt.Sprintf(proxy, objectName, objectName, objectName, objectName), "fangs_shim.js")
	} else {
		ctx.Global().Set(objectName, ctx.Bind(objectName, noop))
	}

	obj, err := ctx.Global().Get(objectName)
	if err != nil {
		logger.Error("Failed to set JS proxy", "obj_name", objectName, "use_deep_proxy", useDeepProxy)
		return
	}

	toString := ctx.Bind("toString", func(in v8.CallbackArgs) (*v8.Value, error) {
		if useDeepProxy || objectName[0] < 91 && objectName[0] > 64 {
			// Starts with uppercase -> object is a class
			return ctx.Create(fmt.Sprintf("function %s () { [native code] }", objectName))
		}

		return ctx.Create(fmt.Sprintf("[object %s]", utils.Capitalize(objectName)))
	})

	obj.Set("toString", toString)

	__proto__, _ := obj.Get("__proto__")
	__proto__.Set("toString", toString)

	prototype, _ := obj.Get("prototype")
	// TODO: maybe this could just be bound to toString?
	// prototype.Set("toString", toString)
	prototype.Set("toString", ctx.Bind("toString", func(in v8.CallbackArgs) (*v8.Value, error) {
		return in.Context.Create(fmt.Sprintf("[object %s]", utils.Capitalize(objectName)))
	}))
}
