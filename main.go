package main

import (
	"fmt"

	"rogchap.com/v8go"
)

// exampleScript is a basic script that presents a function being called with an object
// that has an accessor function attached to it.
// This also presents use of a global "jarvis" object with callback functions.
const exampleScript = `
function thing(request) {
	req_id = request.get_request_id()
  jarvis.add_note(req_id)
}
`

func main() {
	// Setup the notes object we want to update:
	notes := []string{}

	// Setup new isolate (VM) for this execution.
	iso, err := v8go.NewIsolate()
	if err != nil {
		panic(err)
	}

	// Setup the global "jarvis" object with callback.
	// a template that represents a JS function
	addNoteFn, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		for _, v := range info.Args() {
			notes = append(notes, v.String())
		}
		return nil // you can return a value back to the JS caller if required
	})
	if err != nil {
		panic(err)
	}

	// so now we are setting the add_note function to an object.
	jarvisObj, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		panic(err)
	}
	err = jarvisObj.Set("add_note", addNoteFn)
	if err != nil {
		panic(err)
	}

	// Setup the request object with callback.
	// a template that represents a JS function
	getReqIDFn, err := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		val, err := v8go.NewValue(iso, "some_kind_of_req_id")
		if err != nil {
			// Work out error handling
			panic(err)
		}
		return val
	})
	if err != nil {
		panic(err)
	}

	// so now we are setting the add_note function to an object.
	requestObj, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		panic(err)
	}
	err = requestObj.Set("get_request_id", getReqIDFn)
	if err != nil {
		panic(err)
	}

	// Add global objects to global environment ("jarvis" and args)
	global, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		panic(err)
	}
	err = global.Set("__req_object", requestObj)
	if err != nil {
		panic(err)
	}
	err = global.Set("jarvis", jarvisObj)
	if err != nil {
		panic(err)
	}

	// Add arguments to the context
	// TODO: This is currently being done in a slightly hacky way, let me know if there are other ways.

	// Setup V8 context with the existing isolate.
	ctx, err := v8go.NewContext(iso, global)
	if err != nil {
		panic(err)
	}

	// Execute code to setup environment
	_, err = ctx.RunScript(exampleScript, "userScript.js") // return a value in JavaScript back to Go
	if err != nil {
		panic(err)
	}

	fmt.Printf("before: %v\n", notes)

	// Execute hook with arguments
	val, err := ctx.RunScript("thing(__req_object)", "run_thing_hook.js") // return a value in JavaScript back to Go
	if err != nil {
		panic(err)
	}

	fmt.Printf("val: %v\n", val)
	fmt.Printf("after: %v\n", notes)

	// ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
	// ctx.RunScript("const result = add(3, 4)", "main.js") // any functions previously added to the context can be called
	// val, err := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("addition result: %s\n", val)
}
