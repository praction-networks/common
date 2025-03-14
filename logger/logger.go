package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global logger instance
var (
	logInstance         *zap.Logger
	logLevel            = zap.NewAtomicLevel()
	initOnce            sync.Once
	errInit             error
	requestScopedLogger *zap.Logger // Request-specific logger
	mu                  sync.Mutex  // Mutex to protect the request logger
)

// LoggerConfig defines the configuration for the logger
type LoggerConfig struct {
	LogstashHost    string
	LogstashPort    int
	Index           string
	LogstashEnabled bool
	LogLevel        string
}

// CasbinLogger is a wrapper for your custom logger to integrate with Casbin
type CasbinLogger struct {
	enabled bool
}

// NewCasbinLogger creates a new CasbinLogger instance// NewCasbinLogger creates a new CasbinLogger instance
func NewCasbinLogger() *CasbinLogger {
	// Adjust caller skip to avoid redundant logging frames
	logInstance = logInstance.WithOptions(zap.AddCallerSkip(2))

	return &CasbinLogger{
		enabled: true, // Enable logging by default
	}
}

// EnableLog enables or disables logging
func (c *CasbinLogger) EnableLog(enabled bool) {
	c.enabled = enabled
}

// IsEnabled checks if logging is enabled
func (c *CasbinLogger) IsEnabled() bool {
	return c.enabled
}

// LogModel logs Casbin model details
func (c *CasbinLogger) LogModel(model [][]string) {
	if c.IsEnabled() {
		Debug("Casbin Model", "Model", model)
	}
}

// LogPolicy logs Casbin policy details
func (c *CasbinLogger) LogPolicy(policy map[string][][]string) {
	if c.IsEnabled() {
		Debug("Casbin Policy", "Policy", policy)
	}
}

// LogRole logs Casbin role-related information
func (c *CasbinLogger) LogRole(roles []string) {
	if c.IsEnabled() {
		Debug("Casbin Role", "roles", roles)
	}
}

// LogEnforce logs Casbin enforcement results
func (c *CasbinLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	if c.IsEnabled() {
		Info("Casbin Enforcement",
			"matcher", matcher,
			"request", request,
			"result", result,
			"explains", explains,
		)
	}
}

// LogError logs error messages
func (c *CasbinLogger) LogError(err error, v ...string) {
	if c.IsEnabled() {
		Error("Casbin Error",
			"error", err,
			"details", v,
		)
	}
}

// logs general Casbin messages
func (c *CasbinLogger) Log(v ...interface{}) {
	if c.IsEnabled() {
		Info("Casbin Log", v...)
	}
}

// logstashWriter implements zapcore.WriteSyncer for sending logs to Logstash
type logstashWriter struct {
	url string
}

// InitializeLogger initializes the logger with default log level INFO.
func InitializeLogger(config LoggerConfig) (*zap.Logger, error) {

	initOnce.Do(func() {
		if config.LogstashEnabled && (config.LogstashHost == "" || config.LogstashPort == 0) {
			errInit = fmt.Errorf("invalid logstash configuration: host or port is missing")
			return
		}

		// Construct the Logstash URL
		logstashURL := fmt.Sprintf("http://%s:%d", config.LogstashHost, config.LogstashPort)
		// Create atomic log level with default level INFO

		if config.LogLevel == "" {
			config.LogLevel = "info" // Default level
		}

		// Map logLevel string to zapcore.Level
		var zapLogLevel zapcore.Level
		switch config.LogLevel {
		case "debug":
			zapLogLevel = zapcore.DebugLevel
		case "info":
			zapLogLevel = zapcore.InfoLevel
		case "warn":
			zapLogLevel = zapcore.WarnLevel
		case "error":
			zapLogLevel = zapcore.ErrorLevel
		case "fatal":
			zapLogLevel = zapcore.FatalLevel
		default:
			zapLogLevel = zapcore.InfoLevel // Default to INFO level if invalid log level is passed
			fmt.Printf("Invalid log level '%s' provided. Defaulting to INFO level.\n", config.LogLevel)
		}

		// Set the global atomic log level
		logLevel.SetLevel(zapLogLevel)

		var cores []zapcore.Core

		textEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Logger",
			CallerKey:      "Caller",
			MessageKey:     "Message",
			StacktraceKey:  "Stacktrace",
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		})

		textCore := zapcore.NewCore(textEncoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), logLevel)
		cores = append(cores, textCore)

		if config.LogstashEnabled {
			jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
			logstashCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(&logstashWriter{url: logstashURL}), logLevel)
			cores = append(cores, logstashCore)
		}

		var core zapcore.Core
		if len(cores) == 1 {
			core = cores[0]
		} else {
			core = zapcore.NewTee(cores...)
		}

		defaultFields := []zap.Field{
			zap.String("service", config.Index),
		}
		addEnvironmentFields(&defaultFields)

		// Create the logger
		logInstance = zap.New(core, zap.AddCaller()).With(
			defaultFields...,
		)

		// Log initialization message
		logInstance.Info("Logger initialized",
			zap.String("default_level", logLevel.String()),
			zap.Bool("logstash_enabled", config.LogstashEnabled),
		)
	})

	return logInstance, errInit
}

// GetGlobalLogger returns the global logger instance.
func GetGlobalLogger() *zap.Logger {
	return logInstance
}

// SetDefaultRequestLogger sets a logger specific to the current request lifecycle.
func SetDefaultRequestLogger(fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	requestScopedLogger = logInstance.With(fields...)
}

// ClearDefaultRequestLogger clears the request-specific logger.
func ClearDefaultRequestLogger() {
	mu.Lock()
	defer mu.Unlock()
	requestScopedLogger = nil
}

// getDefaultLogger returns the request-specific logger if set, or the global logger.
func getDefaultLogger() *zap.Logger {
	mu.Lock()
	defer mu.Unlock()
	if requestScopedLogger != nil {
		return requestScopedLogger
	}
	return logInstance
}

// addIfNotEmpty adds a field to the defaultFields slice if the value is not empty.
func addIfNotEmpty(fields *[]zap.Field, key, value string) {
	if value != "" {
		*fields = append(*fields, zap.String(key, value))
	}
}

func (w *logstashWriter) Write(p []byte) (n int, err error) {

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", w.url, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logInstance.Warn("Failed to send log to Logstash",
			zap.String("logstash_url", w.url),
			zap.Error(err))
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Logstash returned status %d\n", resp.StatusCode)
		return 0, fmt.Errorf("logstash returned status %d", resp.StatusCode)
	}

	return len(p), nil
}

// UpdateLogLevel updates the log level dynamically.
func UpdateLogLevel(newLevel string) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(newLevel)); err != nil {
		logInstance.Warn("Invalid log level provided",
			zap.String("provided_level", newLevel),
			zap.Error(err))
		return fmt.Errorf("invalid log level: %s", newLevel)
	}
	logLevel.SetLevel(level)
	logInstance.Info("Log level updated", zap.String("new_level", level.String()))
	return nil
}

// Sync implements zapcore.WriteSyncer interface (no-op for Logstash)
func (w *logstashWriter) Sync() error {
	return nil
}

// WithContext returns a logger instance with additional context fields.
func WithContext(fields ...zap.Field) *zap.Logger {
	ensureInitialized()
	return logInstance.With(fields...)
}

func ensureInitialized() {
	if logInstance == nil {
		panic("Logger not initialized. Call InitializeLogger first.")
	}
}

// addEnvironmentFields adds Kubernetes-related environment fields to the logger.
func addEnvironmentFields(fields *[]zap.Field) {
	addIfNotEmpty(fields, "pod_name", os.Getenv("POD_NAME"))
	addIfNotEmpty(fields, "namespace", os.Getenv("POD_NAMESPACE"))
	addIfNotEmpty(fields, "service_name", os.Getenv("DEPLOYMENT_NAME"))
	addIfNotEmpty(fields, "node_name", os.Getenv("NODE_NAME"))
	addIfNotEmpty(fields, "pod_uuid", os.Getenv("POD_UID"))
}

func logWithLevel(level string, msg string, args ...interface{}) {
	ensureInitialized()

	logger := getDefaultLogger() // Use request-specific or global logger
	fields := []zap.Field{}

	// Check for error and add as a field
	if len(args) > 0 {
		if err, ok := args[0].(error); ok {
			fields = append(fields, zap.String("Error", err.Error()))
			args = args[1:] // Remove the error from args
		}
	}

	// // Process remaining arguments as key-value pairs
	fields = append(fields, withFields(args)...) // Add remaining arguments as fields
	fields = append(fields, logLevelField(level))

	// Add stack trace for specific log levels if needed
	if level == "debug" || level == "fatal" || level == "panic" {
		fields = append(fields, stackTraceField())
	}

	// Call the appropriate log method based on the level
	switch level {
	case "info":
		logger.WithOptions(zap.AddCallerSkip(2)).Info(msg, fields...)
	case "warn":
		logger.WithOptions(zap.AddCallerSkip(2)).Warn(msg, fields...)
	case "error":
		logger.WithOptions(zap.AddCallerSkip(2)).Error(msg, fields...)
	case "debug":
		logger.WithOptions(zap.AddCallerSkip(2)).Debug(msg, fields...)
	case "fatal":
		logger.WithOptions(zap.AddCallerSkip(2)).Fatal(msg, fields...)
	case "panic":
		logger.WithOptions(zap.AddCallerSkip(2)).Panic(msg, fields...)
	}
}

func Info(msg string, args ...interface{}) {
	logWithLevel("info", msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logWithLevel("warn", msg, args...)
}

func Error(msg string, args ...interface{}) {
	logWithLevel("error", msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logWithLevel("debug", msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logWithLevel("fatal", msg, args...)
}

func Panic(msg string, args ...interface{}) {
	logWithLevel("panic", msg, args...)
}

// Helper functions

func logLevelField(level string) zap.Field {
	return zap.String("loglevel", level)
}

func stackTraceField() zap.Field {
	return zap.String("stacktrace", string(debug.Stack()))
}
func withFields(args []interface{}) []zap.Field {
	fields := []zap.Field{}

	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			switch key := args[i].(type) {
			case string:
				fields = append(fields, zap.Any(key, args[i+1]))
			default:
				logInstance.WithOptions(zap.AddCallerSkip(3)).Warn("Invalid key in arguments passed to logger", zap.Any("invalid_key", key))
			}
		} else {
			logInstance.WithOptions(zap.AddCallerSkip(4)).Warn("Missing value for key in arguments passed to logger", zap.Any("key", args[i]))
		}
	}

	return fields
}
