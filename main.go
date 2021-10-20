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
	ctx, err := v8go.NewContext(nil) // creates a new V8 context with a new Isolate aka VM
	if err != nil {
		panic(err)
	}
	ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
	ctx.RunScript("const result = add(3, 4)", "main.js") // any functions previously added to the context can be called
	val, err := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
	if err != nil {
		panic(err)
	}
	fmt.Printf("addition result: %s\n", val)
}
