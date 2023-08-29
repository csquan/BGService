package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
)

// AddFileOut only support unix os
func AddFileOut(logFilePath string, level, days int) (err error) {
	/*var absPath string
	if absPath, err = filepath.Abs(logFilePath); err != nil {
		return errors.WithStack(err)
	}

	logDirPath := filepath.Dir(absPath)
	if _, err = os.Stat(logDirPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(logDirPath, 0755); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	if err = errors.WithStack(err); err != nil {
		return
	}
	if err = unix.Access(logDirPath, unix.W_OK); err != nil {
		err = errors.Wrapf(err, "%s is not writable", logDirPath)
		return
	}

	var logf *rotatelogs.RotateLogs
	logf, err = rotatelogs.New(
		absPath+".%Y-%m-%dT%H:%M",
		rotatelogs.WithLinkName(absPath),
		rotatelogs.WithMaxAge(time.Duration(days)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err = errors.WithStack(err); err != nil {
		return err
	}
	hook := &fileHook{
		formatter: &logrus.TextFormatter{
			DisableTimestamp: false,
			CallerPrettyfier: callerPrettyfier,
		},
		levels: getHookLevel(level),
		rotate: logf,
	}

	logrus.AddHook(hook)*/
	return
}

type fileHook struct {
	formatter logrus.Formatter
	levels    []logrus.Level
	rotate    *rotatelogs.RotateLogs
}

func (hook *fileHook) Fire(entry *logrus.Entry) error {
	if enableDefaultFieldMap {
		for key, value := range defaultFieldMap {
			if _, ok := entry.Data[key]; !ok {
				entry.Data[key] = value
			}
		}
	}
	formatBytes, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.rotate.Write(formatBytes)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to write file on filehook(entry.String)%v", err)
		return err
	}
	return nil
}

func (hook *fileHook) Levels() []logrus.Level {
	return hook.levels
}
