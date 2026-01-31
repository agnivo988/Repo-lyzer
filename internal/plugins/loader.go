package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

// Loader handles loading and managing plugins
type Loader struct {
	registry *PluginRegistry
}

// NewLoader creates a new plugin loader
func NewLoader() *Loader {
	return &Loader{
		registry: NewPluginRegistry(),
	}
}

// LoadPlugins scans the specified directory for plugin binaries and loads them
func (l *Loader) LoadPlugins(directory string) error {
	if directory == "" {
		return fmt.Errorf("plugin directory not specified")
	}

	// Check if directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return fmt.Errorf("plugin directory does not exist: %s", directory)
	}

	// Get plugin file extension based on OS
	ext := l.getPluginExtension()

	// Walk through the directory
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-plugin files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ext) {
			return nil
		}

		// Load the plugin
		if err := l.loadPlugin(path); err != nil {
			// Log error but continue loading other plugins
			fmt.Printf("Warning: failed to load plugin %s: %v\n", path, err)
			return nil
		}

		fmt.Printf("Loaded plugin: %s\n", path)
		return nil
	})
}

// loadPlugin loads a single plugin file
func (l *Loader) loadPlugin(path string) error {
	// Open the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look for the Plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin does not export 'Plugin' symbol: %w", err)
	}

	// Assert that it's a Plugin
	pluginInstance, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin does not implement Plugin interface")
	}

	// Register the plugin
	l.registry.Register(pluginInstance)
	return nil
}

// getPluginExtension returns the appropriate plugin file extension for the current OS
func (l *Loader) getPluginExtension() string {
	switch runtime.GOOS {
	case "windows":
		return ".dll"
	case "darwin":
		return ".dylib"
	default: // linux and others
		return ".so"
	}
}

// GetRegistry returns the plugin registry
func (l *Loader) GetRegistry() *PluginRegistry {
	return l.registry
}

// GetLoadedPlugins returns a list of loaded plugin names
func (l *Loader) GetLoadedPlugins() []string {
	plugins := l.registry.List()
	names := make([]string, 0, len(plugins))
	for name := range plugins {
		names = append(names, name)
	}
	return names
}

// IsPluginLoaded checks if a specific plugin is loaded
func (l *Loader) IsPluginLoaded(name string) bool {
	_, exists := l.registry.Get(name)
	return exists
}

// UnloadPlugin removes a plugin from the registry
func (l *Loader) UnloadPlugin(name string) {
	l.registry.Unregister(name)
}
