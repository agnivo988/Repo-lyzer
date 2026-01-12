# Task: Fix MainModel menu field type and add security note to migration guide

## Steps to Complete

- [x] Update `internal/ui/app.go`: Change the type of the `menu` field in MainModel from `EnhancedMenuModel` to `MenuModel`
- [x] Update `docs/API_REFERENCE.md`: Change the documentation for the menu field from `EnhancedMenuModel` to `MenuModel`
- [x] Run `go build` to ensure the changes compile successfully
- [x] Update `docs/MIGRATION_GUIDE.md`: Add a security note next to the YAML example warning against storing sensitive credentials, recommending environment variables or secrets managers, and providing an example of loading the token from environment
