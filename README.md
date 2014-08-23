# Decay

*Decay appears to be a Go library that possibly and randomly corrupts data on get requests*

## WTF?

If you have something that you don't want to exist forever then Decay is what you need. On each call a decision (based on several formulae and (pseudo)random numbers) is made - through time your data transforms into total mess.

Currently text and image resources are supported.

## How to use?

### DB
Basic Decay model looks like {id: int, content: string, callCount: int} and is stored in Redis database in two hashsets for each resource: 

resource_name:content id content
resource_name:call_count id callCount

(For images content use image filename - path to image is currently hardcoded)

### API
Decay has public API which now only supports get request:

Update is unsupported by design, create and delete are to be implemented.

To plug in any of Decay modules write something like this:

```
package main

import (
    "github.com/lafaulx/decay/controller"
    "net/http"
)

func main() {
    textController := &controller.TextController{}
    http.HandleFunc("/a/1/text", textController.ProcessRequest)
    http.ListenAndServe(":8080", nil)
}
```

## TODO?
Lots of stuff but, as it's only a concept, future of Decay development is unpredictable.
