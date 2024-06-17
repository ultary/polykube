package repo

import "ultary.co/kluster/pkg/monokube/otlp"

func (m *Manifests) NewOpenTelemetryAgent() *otlp.OpenTelemetry {
	const name = "otel_agent"
	d := m.dependencies[name]
	v := m.values[name].(map[string]interface{})
	return otlp.NewOpenTelemetry(d, v)
}

func (m *Manifests) NewOpenTelemetryCollector() *otlp.OpenTelemetry {
	const name = "otel_collector"
	d := m.dependencies[name]
	v := m.values[name].(map[string]interface{})
	return otlp.NewOpenTelemetry(d, v)
}

func (m *Manifests) NewOpenTelemetryConsumer() *otlp.OpenTelemetry {
	const name = "otel_consumer"
	d := m.dependencies[name]
	v := m.values[name].(map[string]interface{})
	return otlp.NewOpenTelemetry(d, v)
}
