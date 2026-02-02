package cmd

import (
	"fmt"
	"os"

	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/plugins"
	"github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage custom analysis plugins",
	Long:  `Manage custom analysis plugins for extended repository analysis capabilities.`,
}

var pluginsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available and loaded plugins",
	Long:  `List all plugins found in the plugin directory and their current status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := config.LoadSettings()
		if err != nil {
			return fmt.Errorf("failed to load settings: %w", err)
		}

		loader := plugins.NewLoader()
		err = loader.LoadPlugins(settings.PluginDirectory)
		if err != nil {
			return fmt.Errorf("failed to load plugins: %w", err)
		}

		loadedPlugins := loader.GetLoadedPlugins()
		enabledPlugins := settings.EnabledPlugins

		fmt.Println("📦 Available Plugins:")
		if len(loadedPlugins) == 0 {
			fmt.Println("  No plugins found in directory:", settings.PluginDirectory)
			return nil
		}

		for _, pluginName := range loadedPlugins {
			status := "❌ Disabled"
			if contains(enabledPlugins, pluginName) {
				status = "✅ Enabled"
			}
			fmt.Printf("  %s %s\n", status, pluginName)
		}

		return nil
	},
}

var pluginsEnableCmd = &cobra.Command{
	Use:   "enable <plugin-name>",
	Short: "Enable a specific plugin",
	Long:  `Enable a plugin by name. The plugin must be available in the plugin directory.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginName := args[0]
		settings, err := config.LoadSettings()
		if err != nil {
			return fmt.Errorf("failed to load settings: %w", err)
		}

		loader := plugins.NewLoader()
		err = loader.LoadPlugins(settings.PluginDirectory)
		if err != nil {
			return fmt.Errorf("failed to load plugins: %w", err)
		}

		if !loader.IsPluginLoaded(pluginName) {
			return fmt.Errorf("plugin '%s' not found in directory: %s", pluginName, settings.PluginDirectory)
		}

		if contains(settings.EnabledPlugins, pluginName) {
			fmt.Printf("Plugin '%s' is already enabled\n", pluginName)
			return nil
		}

		settings.EnabledPlugins = append(settings.EnabledPlugins, pluginName)
		err = settings.SaveSettings()
		if err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}

		fmt.Printf("✅ Plugin '%s' enabled successfully\n", pluginName)
		return nil
	},
}

var pluginsDisableCmd = &cobra.Command{
	Use:   "disable <plugin-name>",
	Short: "Disable a specific plugin",
	Long:  `Disable a plugin by name. The plugin will no longer be used in analysis.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginName := args[0]
		settings, err := config.LoadSettings()
		if err != nil {
			return fmt.Errorf("failed to load settings: %w", err)
		}

		if !contains(settings.EnabledPlugins, pluginName) {
			fmt.Printf("Plugin '%s' is not currently enabled\n", pluginName)
			return nil
		}

		// Remove plugin from enabled list
		var newEnabled []string
		for _, enabled := range settings.EnabledPlugins {
			if enabled != pluginName {
				newEnabled = append(newEnabled, enabled)
			}
		}
		settings.EnabledPlugins = newEnabled

		err := settings.SaveSettings()
		if err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}

		fmt.Printf("❌ Plugin '%s' disabled successfully\n", pluginName)
		return nil
	},
}

var pluginsDirCmd = &cobra.Command{
	Use:   "dir [path]",
	Short: "Get or set the plugin directory",
	Long:  `Get the current plugin directory or set a new one.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := config.LoadSettings()
		if err != nil {
			return fmt.Errorf("failed to load settings: %w", err)
		}

		if len(args) == 0 {
			// Get current directory
			fmt.Printf("Current plugin directory: %s\n", settings.PluginDirectory)
			return nil
		}

		// Set new directory
		newDir := args[0]

		// Check if directory exists and is actually a directory
		info, err := os.Stat(newDir)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("directory does not exist: %s", newDir)
			}
			return fmt.Errorf("failed to access path: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path is not a directory: %s", newDir)
		}

		settings.PluginDirectory = newDir
		err = settings.SaveSettings()
		if err != nil {
			return fmt.Errorf("failed to save settings: %w", err)
		}

		fmt.Printf("✅ Plugin directory set to: %s\n", newDir)
		return nil
	},
}

func init() {
	pluginsCmd.AddCommand(pluginsListCmd)
	pluginsCmd.AddCommand(pluginsEnableCmd)
	pluginsCmd.AddCommand(pluginsDisableCmd)
	pluginsCmd.AddCommand(pluginsDirCmd)

	rootCmd.AddCommand(pluginsCmd)
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
