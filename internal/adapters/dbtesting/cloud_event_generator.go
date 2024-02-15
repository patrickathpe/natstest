// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package dbtesting

import (
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

const (
	PropertyCustomerID   = "customerid"
	PropertyTraceParent  = "traceparent"
	PropertyTraceState   = "tracestate"
	PropertyDevice       = "device"
	PropertyOutboxID     = "outboxid"
	PropertyPartitionKey = "partitionkey"
)

type CloudEventGenerator struct {
	// Headers
	ID           uuid.UUID
	AuthHeader   string
	CustomerID   string
	DataSchema   string
	DeviceID     string
	PartitionKey string
	Source       string
	SpecVersion  string
	Subject      string
	Time         time.Time
	TraceParent  string
	TraceState   string
	Type         string

	// Data payload
	Data            interface{}
	DataContentType string
}

func (gen *CloudEventGenerator) InitializeUnsetFields() {
	gen.initializeUnsetHeaders()
	gen.initializeUnsetData()
}

func (gen *CloudEventGenerator) initializeUnsetHeaders() {
	if gen.ID == uuid.Nil {
		gen.ID = uuid.New()
	}
	if gen.AuthHeader == "" {
		gen.AuthHeader = RandomString(generatedIDLen)
	}
	if gen.CustomerID == "" {
		gen.CustomerID = RandomString(generatedIDLen)
	}
	if gen.DataSchema == "" {
		gen.DataSchema = RandomString(generatedIDLen)
	}
	if gen.DeviceID == "" {
		gen.DeviceID = RandomString(generatedIDLen)
	}
	if gen.PartitionKey == "" {
		gen.PartitionKey = RandomString(generatedIDLen)
	}
	if gen.Source == "" {
		gen.Source = RandomString(generatedIDLen)
	}
	if gen.SpecVersion == "" {
		gen.SpecVersion = ce.VersionV1
	}
	if gen.Subject == "" {
		gen.Subject = RandomString(generatedIDLen)
	}
	if gen.Time.IsZero() {
		gen.Time = time.Now().UTC()
	}
	if gen.TraceParent == "" {
		gen.TraceParent = RandomString(generatedNameLen)
	}
	if gen.TraceState == "" {
		gen.TraceState = RandomString(generatedNameLen)
	}
	if gen.Type == "" {
		gen.Type = RandomString(generatedIDLen)
	}
}

func (gen *CloudEventGenerator) initializeUnsetData() {
	if gen.DataContentType == "" && gen.Data == nil {
		gen.DataContentType = ce.ApplicationJSON
		gen.Data = map[string]string{
			"foo": RandomString(generatedIDLen),
			"bar": RandomString(generatedIDLen),
		}
	}
}

// GetCloudEvent returns a CloudEvent entity containing the generated properties
func (gen *CloudEventGenerator) GetCloudEvent() *ce.Event {
	gen.InitializeUnsetFields()
	event := ce.NewEvent(gen.SpecVersion)
	event.SetDataSchema(gen.DataSchema)
	event.SetID(gen.ID.String())
	event.SetSource(gen.Source)
	event.SetSubject(gen.Subject)
	event.SetTime(gen.Time)
	event.SetType(gen.Type)
	event.SetExtension("authheader", gen.AuthHeader)
	event.SetExtension(PropertyCustomerID, gen.CustomerID)
	event.SetExtension(PropertyDevice, gen.DeviceID)
	event.SetExtension(PropertyPartitionKey, gen.PartitionKey)
	event.SetExtension(PropertyTraceParent, gen.TraceParent)
	event.SetExtension(PropertyTraceState, gen.TraceState)
	_ = event.SetData(gen.DataContentType, gen.Data)
	event.DataBase64 = false
	return &event
}
