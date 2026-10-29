// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package serverboot

import "testing"

func TestOTLPDestinationDefaults(t *testing.T) {
	for _, test := range []struct {
		name, endpoint, metrics, traces, exporter string
		connected, wantMetrics, wantTraces        bool
	}{
		{name: "absent destination"},
		{name: "whitespace destination", endpoint: " \t "},
		{name: "shared destination", endpoint: "http://collector.test:4317", wantMetrics: true, wantTraces: true},
		{name: "metrics destination", metrics: "http://collector.test:4317", wantMetrics: true},
		{name: "traces destination", traces: "http://collector.test:4317", wantTraces: true},
		{name: "relay", connected: true, wantMetrics: true, wantTraces: true},
		{name: "explicit disabled destination", endpoint: "http://collector.test:4317", exporter: "none"},
		{name: "explicit disabled relay", connected: true, exporter: "none"},
		{name: "explicit SDK destination", exporter: "otlp", wantMetrics: true, wantTraces: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", test.endpoint)
			t.Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", test.metrics)
			t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", test.traces)
			t.Setenv(metricsExporterEnv, test.exporter)
			t.Setenv(tracesExporterEnv, test.exporter)
			if got := metricsExporters(t.Context(), false, test.connected).Has(ExporterOTLP); got != test.wantMetrics {
				t.Fatalf("metrics push = %t, want %t", got, test.wantMetrics)
			}
			if got := tracesPushEnabled(t.Context(), test.connected); got != test.wantTraces {
				t.Fatalf("traces push = %t, want %t", got, test.wantTraces)
			}
		})
	}
}
