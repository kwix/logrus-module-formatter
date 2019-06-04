# logrus-module-formatter
This package is intended to be used as [logrus](https://github.com/sirupsen/logrus) formatter. 

You can use this as a wrapper over any other logrus formatter with the added benefit of allowing you to filter what actually gets printed by so-called "modules". 

A `module` is whatever you want it to be. Just specify it via fields when logging something `logrus.WithField("module", "my-special-module").Info("<3")` and it will be considered as such.

Ok, ok, too many words ... here's an example:
- for any kind of log that belongs to the module `very-bugged-module-that-needs-a-lot-of-debugging`, I want to see logs `>= logrus.TraceLevel` while for all other modules I only want to see logs `>= logrus.ErrorLevel`. 
 
## Installation
This package supports go modules:
```bash
go get github.com/kwix/logrus-module-formatter
```

## Configuration
The `ModuleFormatter` accepts a `ModulesMap` as configuration.
 
`ModulesMap` is in the form of `map[string]logrus.Level`. Here's an example:
```
ModulesMap{
    "*":     logrus.WarnLevel,
    "test":  logrus.InfoLevel,
    "test1": logrus.DebugLevel,
}
``` 

Notes:
- you should use `"*": [level]` to specify the global logging level; if this doesn't exist, everything that's not in the whitelist will be ingored
    - this formatter automatically calls `logrus.SetLevel(logrus.TraceLevel)` as a workaround to the fact that a log wouldn't reach the formatter if the level is not `>= logrus global level`

## Usage
```go
import (
	modules "github.com/kwix/logrus-module-formatter"
	"github.com/sirupsen/logrus"
)

f, err := modules.New(modules.ModulesMap{
    "*":     logrus.WarnLevel,
    "test":  logrus.InfoLevel,
    "test1": logrus.DebugLevel,
})
if err != nil {
    panic(err)
}

logrus.SetFormatter(f)

logrus.WithField("module", "test").Debug("This should be ignored")
logrus.WithField("module", "test1").Info("This should be displayed")
logrus.WithField("module", "test2").Warn("This should be displayed too")
```

Output: 
```
time="2019-06-04T21:46:44+03:00" level=info msg="This should be displayed"
time="2019-06-04T21:46:44+03:00" level=warning msg="This should be displayed too"
```