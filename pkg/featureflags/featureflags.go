/*
Package featureflags provides support to set up and look up of cluster wide features.

Features are supposed to be have unique global names containing the
same symbols allowed for environment variables. Featureflags are
orginized as name=value pairs where value is primarily designed to
be strings "Enabled" or "Disabled"

The package provides functionality to split features by categories
(alpha/beta) and enable/disable by category.

inspired by:

https://pkg.go.dev/k8s.io/component-base/featuregate

https://github.com/knative/eventing/blob/main/pkg/apis/feature/features.go

All the known features are listed in the array `features` of
FeatureSpec. Unknown features are ignored and not preserved in the
configuration.

The actual feature configuration is internal global state of the
package. The config is updated automically whole at once. Individual
feature changes are not supposed to happen.

Package interface functions are thread-safe, although there is room to
read stale values (especially for Subscriber) if several configuration
updates come too fast. Subscriber should get update signal for every
config update, but it's recommended to get the actual value in the
handler.

!!! On the initialization package creates default configuration with default values.

It means, that the client code should be either ready to change of
feature value in runtime or should not check feature values before
proper update happend.

The package implement simple Publisher/Subscriber interface to track
updates of individual features.

# Possible scenarios

There are different possibilities to use the package.

1. features can be defined in a ConfigMap or some CR and
unmarshaled directly to a map[string]string. The client can watch
the object and for any object update call
UpdateConfigFromMap(). After that it should trigger reevaluation of
the code which depends of the feature value.

2. The client code which depends of the feature subscribes on
feature changes. As above features can be defined in a ConfigMap or
some CR. As soon as feature value updated the Subscriber gets
signal and reavaluate the code.

3. feature configuration initialized from environment by calling
UpdateConfigFromEnvironment() on startup, after that code, which
depends of feature values evaluated.

*/
package featureflags

import (
	"os"
	"slices"
	"sync"
	"sync/atomic"
)

const (
	Enabled = "Enabled"
	Disabled = "Disabled"
)

var (
	allowed = []string{ Enabled, Disabled, }
	current atomic.Value // stores singleton config
	subs = subscribers{
		subs: map[string][]Subscriber{},
	}
)

type subscribers struct {
	m sync.Mutex
	subs map[string][]Subscriber // maps feature name to subscribers
}

type featureSpec struct {
	Default string
	Since string
	Stage string
}

type Subscriber interface {
	Changed(name, old, new string) // signals that feature changed its value
}

// features array stores all defined features and their characteristics
var features = map[string]featureSpec{
	"EXAMPLE": { Default: Enabled, Since: "2.0", Stage: "alpha", },
}

type Config map[string]string

func init() {
	config := newDefaults()
	current.Store(config)
}

func newDefaults() Config {
	config := make(Config)
	for k,v := range features {
		config[k] = v.Default
	}
	return config
}

func signal(f, old, new string) {
	subs.m.Lock()
	defer subs.m.Unlock()

	arr := subs.subs[f]

	for _, s := range arr {
		s.Changed(f, old, new)
	}
}

func updateConfig(c Config) {
	old := current.Swap(c).(Config)
	for k, v := range c {
		if v != old[k] {
			signal(k, old[k], v)
		}
	}
}

func updateConfigFromEnvironment() Config {
	c := newDefaults()
	for k, _ := range features {
		envV := os.Getenv(k)
		if envV != "" {
			c[k] = envV
		}
	}
	return c
}

func UpdateConfigFromEnvironment() {
	config := updateConfigFromEnvironment()
	updateConfig(config)
}

// UpdateConfigFromMap creates config based on the provided map and
// replaces current config with it. It uses default values if any
// feature is not part of the map
//
// The function is supposed to be called from any Reconcile to update
// the feature config at once atomically
func UpdateConfigFromMap(m map[string]string) {
	config := updateConfigFromMap(m)
	updateConfig(config)
}

func updateConfigFromMap(m map[string]string) Config {
	config := newDefaults()

	for k,v := range m {
		_, exists := features[k]
		if !exists {
			continue
		}
		if !slices.Contains(allowed, v) {
			continue
		}
		config[k] = v
	}

	return config
}

func getConfig() Config {
	return current.Load().(Config)
}

func (c Config) IsEnabled(f string) bool {
	return c[f] == Enabled
}

func (c Config) IsDisabled(f string) bool {
	return c[f] != Enabled
}

// IsEnabled is a wrapper on (Config).IsEnabled
// to make access like featureflags.IsEnabled("EXAMPLE")
func IsEnabled(f string) bool {
	c := getConfig()
	return c.IsEnabled(f)
}

// IsDisabled is a wrapper on (Config).IsDisabled
// to make access like featureflags.IsDisabled("EXAMPLE")
func IsDisabled(f string) bool {
	c := getConfig()
	return c.IsDisabled(f)
}

// Subscribe subscribes s on changes of feature named f
func Subscribe(s Subscriber, f string) {
	subs.m.Lock()
	defer subs.m.Unlock()

	arr := subs.subs[f]
	if !slices.Contains(arr, s) {
		arr = append(arr, s)
	}
	subs.subs[f] = arr
}
