package logger

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

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
	asyncWriter         *asyncWriteSyncer
	complianceConfig    *ComplianceConfig
)

// ============================================
// GDPR & COMPLIANCE CONFIGURATION
// ============================================

// ComplianceMode defines the level of data protection
type ComplianceMode string

const (
	// ComplianceModeStrict - Maximum protection, all PII redacted (GDPR, HIPAA)
	ComplianceModeStrict ComplianceMode = "strict"
	// ComplianceModeModerate - PII masked, IDs pseudonymized (SOC 2, ISO 27001)
	ComplianceModeModerate ComplianceMode = "moderate"
	// ComplianceModeMinimal - Basic sensitive data redaction only
	ComplianceModeMinimal ComplianceMode = "minimal"
)

// ComplianceConfig holds GDPR and compliance settings
type ComplianceConfig struct {
	Mode               ComplianceMode
	AnonymizeIP        bool   // Anonymize last octet of IP addresses
	PseudonymizeUserID bool   // Hash user IDs for privacy
	Salt               string // Salt for pseudonymization (should be from env)
	RetentionDays      int    // Log retention period (for documentation)
	DataRegion         string // Data residency region (EU, US, etc.)
	Enabled            bool   // Enable/disable compliance features (for dev mode)
}

// DefaultComplianceConfig returns a GDPR-compliant default configuration
func DefaultComplianceConfig() *ComplianceConfig {
	env := os.Getenv("ENVIRONMENT")
	// Enable compliance by default in production, disable in dev
	enabled := env == "production" || env == "staging" || env == "prod"

	return &ComplianceConfig{
		Mode:               ComplianceModeModerate,
		AnonymizeIP:        enabled,
		PseudonymizeUserID: false,
		Salt:               os.Getenv("LOG_PSEUDONYMIZATION_SALT"),
		RetentionDays:      90,
		DataRegion:         os.Getenv("DATA_REGION"),
		Enabled:            enabled,
	}
}

// ============================================
// SENSITIVE & PERSONAL DATA DEFINITIONS
// ============================================

// sensitiveKeys - Data that must NEVER be logged (complete redaction)
// Compliance: GDPR Art. 32, PCI-DSS, HIPAA
var sensitiveKeys = []string{
	// Authentication & Security
	"password", "passwd", "pwd", "pass",
	"secret", "secretkey", "secret_key",
	"token", "accesstoken", "access_token", "refreshtoken", "refresh_token",
	"apikey", "api_key", "apitoken", "api_token",
	"bearer", "authorization", "auth",
	"otp", "totp", "mfa", "2fa", "pin",
	"privatekey", "private_key", "privkey",
	"passphrase", "passcode",
	"sessionid", "session_id", "sessiontoken", "session_token",

	// Financial Data (PCI-DSS)
	"creditcard", "credit_card", "cardnumber", "card_number",
	"cvv", "cvc", "cvv2", "securitycode", "security_code",
	"bankaccount", "bank_account", "accountnumber", "account_number",
	"routingnumber", "routing_number", "iban", "swift", "bic",
	"expiry", "expirydate", "expiry_date", "expirationdate",

	// Indian PII (DPDP Act 2023)
	"pan", "pannumber", "pan_number",
	"aadhaar", "aadhar", "uidai", "aadharnumber", "aadhaar_number",
	"voterid", "voter_id", "epicnumber", "epic_number",
	"drivinglicense", "driving_license", "dlnumber", "dl_number",
	"passport", "passportnumber", "passport_number",
	"ration", "rationcard", "ration_card",

	// US PII (CCPA, HIPAA)
	"ssn", "socialsecurity", "social_security", "socialsecuritynumber",
	"taxid", "tax_id", "ein", "itin",

	// Biometric Data (GDPR Art. 9 - Special Category)
	"biometric", "fingerprint", "faceid", "face_id", "facedata", "face_data",
	"retina", "iris", "voiceprint", "voice_print",

	// Health Data (HIPAA, GDPR Art. 9)
	"healthdata", "health_data", "medicalrecord", "medical_record",
	"diagnosis", "prescription", "insurance", "insuranceid",

	// WiFi/Network Credentials
	"wifi", "wifipassword", "wifi_password", "wifikey", "wifi_key",
	"networkkey", "network_key", "psk", "wpakey", "wpa_key",

	// Encryption Keys & Hashes (even hashes shouldn't be logged)
	"hash", "passwordhash", "password_hash", "encryptionkey", "encryption_key",
	"salt", "iv", "nonce", "cipher",
}

// personalKeys - PII that should be masked but may need partial visibility
// Compliance: GDPR Art. 5 (Data Minimization)
var personalKeys = []string{
	// Contact Information
	"email", "emailaddress", "email_address", "mail",
	"phone", "phonenumber", "phone_number", "telephone", "tel",
	"mobile", "mobilenumber", "mobile_number", "cell", "cellphone",
	"whatsapp", "whatsappnumber", "whatsapp_number",

	// Personal Identifiers
	"username", "user_name", "userid", "user_id",
	"firstname", "first_name", "fname",
	"lastname", "last_name", "lname", "surname",
	"middlename", "middle_name", "mname",
	"fullname", "full_name", "name", "displayname", "display_name",

	// Location Data (GDPR considers this PII)
	"address", "streetaddress", "street_address",
	"city", "state", "province", "country", "region",
	"zipcode", "zip_code", "postalcode", "postal_code", "pincode", "pin_code",
	"latitude", "lat", "longitude", "lng", "lon", "geolocation", "geo",
	"location", "coordinates", "coords",

	// Date of Birth & Age
	"dob", "dateofbirth", "date_of_birth", "birthdate", "birth_date",
	"age", "birthday",

	// Device & Network Identifiers
	"ip", "ipaddress", "ip_address", "clientip", "client_ip", "remoteip", "remote_ip",
	"useragent", "user_agent", "ua",
	"deviceid", "device_id", "macaddress", "mac_address", "imei", "udid",

	// Social & Professional
	"company", "employer", "organization", "org",
	"jobtitle", "job_title", "position", "role",
	"linkedin", "twitter", "facebook", "instagram", "social",
}

// ipv4Regex matches IPv4 addresses for anonymization
var ipv4Regex = regexp.MustCompile(`\b(\d{1,3}\.\d{1,3}\.\d{1,3})\.\d{1,3}\b`)

// ipv6Regex matches IPv6 addresses for anonymization
var ipv6Regex = regexp.MustCompile(`\b([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}\b`)

// emailRegex matches email addresses for masking
var emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

// localTimeEncoder encodes time in local timezone with ISO8601 format
// Example: 2025-12-03T10:54:27+05:30 (IST) or 2025-12-03T06:24:27+01:00 (CET)
func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Local().Format(time.RFC3339))
}

// ============================================
// ASYNC WRITE SYNCER
// ============================================

// asyncWriteSyncer wraps a WriteSyncer to make it asynchronous
type asyncWriteSyncer struct {
	ws       zapcore.WriteSyncer
	queue    chan []byte
	done     chan struct{}
	wg       sync.WaitGroup
	stopOnce sync.Once
}

// newAsyncWriteSyncer creates a new asynchronous WriteSyncer
func newAsyncWriteSyncer(ws zapcore.WriteSyncer, bufferSize int) *asyncWriteSyncer {
	aws := &asyncWriteSyncer{
		ws:    ws,
		queue: make(chan []byte, bufferSize),
		done:  make(chan struct{}),
	}

	aws.wg.Add(1)
	go aws.run()

	return aws
}

// run processes log entries from the queue
func (aws *asyncWriteSyncer) run() {
	defer aws.wg.Done()
	for {
		select {
		case entry := <-aws.queue:
			_, _ = aws.ws.Write(entry)
		case <-aws.done:
			// Drain remaining entries
			for {
				select {
				case entry := <-aws.queue:
					_, _ = aws.ws.Write(entry)
				default:
					return
				}
			}
		}
	}
}

// Write implements io.Writer - non-blocking, returns immediately
func (aws *asyncWriteSyncer) Write(p []byte) (n int, err error) {
	// Make a copy of the data to avoid race conditions
	entry := make([]byte, len(p))
	copy(entry, p)

	select {
	case aws.queue <- entry:
		return len(p), nil
	default:
		// Queue is full, fallback to synchronous write to prevent blocking
		return aws.ws.Write(p)
	}
}

// Sync flushes all pending writes
func (aws *asyncWriteSyncer) Sync() error {
	aws.stopOnce.Do(func() {
		close(aws.done)
	})
	aws.wg.Wait()
	return aws.ws.Sync()
}

// ============================================
// LOGGER CONFIGURATION
// ============================================

type LoggerConfig struct {
	LogLevel   string
	Compliance *ComplianceConfig
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

		// Set compliance config
		if config.Compliance != nil {
			complianceConfig = config.Compliance
		} else {
			complianceConfig = DefaultComplianceConfig()
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
			EncodeTime:     localTimeEncoder, // Uses server's local timezone (IST, CET, etc.)
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		})

		// Create async writer with buffer size of 256KB (can handle ~1000 log entries)
		asyncWriter = newAsyncWriteSyncer(zapcore.AddSync(os.Stdout), 256*1024)
		core := zapcore.NewCore(jsonEncoder, asyncWriter, logLevel)

		defaultFields := []zap.Field{}
		addEnvironmentFields(&defaultFields)

		// Add compliance metadata
		if complianceConfig.DataRegion != "" {
			defaultFields = append(defaultFields, zap.String("data_region", complianceConfig.DataRegion))
		}
		defaultFields = append(defaultFields, zap.String("compliance_mode", string(complianceConfig.Mode)))

		logInstance = zap.New(core, zap.AddCaller()).With(defaultFields...)

		logInstance.Info("Logger initialized",
			zap.String("default_level", logLevel.String()),
			zap.String("compliance_mode", string(complianceConfig.Mode)),
			zap.Bool("ip_anonymization", complianceConfig.AnonymizeIP),
			zap.Int("retention_days", complianceConfig.RetentionDays),
		)
	})

	return errInit
}

// SetComplianceMode allows runtime compliance mode changes
func SetComplianceMode(mode ComplianceMode) {
	if complianceConfig != nil {
		complianceConfig.Mode = mode
		logInstance.Info("Compliance mode updated", zap.String("new_mode", string(mode)))
	}
}

// GetComplianceConfig returns current compliance configuration
func GetComplianceConfig() *ComplianceConfig {
	return complianceConfig
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
	if asyncWriter != nil {
		_ = asyncWriter.Sync()
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
	addIfNotEmpty(fields, "env", os.Getenv("ENVIRONMENT"))
	addIfNotEmpty(fields, "service", os.Getenv("SERVICE_NAME"))
	addIfNotEmpty(fields, "version", os.Getenv("SERVICE_VERSION"))
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

// ============================================
// AUDIT LOGGING (SOC 2, ISO 27001)
// ============================================

// Audit logs security-relevant events with compliance metadata
func Audit(action string, args ...interface{}) {
	ensureInitialized()
	// Audit logs always go through regardless of log level
	auditArgs := append([]interface{}{
		"audit", true,
		"action", action,
		"audit_ts", time.Now().UTC().Format(time.RFC3339Nano),
	}, args...)
	logWithLevel("info", "AUDIT: "+action, auditArgs...)
}

// AuditAccess logs data access events (GDPR Art. 30)
func AuditAccess(userID, resource, action string, args ...interface{}) {
	auditArgs := append([]interface{}{
		"audit_type", "data_access",
		"accessor_id", processValue("userid", userID),
		"resource", resource,
		"access_action", action,
	}, args...)
	Audit("DATA_ACCESS", auditArgs...)
}

// AuditConsent logs consent-related events (GDPR Art. 7)
func AuditConsent(userID, consentType, status string, args ...interface{}) {
	auditArgs := append([]interface{}{
		"audit_type", "consent",
		"subject_id", processValue("userid", userID),
		"consent_type", consentType,
		"consent_status", status,
	}, args...)
	Audit("CONSENT_"+strings.ToUpper(status), auditArgs...)
}

// AuditDataSubjectRequest logs GDPR data subject requests
func AuditDataSubjectRequest(requestType, subjectID, status string, args ...interface{}) {
	auditArgs := append([]interface{}{
		"audit_type", "dsr",
		"request_type", requestType, // access, rectification, erasure, portability
		"subject_id", processValue("userid", subjectID),
		"request_status", status,
	}, args...)
	Audit("DSR_"+strings.ToUpper(requestType), auditArgs...)
}

func stackTraceField() zap.Field {
	return zap.String("stacktrace", string(debug.Stack()))
}

// ============================================
// DATA PROTECTION FUNCTIONS
// ============================================

func withFields(args []interface{}) []zap.Field {
	fields := []zap.Field{}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			switch key := args[i].(type) {
			case string:
				value := args[i+1]

				if value == nil {
					continue
				}

				processedValue := processValue(key, value)
				if processedValue == nil {
					continue
				}

				switch v := processedValue.(type) {
				case string:
					if v == "" && isOptionalField(key) {
						continue
					}
					fields = append(fields, zap.String(key, v))
				case time.Duration:
					fields = append(fields, zap.String(key, v.String()))
				case time.Time:
					fields = append(fields, zap.Time(key, v))
				default:
					fields = append(fields, zap.Any(key, v))
				}
			default:
				if key != nil {
					logInstance.WithOptions(zap.AddCallerSkip(3)).Warn("Invalid key in arguments passed to logger", zap.Any("invalid_key", key))
				}
			}
		} else {
			if key, ok := args[i].(string); ok {
				logInstance.WithOptions(zap.AddCallerSkip(4)).Warn("Missing value for key in arguments passed to logger", zap.String("key", key))
			}
		}
	}
	return fields
}

// processValue applies appropriate protection based on key type and compliance mode
func processValue(key string, value interface{}) interface{} {
	if value == nil {
		return nil
	}

	keyLower := strings.ToLower(key)

	// Handle string values
	if strVal, ok := value.(string); ok {
		// Sensitive data - always redact
		if containsAny(keyLower, sensitiveKeys) {
			return "[REDACTED]"
		}

		// Personal data - mask based on compliance mode
		if containsAny(keyLower, personalKeys) {
			return maskPII(keyLower, strVal)
		}

		// Check for embedded PII in values (defense in depth)
		return sanitizeValue(strVal)
	}

	return value
}

// maskPII applies appropriate masking based on data type and compliance mode
func maskPII(keyLower, value string) string {
	if complianceConfig == nil || value == "" {
		return mask(value)
	}

	// In dev mode with compliance disabled, return value as-is for easier debugging
	if !complianceConfig.Enabled {
		return value
	}

	switch complianceConfig.Mode {
	case ComplianceModeStrict:
		// Complete redaction for strict mode
		return "[PII_REDACTED]"

	case ComplianceModeModerate:
		// Smart masking based on data type
		if isIPField(keyLower) && complianceConfig.AnonymizeIP {
			return anonymizeIP(value)
		}
		if isEmailField(keyLower) {
			return maskEmail(value)
		}
		if isPhoneField(keyLower) {
			return maskPhone(value)
		}
		if isUserIDField(keyLower) && complianceConfig.PseudonymizeUserID {
			return pseudonymize(value)
		}
		return mask(value)

	case ComplianceModeMinimal:
		// Basic masking
		return mask(value)

	default:
		return mask(value)
	}
}

// sanitizeValue checks for embedded PII in arbitrary strings
func sanitizeValue(value string) string {
	result := value

	// Anonymize embedded IP addresses
	if complianceConfig != nil && complianceConfig.AnonymizeIP {
		// IPv4: Replace last octet with xxx
		result = ipv4Regex.ReplaceAllString(result, "$1.xxx")
		// IPv6: Replace with anonymized version
		result = ipv6Regex.ReplaceAllStringFunc(result, func(ip string) string {
			return anonymizeIP(ip)
		})
	}

	// Mask embedded email addresses
	result = emailRegex.ReplaceAllStringFunc(result, func(email string) string {
		return maskEmail(email)
	})

	return result
}

// ============================================
// MASKING & ANONYMIZATION HELPERS
// ============================================

// mask provides basic character masking
func mask(input string) string {
	if len(input) <= 2 {
		return "**"
	}
	return input[:1] + strings.Repeat("*", len(input)-2) + input[len(input)-1:]
}

// maskEmail masks email while preserving domain hint
// j***n@e***.com
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return mask(email)
	}

	localPart := parts[0]
	domain := parts[1]

	maskedLocal := mask(localPart)

	// Mask domain but keep TLD
	domainParts := strings.Split(domain, ".")
	if len(domainParts) >= 2 {
		maskedDomain := mask(domainParts[0])
		tld := domainParts[len(domainParts)-1]
		return maskedLocal + "@" + maskedDomain + "." + tld
	}

	return maskedLocal + "@" + mask(domain)
}

// maskPhone masks phone number keeping country code and last 2 digits
// +91******90
func maskPhone(phone string) string {
	// Remove non-digit characters for processing
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(digits) < 4 {
		return "****"
	}

	// Keep first 2-3 digits (country code) and last 2 digits
	prefixLen := 2
	if len(digits) > 10 {
		prefixLen = len(digits) - 10 + 2 // Country code + 2
	}

	prefix := digits[:prefixLen]
	suffix := digits[len(digits)-2:]
	maskLen := len(digits) - prefixLen - 2

	if strings.HasPrefix(phone, "+") {
		return "+" + prefix + strings.Repeat("*", maskLen) + suffix
	}
	return prefix + strings.Repeat("*", maskLen) + suffix
}

// anonymizeIP anonymizes IP address by zeroing last octet (GDPR compliant)
// 192.168.1.100 -> 192.168.1.0
func anonymizeIP(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return mask(ip)
	}

	// IPv4
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		ipv4[3] = 0
		return ipv4.String()
	}

	// IPv6 - zero last 80 bits (keep /48)
	if len(parsedIP) == 16 {
		for i := 6; i < 16; i++ {
			parsedIP[i] = 0
		}
		return parsedIP.String()
	}

	return mask(ip)
}

// pseudonymize creates a consistent hash for user tracking without exposing ID
func pseudonymize(value string) string {
	if complianceConfig == nil || complianceConfig.Salt == "" {
		return mask(value)
	}

	hash := sha256.Sum256([]byte(value + complianceConfig.Salt))
	return "pseudo_" + hex.EncodeToString(hash[:8]) // First 16 chars of hash
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func containsAny(item string, slice []string) bool {
	for _, s := range slice {
		if strings.Contains(item, s) {
			return true
		}
	}
	return false
}

func isOptionalField(key string) bool {
	keyLower := strings.ToLower(key)
	return keyLower == "userid" || keyLower == "user_id" || keyLower == "tenantid" || keyLower == "tenant_id"
}

func isIPField(key string) bool {
	return strings.Contains(key, "ip") || strings.Contains(key, "address")
}

func isEmailField(key string) bool {
	return strings.Contains(key, "email") || strings.Contains(key, "mail")
}

func isPhoneField(key string) bool {
	return strings.Contains(key, "phone") || strings.Contains(key, "mobile") ||
		strings.Contains(key, "tel") || strings.Contains(key, "cell") ||
		strings.Contains(key, "whatsapp")
}

func isUserIDField(key string) bool {
	return strings.Contains(key, "userid") || strings.Contains(key, "user_id") ||
		key == "id" || strings.Contains(key, "subject")
}
