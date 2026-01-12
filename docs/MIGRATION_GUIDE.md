# Migration Guide

This guide provides instructions for upgrading Repo-lyzer from older versions to newer ones. It covers breaking changes, deprecated features, new configurations, and step-by-step migration steps to ensure a smooth transition.

## Upgrading from Version 1.0.x to 1.1.x

### Breaking Changes

- **GitHub API Rate Limiting**: The default rate limit handling has been updated to be more aggressive in preventing API exhaustion. This may result in slower analysis for repositories with high activity if no personal access token is provided.
- **Output Formatting**: The JSON export format now includes additional metadata fields. Existing scripts parsing the output may need updates.

### Deprecated Features

- **Legacy Configuration File**: Support for the old `.repo-analyzer-config` file format is deprecated. Use the new environment variable-based configuration instead.
- **Manual Token Input**: Prompting for GitHub tokens during runtime is deprecated. Set the `GITHUB_TOKEN` environment variable before running the tool.

### New Configurations

- **Environment Variables**: New environment variables for customizing behavior:
  - `REPO_LYZER_TIMEOUT`: Set API request timeout (default: 30s)
  - `REPO_LYZER_CACHE_DIR`: Specify cache directory for API responses
- **Config File**: New YAML-based configuration file support at `~/.repo-lyzer/config.yaml`

### Migration Steps

1. **Update Installation**:
   ```bash
   go install github.com/agnivo988/Repo-lyzer@v1.1.0
   ```

2. **Update Environment Variables**:
   - Replace any manual token prompts with:
     ```bash
     export GITHUB_TOKEN=your_token_here
     ```

3. **Migrate Configuration**:
   - If using old config file, convert to new format:
     ```yaml
     # ~/.repo-lyzer/config.yaml
     github:
       token: your_token_here
       timeout: 30s
     output:
       format: json
     ```

4. **Test Analysis**:
   - Run a test analysis on a small repository to verify output format compatibility.

5. **Update Scripts**:
   - Review any automation scripts that parse Repo-lyzer output and update them to handle new JSON fields.

## Upgrading from Version 0.x to 1.0.x

### Breaking Changes

- **Command Structure**: The CLI command structure has changed from `repo-analyzer` to `repo-lyzer`.
- **Output Directory**: Analysis results are now saved to the current working directory by default instead of a temporary folder.

### Deprecated Features

- **Old Scoring Algorithm**: The legacy health scoring method has been replaced with the new modular analyzer system.

### New Configurations

- **Modular Analyzers**: Configure which analyzers to run via command-line flags or config file.

### Migration Steps

1. **Update Commands**:
   - Change `repo-analyzer analyze` to `repo-lyzer analyze`

2. **Update Scripts**:
   - Modify any scripts or aliases to use the new command name.

3. **Verify Output Location**:
   - Check that analysis files are saved in the expected location.

## General Migration Tips

- **Backup Data**: Always backup important analysis results before upgrading.
- **Test in Staging**: Test the new version on non-production repositories first.
- **Check Dependencies**: Ensure Go version meets the minimum requirements (1.21+).
- **Review Changes**: Check the changelog for any version-specific notes.

## Troubleshooting

### Common Issues

- **Rate Limit Errors**: Set a GitHub personal access token to increase API limits.
- **Configuration Not Loading**: Ensure config files are in the correct location and have proper permissions.
- **Output Parsing Errors**: Update scripts to handle new JSON structure.

### Getting Help

If you encounter issues during migration:
- Check the [GitHub Issues](https://github.com/agnivo988/Repo-lyzer/issues) for similar problems.
- Open a new issue with your migration details and error messages.
- Review the [Contributing Guide](CONTRIBUTING.md) for additional support channels.
