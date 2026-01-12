# Settings Backend Implementation

## Overview

The Settings Backend provides full persistence for user preferences including theme selection, export format configuration, and GitHub token management. All settings are stored in `~/.repo-lyzer/settings.json`.

---

## Components

### Config Module (`internal/config/settings.go`)

| Function | Purpose |
|----------|---------|
| `LoadSettings()` | Load settings from disk or return defaults |
| `SaveSettings()` | Persist current settings to disk |
| `ResetToDefaults()` | Reset all settings to defaults |
| `SetTheme()` | Update and persist theme selection |
| `SetExportFormat()` | Update and persist export format |
| `CycleExportFormat()` | Cycle through available export formats |
| `SetGitHubToken()` | Save GitHub token |
| `ClearGitHubToken()` | Remove saved token |
| `GetMaskedToken()` | Return token with characters masked for display |

---

## Settings File Structure

**Location**: `~/.repo-lyzer/settings.json`

```json
{
  "theme_name": "Catppuccin Mocha",
  "default_export_format": "json",
  "export_directory": "/Users/username/Downloads",
  "github_token": "",
  "default_analysis_type": "quick"
}
```

---

## Keyboard Shortcuts

### Theme Settings
| Key | Action |
|-----|--------|
| `1-7` | Select theme by number |
| `t` | Cycle to next theme |

### Export Options
| Key | Action |
|-----|--------|
| `f` | Cycle through formats (JSON → Markdown → CSV → HTML) |

### GitHub Token
| Key | Action |
|-----|--------|
| `i` | Enter token input mode |
| `c` | Clear saved token |
| `Enter` | Save token (in input mode) |
| `ESC` | Cancel input (in input mode) |

### Reset to Defaults
| Key | Action |
|-----|--------|
| `y` | Confirm reset |
| `ESC` | Cancel |

---

## Export Formats

| Format | Extension | Description |
|--------|-----------|-------------|
| JSON | `.json` | Structured data for programmatic use |
| Markdown | `.md` | Human-readable reports |
| CSV | `.csv` | Spreadsheet compatible |
| HTML | `.html` | Web-ready reports |

---

## Integration

Settings are loaded on app startup in `NewMainModel()`:

```go
appConfig, _ := config.LoadSettings()
if appConfig != nil && appConfig.ThemeName != "" {
    SetThemeByName(appConfig.ThemeName)
}
```

Theme changes are saved automatically:
```go
theme := CycleTheme()
if m.appConfig != nil {
    m.appConfig.SetTheme(theme.Name)
}
```

---

## Security

- GitHub tokens are stored in plaintext in the settings file
- Token display is masked (shows first 4 and last 4 characters)
- Environment variable `GITHUB_TOKEN` is still supported as fallback

---

**Status**: ✅ Complete  
**Date**: 2026-01-12  
**Files Added**: `internal/config/settings.go` (207 lines)  
**Files Modified**: `internal/ui/app.go`, `internal/ui/themes.go`
