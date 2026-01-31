package plugins

import (
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// Plugin defines the interface that all custom analysis plugins must implement
type Plugin interface {
	// Name returns the unique name of the plugin
	Name() string

	// Version returns the version of the plugin
	Version() string

	// Description returns a human-readable description of what the plugin does
	Description() string

	// Analyze performs the custom analysis and returns results
	Analyze(repo *github.Repo, commits []github.Commit, contributors []github.Contributor) (*PluginResult, error)
}

// PluginResult contains the results of a plugin's analysis
type PluginResult struct {
	// Score is a numerical score (0-100) contributed by this plugin
	Score int `json:"score"`

	// Weight determines how much this plugin's score affects the overall score (0.0-1.0)
	Weight float64 `json:"weight"`

	// Metrics contains key-value pairs of custom metrics
	Metrics map[string]interface{} `json:"metrics"`

	// Recommendations contains actionable suggestions from the plugin
	Recommendations []string `json:"recommendations"`

	// RiskLevel indicates the risk level determined by this plugin ("Low", "Medium", "High")
	RiskLevel string `json:"risk_level"`

	// Category groups the plugin (e.g., "Security", "Performance", "Code Quality")
	Category string `json:"category"`
}

// PluginRegistry manages loaded plugins
type PluginRegistry struct {
	plugins map[string]Plugin
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]Plugin),
	}
}

// Register adds a plugin to the registry
func (r *PluginRegistry) Register(plugin Plugin) {
	r.plugins[plugin.Name()] = plugin
}

// Get retrieves a plugin by name
func (r *PluginRegistry) Get(name string) (Plugin, bool) {
	plugin, exists := r.plugins[name]
	return plugin, exists
}

// List returns all registered plugins
func (r *PluginRegistry) List() map[string]Plugin {
	return r.plugins
}

// Unregister removes a plugin from the registry
func (r *PluginRegistry) Unregister(name string) {
	delete(r.plugins, name)
}
