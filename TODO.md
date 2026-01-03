# TODO: Fix Tree View Integration and Printing Issues

## Tree View Integration
- [ ] Add stateTree constant to app.go
- [ ] Add tree TreeModel field to MainModel struct
- [ ] Initialize tree model in NewMainModel
- [ ] Add tree state handling in Update method
- [ ] Add tree view in View method
- [ ] Handle 'f' key in dashboard.go to switch to tree state
- [ ] Update BuildFileTree to use actual repo file data from AnalysisResult

## CLI Printing Options
- [ ] Add CLI flags to analyze command (--repo, --langs, --activity, --health, --api, --recruiter, --all)
- [ ] Modify analyze.go to conditionally print based on flags
- [ ] Test all printing functions work correctly

## Comparison Fixes
- [ ] Review compare.go for any issues
- [ ] Fix any bugs in comparison logic
- [ ] Ensure comparison works properly

## Testing
- [ ] Test tree navigation in UI
- [ ] Test selective printing in CLI
- [ ] Test comparison functionality
- [ ] Verify all information prints correctly
