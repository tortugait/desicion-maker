package log

import "go.uber.org/zap"

//nolint:gochecknoglobals
var Logger *zap.SugaredLogger

//nolint:gochecknoinits
func init() {
	if Logger != nil {
		return
	}

	l, _ := zap.NewDevelopment()
	Logger = l.Sugar()
}
