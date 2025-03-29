package logger

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logInstance         *zap.Logger
	logLevel            = zap.NewAtomicLevel()
	initOnce            sync.Once
	errInit             error
	requestScopedLogger *zap.Logger
	mu                  sync.Mutex
)

var sensitiveKeys = []string{"password", "secret", "token", "apikey", "otp", "pan", "aadhaar", "mobile", "phone", "wifi", "wifi_password"}
var personalKeys = []string{"email", "phone", "mobile", "username"}

type LoggerConfig struct {
	LogLevel string
}

type CasbinLogger struct {
	enabled bool
}

func NewCasbinLogger() *CasbinLogger {
	logInstance = logInstance.WithOptions(zap.AddCallerSkip(2))
	return &CasbinLogger{enabled: true}
}

func (c *CasbinLogger) EnableLog(enabled bool) { c.enabled = enabled }
func (c *CasbinLogger) IsEnabled() bool        { return c.enabled }

func (c *CasbinLogger) LogModel(model [][]string) {
	if c.IsEnabled() {
		Debug("Casbin Model", "model", model)
	}
}

func (c *CasbinLogger) LogPolicy(policy map[string][][]string) {
	if c.IsEnabled() {
		Debug("Casbin Policy", "policy", policy)
	}
}

func (c *CasbinLogger) LogRole(roles []string) {
	if c.IsEnabled() {
		Debug("Casbin Role", "roles", roles)
	}
}

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

func (c *CasbinLogger) LogError(err error, v ...string) {
	if c.IsEnabled() {
		Error("Casbin Error",
			"error", err,
			"details", v,
		)
	}
}

func (c *CasbinLogger) Log(v ...interface{}) {
	if c.IsEnabled() {
		Info("Casbin Log", v...)
	}
}

func InitializeLogger(config LoggerConfig) error {
	initOnce.Do(func() {
		if config.LogLevel == "" {
			config.LogLevel = "info"
		}

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
			zapLogLevel = zapcore.InfoLevel
			fmt.Printf("Invalid log level '%s' provided. Defaulting to INFO level.\n", config.LogLevel)
		}

		logLevel.SetLevel(zapLogLevel)

		jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		})

		core := zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), logLevel)

		defaultFields := []zap.Field{}
		addEnvironmentFields(&defaultFields)

		logInstance = zap.New(core, zap.AddCaller()).With(defaultFields...)

		logInstance.Info("Logger initialized", zap.String("default_level", logLevel.String()))
	})

	return errInit
}

func GetGlobalLogger() *zap.Logger {
	if logInstance == nil {
		panic("Logger not initialized")
	}
	return logInstance
}

func SetDefaultRequestLogger(fields ...zap.Field) {
	mu.Lock()
	defer mu.Unlock()
	requestScopedLogger = logInstance.With(fields...)
}

func ClearDefaultRequestLogger() {
	mu.Lock()
	defer mu.Unlock()
	requestScopedLogger = nil
}

func getDefaultLogger() *zap.Logger {
	mu.Lock()
	defer mu.Unlock()
	if requestScopedLogger != nil {
		return requestScopedLogger
	}
	return logInstance
}

func UpdateLogLevel(newLevel string) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(newLevel)); err != nil {
		logInstance.Warn("Invalid log level provided", zap.String("provided_level", newLevel), zap.Error(err))
		return fmt.Errorf("invalid log level: %s", newLevel)
	}
	logLevel.SetLevel(level)
	logInstance.Info("Log level updated", zap.String("new_level", level.String()))
	return nil
}

func Sync() {
	if logInstance != nil {
		_ = logInstance.Sync()
	}
}

func WithContext(fields ...zap.Field) *zap.Logger {
	ensureInitialized()
	return logInstance.With(fields...)
}

func ensureInitialized() {
	if logInstance == nil {
		panic("Logger not initialized. Call InitializeLogger first.")
	}
}

func addEnvironmentFields(fields *[]zap.Field) {
	// addIfNotEmpty(fields, "pod", os.Getenv("POD_NAME"))
	// addIfNotEmpty(fields, "namespace", os.Getenv("POD_NAMESPACE"))
	// addIfNotEmpty(fields, "service", os.Getenv("DEPLOYMENT_NAME"))
	// addIfNotEmpty(fields, "node", os.Getenv("NODE_NAME"))
	// addIfNotEmpty(fields, "pod_uuid", os.Getenv("POD_UID"))
	addIfNotEmpty(fields, "env", os.Getenv("ENVIRONMENT"))
}

func addIfNotEmpty(fields *[]zap.Field, key, value string) {
	if value != "" {
		*fields = append(*fields, zap.String(key, value))
	}
}

func logWithLevel(level string, msg string, args ...interface{}) {
	ensureInitialized()
	logger := getDefaultLogger()
	fields := []zap.Field{}

	errorCount := 0
	cleanedArgs := []interface{}{}

	for _, arg := range args {
		if err, ok := arg.(error); ok {
			key := "error"
			if errorCount > 0 {
				key = fmt.Sprintf("error_%d", errorCount+1)
			}
			fields = append(fields, zap.String(key, err.Error()))
			errorCount++
		} else {
			cleanedArgs = append(cleanedArgs, arg)
		}
	}

	fields = append(fields, withFields(cleanedArgs)...)
	// fields = append(fields, logLevelField(level))

	if level == "debug" || level == "fatal" || level == "panic" {
		fields = append(fields, stackTraceField())
	}

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

func Info(msg string, args ...interface{})  { logWithLevel("info", msg, args...) }
func Warn(msg string, args ...interface{})  { logWithLevel("warn", msg, args...) }
func Error(msg string, args ...interface{}) { logWithLevel("error", msg, args...) }
func Debug(msg string, args ...interface{}) { logWithLevel("debug", msg, args...) }
func Fatal(msg string, args ...interface{}) { logWithLevel("fatal", msg, args...) }
func Panic(msg string, args ...interface{}) { logWithLevel("panic", msg, args...) }

// func logLevelField(level string) zap.Field {
// 	return zap.String("loglevel", level)
// }

func stackTraceField() zap.Field {
	return zap.String("stacktrace", string(debug.Stack()))
}

func withFields(args []interface{}) []zap.Field {
	fields := []zap.Field{}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			switch key := args[i].(type) {
			case string:
				value := args[i+1]
				keyLower := strings.ToLower(key)

				if strVal, ok := value.(string); ok {
					if contains(sensitiveKeys, keyLower) {
						fields = append(fields, zap.String(key, "REDACTED"))
					} else if contains(personalKeys, keyLower) {
						fields = append(fields, zap.String(key, mask(strVal)))
					} else {
						fields = append(fields, zap.String(key, strVal))
					}
				} else {
					fields = append(fields, zap.Any(key, value))
				}
			default:
				logInstance.WithOptions(zap.AddCallerSkip(3)).Warn("Invalid key in arguments passed to logger", zap.Any("invalid_key", key))
			}
		} else {
			logInstance.WithOptions(zap.AddCallerSkip(4)).Warn("Missing value for key in arguments passed to logger", zap.Any("key", args[i]))
		}
	}
	return fields
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(item, s) {
			return true
		}
	}
	return false
}

func mask(input string) string {
	if len(input) <= 2 {
		return "**"
	}
	return input[:1] + strings.Repeat("*", len(input)-2) + input[len(input)-1:]
}
