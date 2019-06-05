package formatter

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type ModuleFormatter struct {
	defaultFormatter logrus.Formatter
	whitelist        ModulesMap
	allowAllModules  bool
}

type ModulesMap map[string]logrus.Level

func New(modules ModulesMap) (*ModuleFormatter, error) {
	return NewWithFormatter(modules, &logrus.TextFormatter{})
}

func NewWithFormatter(modules ModulesMap, formatter logrus.Formatter) (*ModuleFormatter, error) {
	logrus.SetLevel(logrus.TraceLevel)

	return &ModuleFormatter{
		defaultFormatter: formatter,
		whitelist:        modules,
		allowAllModules:  len(modules) == 0,
	}, nil
}

func (f *ModuleFormatter) Format(e *logrus.Entry) ([]byte, error) {
	entry := *e
	module, entryHasModuleField := entry.Data["module"]

	// allow any logs that don't have the module field
	if f.allowAllModules {
		return f.defaultFormatter.Format(&entry)
	}

	if entryHasModuleField {
		delete(entry.Data, "module")
		entry.Message = fmt.Sprintf("[%s] %s", module.(string), entry.Message)

		// for the whitelisted modules, allow only the entries with level >= configured
		level, whitelisted := f.whitelist[module.(string)]
		if whitelisted {
			if entry.Level <= level {
				return f.defaultFormatter.Format(&entry)
			} else {
				return nil, nil
			}
		}
	}

	// for non-whitelisted modules, apply the global level
	if globalLevel, ok := f.whitelist["*"]; ok && entry.Level <= globalLevel {
		return f.defaultFormatter.Format(&entry)
	}

	return nil, nil
}
