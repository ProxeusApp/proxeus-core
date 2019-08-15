# Proxeus UI

### Prerequisites
+ yarn (1.12.3+)
+ node (8.11.3+)
+ vue-cli

## Proxeus Core
For more info check the Proxeus Core [Readme](core/README.md)

## Important
+ Only use yarn 1.12.3. For linking together local dependencies we use Yarn Workspaces:
https://yarnpkg.com/lang/en/docs/workspaces/
+ Use yarn only from the /core/central/ui directory, as the common dependencies will be stored in
/core/central/ui/node_modules instead of the package subfolders (./core, etc.).
