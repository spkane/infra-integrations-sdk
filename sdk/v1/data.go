package v1

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/cache"
	"github.com/newrelic/infra-integrations-sdk/metric"
)

type inventoryItem map[string]interface{}

// Inventory is the data type for inventory data produced by an integration data
// source and emitted to the agent's inventory data store.
type Inventory map[string]inventoryItem

// SetItem stores a value into the inventory data structure.
func (i Inventory) SetItem(key string, field string, value interface{}) {
	if _, ok := i[key]; ok {
		i[key][field] = value
	} else {
		i[key] = inventoryItem{field: value}
	}

}

// Event is the data type to represent arbitrary, one-off messages for key
// activities on a system.
type Event struct {
	Summary  string `json:"summary"`
	Category string `json:"category,omitempty"`
}

// DefaultEventCategory is a default value for event category
const DefaultEventCategory = "notifications"

// Integration defines the format of the output JSON that integrations will
// return.
// Deprecated: use v2.Integration
type Integration struct {
	Name               string       `json:"name"`
	ProtocolVersion    string       `json:"protocol_version"`
	IntegrationVersion string       `json:"integration_version"`
	Metrics            []metric.Set `json:"metrics"`
	Inventory          Inventory    `json:"inventory"`
	Events             []Event      `json:"events"`
	prettyOutput       bool
}

// NewIntegration initializes a new instance of integration data.
// Deprecated: use v2.NewIntegration
func NewIntegration(name string, version string, arguments interface{}) (*Integration, error) {
	err := args.SetupArgs(arguments)
	if err != nil {
		return nil, err
	}
	defaultArgs := args.GetDefaultArgs(arguments)

	cache.GlobalLog.SetDebug(defaultArgs.Verbose)

	// Avoid working with an uninitialized or in error state cache
	if err = cache.Status(); err != nil {
		return nil, err
	}

	integration := &Integration{
		Name:               name,
		ProtocolVersion:    "1",
		IntegrationVersion: version,
		Inventory:          make(Inventory),
		Metrics:            make([]metric.Set, 0),
		Events:             make([]Event, 0),
		prettyOutput:       defaultArgs.Pretty,
	}

	return integration, nil
}

// NewMetricSet returns a new instance of Set with its sample attached to
// the IntegrationData.
func (integration *Integration) NewMetricSet(eventType string) metric.Set {
	ms := metric.NewSet(eventType)
	integration.Metrics = append(integration.Metrics, ms)
	return ms
}

// AddNotificationEvent method adds a new Event with default event category.
func (integration *Integration) AddNotificationEvent(summary string) error {
	return integration.AddEvent(Event{Summary: summary, Category: DefaultEventCategory})
}

// AddEvent method adds a new Event.
func (integration *Integration) AddEvent(e Event) error {
	if e.Summary == "" {
		return fmt.Errorf("summary of the event cannot be empty")
	}

	integration.Events = append(integration.Events, e)

	return nil
}

// Publish runs all necessary tasks before publishing the data. Currently, it
// stores the cache, prints the JSON representation of the integration to stdout
// and re-initializes the integration object (allowing re-use it during the
// execution of your code).
func (integration *Integration) Publish() error {
	if err := cache.Save(); err != nil {
		return err
	}

	output, err := integration.toJSON(integration.prettyOutput)
	if err != nil {
		return err
	}

	fmt.Println(output)
	integration.Clear()

	return nil
}

// Clear re-initializes the Inventory, Metrics and Events for this integration.
// Used after publishing so the object can be reused.
func (integration *Integration) Clear() {
	integration.Inventory = make(Inventory)
	integration.Metrics = make([]metric.Set, 0)
	integration.Events = make([]Event, 0)
}

// toJSON returns the integration as a JSON string. If the pretty attribute is
// set to true, the JSON will be indented for easy reading.
func (integration *Integration) toJSON(pretty bool) (string, error) {
	var output []byte
	var err error

	if pretty {
		output, err = json.MarshalIndent(integration, "", "\t")
	} else {
		output, err = json.Marshal(integration)
	}

	if err != nil {
		return "", fmt.Errorf("error marshalling to JSON: %s", err)
	}

	if string(output) == "null" {
		return "[]", nil
	}

	return string(output), nil
}
