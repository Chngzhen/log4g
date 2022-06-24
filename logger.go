package log4g

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/Chngzhen/gxlib/xdates"
)

var config *Config
var defLogger = &LoggerDefinition{
	name:     "*",
	appender: &StandardAppenderDefinition{"%d{yyyy-MM-dd HH:mm:ss.SSS} %l %f[%L] - %m%n", INFO},
}

// Build 构造日志系统。
func Build(pConfig *Config) error {
	if pConfig != nil {
		config = pConfig
	}
	return nil
}

// Info 以INFO等级打印日志。
func Info(message string, args ...interface{}) {
	var callerInfo *CallerInfo
	pkgName, fileName, funcName, lineNo, ok := getCallerInfo()
	if ok {
		callerInfo = &CallerInfo{pkgName, fileName, funcName, lineNo}
	}
	log(callerInfo, INFO, message, args...)
}

// Warn 以WARN等级打印日志。
func Warn(message string, args ...interface{}) {
	var callerInfo *CallerInfo
	pkgName, fileName, funcName, lineNo, ok := getCallerInfo()
	if ok {
		callerInfo = &CallerInfo{pkgName, fileName, funcName, lineNo}
	}
	log(callerInfo, WARN, message, args...)
}

// Error 以ERROR等级打印日志。
func Error(message string, args ...interface{}) {
	var callerInfo *CallerInfo
	pkgName, fileName, funcName, lineNo, ok := getCallerInfo()
	if ok {
		callerInfo = &CallerInfo{pkgName, fileName, funcName, lineNo}
	}
	log(callerInfo, ERROR, message, args...)
}

// Debug 以Debug等级打印日志。
func Debug(message string, args ...interface{}) {
	var callerInfo *CallerInfo
	pkgName, fileName, funcName, lineNo, ok := getCallerInfo()
	if ok {
		callerInfo = &CallerInfo{pkgName, fileName, funcName, lineNo}
	}
	log(callerInfo, DEBUG, message, args...)
}

func log(caller *CallerInfo, level Level, message string, args ...interface{}) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	loggers := chooseLoggers(caller, level)
	if loggers == nil || len(loggers) == 0 {
		loggers = []*LoggerDefinition{defLogger}
	}
	for _, logger := range loggers {
		str := parseLayout(logger.appender.getLayout(), caller, level, message)
		logger.appender.append(str, level)
	}
}

func chooseLoggers(pCaller *CallerInfo, level Level) []*LoggerDefinition {
	if config == nil {
		return nil
	}

	pTopName := pCaller.packageName
	pSecondName := pTopName + "." + pCaller.fileName
	pFullName := pSecondName + "." + pCaller.funcName

	// 若存在多个相同name的logger，则后定义的优先于先定义的；若存在多个层次的logger，则优先匹配更精确的。
	var loggers []*LoggerDefinition
	for _, logger := range config.loggers {
		if logger.name == pFullName || logger.name == pSecondName || logger.name == pTopName {
			loggers = append(loggers, logger)
		}
	}
	return loggers
}

// 标准输出颜色设置。
// 格式：%c[isHighLight;background;fontColorm..%c[0m
// 说明：`%c`取值`0x1B`时表示临时设置颜色；`[..m`是颜色定义体，内容为`0`时表示清空临时设置，恢复成默认颜色。
//   isHighLight 是否高亮。取值：0-默认，1-高亮，4-下划线，5-闪烁，7-反白，8-不可见。
//   background  背景颜色。取值：0-默认，40-黑，41-红，42-绿，43-黄，44-蓝，45-紫红，46-青蓝，47-白
//   fontColor   字体颜色。取值：0-默认，30-黑，31-红，32-绿，33-黄，34-蓝，35-紫红，36-青蓝，37-白
func parseLayout(layout string, caller *CallerInfo, level Level, message string) string {
	if len(layout) == 0 {
		return ""
	}
	var format strings.Builder
	if layout[0] != '%' {
		format.WriteByte(layout[0])
		format.WriteString(parseLayout(layout[1:], caller, level, message))
		return format.String()
	}

	placeholder := layout[0:2]
	switch placeholder {
	// 日期时间。后面必须接格式。
	case "%d":
		noDateFormat := true
		dateFormatEndIndex := 3
		for noDateFormat && dateFormatEndIndex < len(layout)-2 {
			char := layout[dateFormatEndIndex]
			if char == '}' {
				noDateFormat = false
			} else {
				dateFormatEndIndex++
			}
		}
		dateFormat := layout[3:dateFormatEndIndex]
		format.WriteString(xdates.FormatNow2String(dateFormat))
		format.WriteString(parseLayout(layout[(dateFormatEndIndex+1):], caller, level, message))
		return format.String()
	// 行号。
	case "%L":
		format.WriteString(strconv.Itoa(caller.lineNo))
		break
	// 日志等级
	case "%l":
		switch level {
		case INFO:
			format.WriteString("INFO")
			break
		case WARN:
			format.WriteString("WARN")
			break
		case ERROR:
			format.WriteString("ERROR")
			break
		case DEBUG:
			format.WriteString("DEBUG")
			break
		}
		break
	// 函数
	case "%f":
		format.WriteString(caller.packageName + "." + caller.fileName + "." + caller.funcName)
		break
	// 日志文本
	case "%m":
		format.WriteString(message)
		break
	// 换行符
	case "%n":
		format.WriteString("\n")
		break
	}
	format.WriteString(parseLayout(layout[2:], caller, level, message))
	return format.String()
}

// getCallerInfo 获取当前方法的调用者信息，包括包名、文件名、方法名盒行号。若无法获取，则返回false。
func getCallerInfo() (string, string, string, int, bool) {
	pc, filePath, lineNo, ok := runtime.Caller(2)
	if ok {
		fileName := path.Base(filePath)
		fileExt := path.Ext(filePath)
		fileNameWithoutExt := fileName[0 : len(fileName)-len(fileExt)]

		function := runtime.FuncForPC(pc)
		functionName := function.Name()
		functionNames := strings.Split(functionName, ".")
		return functionNames[0], fileNameWithoutExt, functionNames[1], lineNo, true
	}
	return "", "", "", 0, false
}