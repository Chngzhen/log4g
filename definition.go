package log4g

import "fmt"

type Level int8

const (
	ERROR Level = 4
	WARN  Level = 3
	INFO  Level = 2
	DEBUG Level = 1
)

type IAppenderDefinition interface {
	getLayout() string
	getLevel() Level
	append(message string, level Level)
}

type StandardAppenderDefinition struct {
	// %d{..} - 日期时间格式。如：%d{yyyy-MM-dd HH:mm:ss.SSS}
	// %l     - 日志等级
	// %f     - 函数
	// %L     - 源码行号
	// %m     - 日志文本
	// %n     - 换行
	layout string
	level  Level
}

func (t *StandardAppenderDefinition) getLayout() string {
	return t.layout
}

func (t *StandardAppenderDefinition) getLevel() Level {
	return t.level
}

func (t *StandardAppenderDefinition) append(message string, level Level) {
	if level >= t.level {
		fmt.Print(message)
	}
}

type LoggerDefinition struct {
	name     string
	appender IAppenderDefinition
}

// Config log4g的配置。
type Config struct {
	loggers []*LoggerDefinition
}

// CallerInfo 调用者信息
type CallerInfo struct {
	// 包名
	packageName string
	// 文件名
	fileName string
	// 函数名
	funcName string
	// 行号
	lineNo int
}
