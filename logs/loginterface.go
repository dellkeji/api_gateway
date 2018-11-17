package logs

// Log : describe function for log
type Log interface {
	Debug(format string, args ...interface{}) // ...表示接收可变参数
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Close() // 文件需要进行关闭操作
}
