package contracts

type LoggerInterface interface {
	Print(log string)
	Printf(format string, args ...interface{})
}
