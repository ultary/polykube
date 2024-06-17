package helm

import "ultary.co/kluster/pkg/monokube/otlp"

func (c *Chart) NewOpenTelemetryAgent() *otlp.OpenTelemetry {
	const name = "otel_agent"
	d := c.dependencies[name]
	v := c.values[name].(map[string]interface{})
	return otlp.NewOpenTelemetry(d, v)
}

func (c *Chart) NewOpenTelemetryCollector() *otlp.OpenTelemetry {
	const name = "otel_collector"
	d := c.dependencies[name]
	v := c.values[name].(map[string]interface{})
	return otlp.NewOpenTelemetry(d, v)
}

func (c *Chart) NewOpenTelemetryConsumer() *otlp.OpenTelemetry {
	const name = "otel_consumer"
	d := c.dependencies[name]
	v := c.values[name].(map[string]interface{})
	return otlp.NewOpenTelemetry(d, v)
}
