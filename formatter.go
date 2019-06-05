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

func (f *ModuleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	module, entryHasModuleField := entry.Data["module"]

	// allow any logs that don't have the module field
	if f.allowAllModules {
		return f.defaultFormatter.Format(entry)
	}

	if entryHasModuleField {
		entry.Data = rebuildData(entry.Data)
		entry.Message = fmt.Sprintf("[%s] %s", module.(string), entry.Message)

		// for the whitelisted modules, allow only the entries with level >= configured
		level, whitelisted := f.whitelist[module.(string)]
		if whitelisted {
			if entry.Level <= level {
				return f.defaultFormatter.Format(entry)
			} else {
				return nil, nil
			}
		}
	}

	// for non-whitelisted modules, apply the global level
	if globalLevel, ok := f.whitelist["*"]; ok && entry.Level <= globalLevel {
		return f.defaultFormatter.Format(entry)
	}

	return nil, nil
}

func rebuildData(data logrus.Fields) logrus.Fields {
	newData := logrus.Fields{}

	for k, v := range data {
		if k == "module" {
			continue
		}
		newData[k] = v
	}

	return newData
}
