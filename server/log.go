package server

import (
	"io"
	"os"
	"sync/atomic"

	srvlog "github.com/elitecodegroovy/gmessage/logger"
)

// Logger interface of the NATS Server
type Logger interface {

	// Log a notice statement
	Noticef(format string, v ...interface{})

	// Log a fatal error
	Fatalf(format string, v ...interface{})

	// Log an error
	Errorf(format string, v ...interface{})

	// Log a debug statement
	Debugf(format string, v ...interface{})

	// Log a trace statement
	Tracef(format string, v ...interface{})
}

// 设置服务器的日志
func (s *Server) ConfigureLogger() {
	var (
		log Logger

		// 镜像服务操作参数
		opts = s.getOpts()
	)

	syslog := opts.Syslog
	if isWindowsService() && opts.LogFile == "" {
		// 设置系统日志标识
		syslog = true
	}

	if opts.LogFile != "" {
		//指定日志文件
		log = srvlog.NewFileLogger(opts.LogFile, opts.Logtime, opts.Debug, opts.Trace, true)
	} else if opts.RemoteSyslog != "" {
		//指定远程系统日志
		log = srvlog.NewRemoteSysLogger(opts.RemoteSyslog, opts.Debug, opts.Trace)
	} else if syslog {
		//指定系统日志文件
		log = srvlog.NewSysLogger(opts.Debug, opts.Trace)
	} else {
		colors := true
		// 日志文字颜色设置
		stat, err := os.Stderr.Stat()
		if err != nil || (stat.Mode()&os.ModeCharDevice) == 0 {
			colors = false
		}
		log = srvlog.NewStdLogger(opts.Logtime, opts.Debug, opts.Trace, colors, true)
	}

	s.SetLogger(log, opts.Debug, opts.Trace)
}

//设置服务的日志
func (s *Server) SetLogger(logger Logger, debugFlag, traceFlag bool) {
	if debugFlag {
		atomic.StoreInt32(&s.logging.debug, 1)
	} else {
		atomic.StoreInt32(&s.logging.debug, 0)
	}
	if traceFlag {
		atomic.StoreInt32(&s.logging.trace, 1)
	} else {
		atomic.StoreInt32(&s.logging.trace, 0)
	}
	s.logging.Lock()
	if s.logging.logger != nil {
		// 检测日志是否实现io.Closer.
		if l, ok := s.logging.logger.(io.Closer); ok {
			if err := l.Close(); err != nil {
				s.Errorf("Error closing logger: %v", err)
			}
		}
	}
	s.logging.logger = logger
	s.logging.Unlock()
}

// If the logger is a file based logger, close and re-open the file.
// This allows for file rotation by 'mv'ing the file then signaling
// the process to trigger this function.
func (s *Server) ReOpenLogFile() {
	// Check to make sure this is a file logger.
	s.logging.RLock()
	ll := s.logging.logger
	s.logging.RUnlock()

	if ll == nil {
		s.Noticef("File log re-open ignored, no logger")
		return
	}

	// Snapshot server options.
	opts := s.getOpts()

	if opts.LogFile == "" {
		s.Noticef("File log re-open ignored, not a file logger")
	} else {
		fileLog := srvlog.NewFileLogger(opts.LogFile,
			opts.Logtime, opts.Debug, opts.Trace, true)
		s.SetLogger(fileLog, opts.Debug, opts.Trace)
		s.Noticef("File log re-opened")
	}
}

// 记录通知语句
func (s *Server) Noticef(format string, v ...interface{}) {
	s.executeLogCall(func(logger Logger, format string, v ...interface{}) {
		logger.Noticef(format, v...)
	}, format, v...)
}

// 记录错误
func (s *Server) Errorf(format string, v ...interface{}) {
	s.executeLogCall(func(logger Logger, format string, v ...interface{}) {
		logger.Errorf(format, v...)
	}, format, v...)
}

// 记录致命错误
func (s *Server) Fatalf(format string, v ...interface{}) {
	s.executeLogCall(func(logger Logger, format string, v ...interface{}) {
		logger.Fatalf(format, v...)
	}, format, v...)
}

// 记录调试错误
func (s *Server) Debugf(format string, v ...interface{}) {
	if atomic.LoadInt32(&s.logging.debug) == 0 {
		return
	}

	s.executeLogCall(func(logger Logger, format string, v ...interface{}) {
		logger.Debugf(format, v...)
	}, format, v...)
}

// 记录追踪信息
func (s *Server) Tracef(format string, v ...interface{}) {
	if atomic.LoadInt32(&s.logging.trace) == 0 {
		return
	}

	s.executeLogCall(func(logger Logger, format string, v ...interface{}) {
		logger.Tracef(format, v...)
	}, format, v...)
}

// 内容日志调用方法
func (s *Server) executeLogCall(f func(logger Logger, format string, v ...interface{}), format string, args ...interface{}) {
	s.logging.RLock()
	defer s.logging.RUnlock()
	if s.logging.logger == nil {
		return
	}

	f(s.logging.logger, format, args...)
}
