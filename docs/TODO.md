# TODO: Address High-Priority Issues - Tree View Integration, Comparison Fixes, Analysis Printing Issues

## Tree View Integration
- [x] Update BuildFileTree in internal/ui/tree.go to parse flat FileTree []github.TreeEntry into hierarchical FileNode tree

## Analysis Printing Issues
- [x] Create new print function for file tree in internal/output/tree.go
- [x] Update cmd/analyze.go to fetch file tree and add print call

## Comparison Fixes
- [x] Modify cmd/compare.go to fetch languages for both repos
- [x] Modify cmd/compare.go to fetch file tree for both repos
- [x] Enhance comparison output to include languages and file tree summary

## Testing
- [ ] Test analyze command with file tree printing
- [ ] Test compare command with new data
- [ ] Test tree view in TUI with real data
