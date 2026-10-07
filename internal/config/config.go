package config

import (
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Log  Log
	HTTP HTTP
	DB   DB
}

type Log struct {
	Level string
}

type HTTP struct {
	Address           string
	TrustedProxies    []netip.Prefix
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

type DB struct {
	Address         string
	MaxConns        *int32
	MaxConnLifetime *time.Duration
	MaxConnIdleTime *time.Duration
	MinConns        *int32
}

func InitConfig() (Config, error) {
	proxies, err := trustedProxies(envOrDefault("HTTP_TRUSTED_PROXIES", ""))
	if err != nil {
		return Config{}, err
	}
	maxConns, err := optionalInt("MAX_CONNS")
	if err != nil {
		return Config{}, err
	}
	minConns, err := optionalInt("MIN_CONNS")
	if err != nil {
		return Config{}, err
	}
	maxConnLifetime, err := optionalDuration("MAX_CONN_LIFETIME")
	if err != nil {
		return Config{}, err
	}
	maxConnIdleTime, err := optionalDuration("MAX_CONN_IDLE_TIME")
	if err != nil {
		return Config{}, err
	}
	readHeaderTimeout, err := durationOrDefault("HTTP_READ_HEADER_TIMEOUT", 2*time.Second)
	if err != nil {
		return Config{}, err
	}
	readTimeout, err := durationOrDefault("HTTP_READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := durationOrDefault("HTTP_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	idleTimeout, err := durationOrDefault("HTTP_IDLE_TIMEOUT", 30*time.Second)
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := durationOrDefault("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Log: Log{
			Level: envOrDefault("LOG_LEVEL", "DEBUG"),
		},
		HTTP: HTTP{
			Address:           envOrDefault("HTTP_ADDRESS", "0.0.0.0:8080"),
			TrustedProxies:    proxies,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
			ShutdownTimeout:   shutdownTimeout,
		},
		DB: DB{
			Address:         envOrDefault("DATABASE_URL", "postgres://postgres:postgres@db:5432/postgres"),
			MaxConns:        maxConns,
			MinConns:        minConns,
			MaxConnLifetime: maxConnLifetime,
			MaxConnIdleTime: maxConnIdleTime,
		},
	}
	return cfg, nil
}

func optionalInt(env string) (*int32, error) {
	value := strings.TrimSpace(os.Getenv(env))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", env, err)
	}
	res := int32(parsed)
	return &res, nil
}

func optionalDuration(env string) (*time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(env))
	if value == "" {
		return nil, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", env, err)
	}
	return &parsed, nil
}

func trustedProxies(value string) ([]netip.Prefix, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	values := strings.Split(value, ",")
	proxies := make([]netip.Prefix, 0, len(values))
	for _, value := range values {
		proxy, err := netip.ParsePrefix(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("parse HTTP_TRUSTED_PROXIES %q: %w", value, err)
		}
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}

func envOrDefault(env string, def string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		return strings.TrimSpace(def)
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return strings.TrimSpace(def)
	}
	return value
}

func durationOrDefault(env string, def time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(env))
	if value == "" {
		return def, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", env, err)
	}
	return parsed, nil
}
