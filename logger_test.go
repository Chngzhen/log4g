package log4g

import "testing"

func TestInfo(t *testing.T) {
	err := Build(&Config{
		loggers: []*LoggerDefinition{
			{
				name:     "log4g",
				appender: &StandardAppenderDefinition{"%d{yyyy-MM-dd HH:mm:ss} %l %f[%L] - %m%n", INFO},
			},
			{
				name:     "log4g.log4g_test",
				appender: &StandardAppenderDefinition{"%d{yyyy-MM-dd HH:mm:ss.SS} %l %f[%L] - %m%n", DEBUG},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	Debug("Hello, world")
	Info("Hello, world")
}
