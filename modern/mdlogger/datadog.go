package mdlogger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/short-d/app/fw"
)

// DataDog logging API =>
//   https://docs.datadoghq.com/api/?lang=bash#logs
const dataDogLoggingApi = "https://http-intake.logs.datadoghq.com/v1/input"

var _ EntryRepository = (*DataDogEntryRepo)(nil)

type DataDogEntryRepo struct {
	apiKey      string
	httpRequest fw.HTTPRequest
	env         fw.ServerEnv
}

func (d DataDogEntryRepo) createLogEntry(
	level fw.LogLevel,
	prefix string,
	line int,
	filename string,
	message string,
	date time.Time) {
	headers := d.authHeaders()

	body := d.requestBody(level, prefix, line, filename, message)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	var res = make(map[string]interface{})
	_ = d.httpRequest.JSON(http.MethodPost, dataDogLoggingApi, headers, string(jsonBody), &res)
}

func getSeverity(level fw.LogLevel) string {
	switch level {
	case fw.LogFatal:
		return "critical"
	case fw.LogError:
		return "error"
	case fw.LogWarn:
		return "warning"
	case fw.LogInfo:
		return "info"
	case fw.LogDebug:
		return "debug"
	case fw.LogTrace:
		return "debug"
	default:
		return "Should not happen"
	}
}

func (d DataDogEntryRepo) requestBody(
	level fw.LogLevel,
	prefix string,
	line int,
	filename string,
	message string) map[string]string {
	severity := getSeverity(level)
	tags := map[string]string{
		"env":       string(d.env),
		"line":      fmt.Sprintf("%d", line),
		"file-name": filename,
	}
	return map[string]string{
		"service": prefix,
		"status":  severity,
		"message": message,
		"ddtags":  d.dataDogTags(tags),
	}
}

func (d DataDogEntryRepo) dataDogTags(tags map[string]string) string {
	var tagsList []string

	for key, val := range tags {
		pair := fmt.Sprintf("%s:%s", key, val)
		tagsList = append(tagsList, pair)
	}
	return strings.Join(tagsList, ",")
}

func (d DataDogEntryRepo) authHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"DD-API-KEY":   d.apiKey,
	}
}

func NewDataDogEntryRepo(apiKey string, httpRequest fw.HTTPRequest, env fw.ServerEnv) DataDogEntryRepo {
	return DataDogEntryRepo{
		apiKey:      apiKey,
		httpRequest: httpRequest,
		env:         env,
	}
}
