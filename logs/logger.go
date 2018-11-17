package logs

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	utils "apigw_golang/utils"
)

// FileLog : struct for file configs
type FileLog struct {
	logPath  string
	logName  string
	file     *os.File
	warnFile *os.File
}

// NewFileLog : get file log configs
func NewFileLog(logPath, logName string) Log {
	log := &FileLog{
		logPath: logPath,
		logName: logName,
	}

	log.init()

	return log
}

func (f *FileLog) init() {
	// 一般日志
	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755) // os.O_CREATE 创建文件 os.O_APPEND 追加写入 os.O_WRONLY 只写操作
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed, err: %v", filename, err))
	}

	f.file = file

	// 错误日志
	warnfilename := fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	warnfile, err := os.OpenFile(warnfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755) // os.O_CREATE 创建文件 os.O_APPEND 追加写入 os.O_WRONLY 只写操作
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed, err: %v", warnfilename, err))
	}

	f.warnFile = warnfile
}

func (f *FileLog) writeLog(file *os.File, level int, format string, args ...interface{}) {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")
	// 这个数字格式是固定的不能改变的，但是-和:可以更换
	levelStr := LogLevelString(level)
	fileName, funcName, lineNo := utils.GetLineInfo()
	//由于这里返回的是全路径的，但是我们不需要，所以我们只需要文件名以及相关的即可
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(file, "%s %s [%s/%s:%d] %s\n", nowStr, levelStr, fileName, funcName, lineNo, msg)
}

// Debug :
func (f *FileLog) Debug(format string, args ...interface{}) {
	f.writeLog(f.file, DebugLevel, format, args...)
}

// Trace :
func (f *FileLog) Trace(format string, args ...interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	stackInfo := fmt.Sprintf("%s", string(buf[:n]))
	fmt.Fprintf(f.warnFile, stackInfo, args...)
}

// Info :
func (f *FileLog) Info(format string, args ...interface{}) {
	f.writeLog(f.file, InfoLevel, format, args...)
}

// Warn :
func (f *FileLog) Warn(format string, args ...interface{}) {
	f.writeLog(f.warnFile, WarnLevel, format, args...)
}

// Error :
func (f *FileLog) Error(format string, args ...interface{}) {
	f.writeLog(f.warnFile, ErrorLevel, format, args...)
}

// Fatal :
func (f *FileLog) Fatal(format string, args ...interface{}) {
	f.writeLog(f.warnFile, FatalLevel, format, args...)
}

// Close :
func (f *FileLog) Close() {
	f.file.Close()
	f.warnFile.Close()
}

// LogLevelString :
func LogLevelString(level int) (levelStr string) {
	switch level {
	case DebugLevel:
		levelStr = "DEBUG"
	case TraceLevel:
		levelStr = "TRACE"
	case InfoLevel:
		levelStr = "INFO"
	case WarnLevel:
		levelStr = "WARN"
	case ErrorLevel:
		levelStr = "ERROR"
	case FatalLevel:
		levelStr = "FATAL"
	}
	return
}
