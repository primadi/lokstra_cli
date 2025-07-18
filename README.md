# Lokstra CLI

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**Lokstra CLI** is the official command-line tool for building, scaffolding, and managing Lokstra-based backend applications. It provides an opinionated toolkit to create scalable, structured Go backend services with built-in linting and project management capabilities.

## Quick Start

```bash
# Install Lokstra CLI
go install github.com/primadi/lokstra_cli/lokstra@latest

# Create a new server project
lokstra init server my-app

# Navigate and run
cd my-app
go run cmd/main.go

# Lint your project
lokstra lint
```

## Features

- 🚀 **Project Scaffolding**: Quickly generate new Lokstra projects of various types
- 🔍 **Built-in Linting**: Validate your code and configuration files
- 📁 **Multiple Project Types**: Support for servers, modules, services, middleware, and plugins (.so binary)
- 🎯 **Service URI Validation**: Ensure proper lokstra:// URI format compliance
- ⚙️ **Template-based Generation**: Customizable project templates
- 🛠️ **Go Module Management**: Automatic dependency management and tidying

## Installation

### Using Go Install

```bash
go install github.com/primadi/lokstra_cli/lokstra@latest
```

### From Source

```bash
git clone https://github.com/primadi/lokstra_cli.git
cd lokstra_cli
go build -o lokstra ./lokstra/
```

### Binary Release

Download the latest binary from the [releases page](https://github.com/primadi/lokstra_cli/releases).

## Commands

Lokstra CLI provides two main commands for project management:

### `lokstra init` - Initialize New Projects

Create a new Lokstra project with comprehensive scaffolding and configuration.

#### Syntax
```bash
lokstra init [project-type] [name] [flags]
```

#### Available Project Types
- **`server`** - Complete Lokstra server application with routing, handlers, middleware, and services
- **`module`** - Collection of handlers, middlewares, and services that can be reused by other projects
- **`service`** - Single service component that can be reused by other projects
- **`middleware`** - Single middleware component that can be reused by other projects
- **`plugin`** - Collection of handlers, middlewares, and services designed to be compiled as binary (.so) plugins

#### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--module` | | Go module name | `github.com/example/<name>` |
| `--template` | | Template directory path or name | `default` |
| `--output` | `-o` | Output directory for the project | current directory |

#### Template Resolution

The `--template` flag supports flexible template resolution:

1. **Environment Variable**: If no template specified, uses `LOKSTRA_TEMPLATE` env var
2. **Direct Path**: If template is a valid directory path, uses it directly
3. **Built-in Templates**: Looks for template under `./scaffold/<template>/[project-type]/` directory
4. **Fallback**: Uses `default` template (`./scaffold/default/[project-type]/`) if none found

**Template Path Examples for server project:**
- `--template default` → `./scaffold/default/server/`
- `--template production` → `./scaffold/production/server/`
- `--template /custom/path` → `/custom/path/server/`
- No template specified → `./scaffold/default/server/`

#### Output Directory

The `--output` flag allows you to specify where the project should be created:

- **Default**: Creates project in `./project-name/`
- **Custom**: Creates project in `output-dir/project-name/`
- **Auto-create**: Automatically creates the output directory if it doesn't exist

#### Examples

```bash
# Basic server creation
lokstra init server my-api

# Custom module name
lokstra init server my-api --module github.com/mycompany/my-api

# Custom output directory
lokstra init server my-api -o /path/to/projects

# Custom template
lokstra init server my-api --template custom-server

# Custom template from absolute path
lokstra init server my-api --template /path/to/my-templates

# Multiple flags combined
lokstra init server my-api \
  --module github.com/mycompany/my-api \
  --template production \
  --output ../projects

# Other project types
lokstra init service user-service --module github.com/mycompany/services/user
lokstra init middleware auth-middleware -o ./middleware  
lokstra init module common-handlers --module github.com/mycompany/modules/common
lokstra init plugin logger-plugin --template custom-plugin
```

#### Project Type Details

##### Server (`lokstra init server my-app`)
Creates a complete Lokstra server application with:
- HTTP routing configuration
- Request handlers implementation
- Middleware stack setup
- Service layer architecture
- Ready-to-run main application

##### Module (`lokstra init module my-module`)
Creates a reusable collection containing:
- Multiple handlers for different endpoints
- Multiple middleware components
- Multiple service implementations
- Designed to be imported by other projects
- Provides shared functionality across applications

##### Service (`lokstra init service my-service`)
Creates a single service component with:
- Focused service interface and implementation
- Business logic for specific domain
- Can be imported and used by other projects
- Lightweight and single-purpose

##### Middleware (`lokstra init middleware my-middleware`)
Creates a single middleware component with:
- HTTP request/response processing logic
- Reusable across different applications
- Follows Lokstra middleware patterns
- Can be imported by servers and modules

##### Plugin (`lokstra init plugin my-plugin`)
Creates a plugin collection similar to module but:
- Designed to be compiled as shared object (.so) binary


#### What Gets Generated

When you run `lokstra init server my-app`, it:

1. **Creates project directory** structure
2. **Generates `go.mod`** with specified module name
3. **Copies template files** and processes `.tpl` files
4. **Downloads dependencies** (`go get github.com/primadi/lokstra@latest`)
5. **Tidies modules** (`go mod tidy`)

### `lokstra lint` - Code Quality Validation

Comprehensive linting for Lokstra projects, checking Go code, YAML configuration, and Lokstra-specific patterns.

#### Syntax
```bash
lokstra lint
```

#### What Gets Checked

##### Go Files (`*.go`)
- **Service URI Format**: Validates `lokstra://` URI patterns
- **Interface Naming**: Ensures CamelCase naming conventions
- **Code Structure**: Checks for common Go best practices

##### YAML Files (`*.yaml`, `*.yml`)
- **Syntax Validation**: Ensures valid YAML structure
- **Configuration Schema**: Validates Lokstra configuration format
- **Service Definitions**: Checks service configuration completeness

##### Service URI Validation
The linter specifically validates Lokstra service URIs:

```
lokstra://[package.]ServiceType/instance-name
```

**Valid Examples:**
- `lokstra://UserService/main`
- `lokstra://auth.TokenValidator/jwt`
- `lokstra://db.Repository/users`

**Invalid Examples:**
- `lokstra://userService/main` (should be CamelCase)
- `lokstra://User_Service/main` (no underscores)
- `lokstra://UserService` (missing instance name)

#### Excluded Directories

The linter automatically skips these directories:
- `.git` - Git version control
- `vendor` - Go vendor dependencies
- `node_modules` - Node.js dependencies
- `bin` - Binary files
- `dist` - Distribution/build files

#### Output Examples

**Success:**
```bash
🔍 Found 15 Go files, 3 YAML files
✅ No lint issues found
```

**With Issues:**
```bash
🔍 Found 15 Go files, 3 YAML files
❌ internal/service.go: invalid service URI format: lokstra://userService/main
❌ config/lokstra.yaml: invalid YAML syntax at line 12
❌ Total lint issues: 2
```

#### Exit Codes
- **0**: No issues found
- **1**: Issues found or linting error

#### Integration with CI/CD

The lint command is perfect for continuous integration:

```yaml
# GitHub Actions example
- name: Lint Lokstra Project
  run: lokstra lint
```

```bash
# Pre-commit hook
#!/bin/sh
lokstra lint || exit 1
```

### Template System

Lokstra CLI uses a template-based generation system located in `./scaffold/` directory:

```
scaffold/
├── default/              # Default template collection
│   ├── server/           # Server project templates
│   ├── module/           # Module project templates  
│   ├── service/          # Service project templates
│   ├── middleware/       # Middleware project templates
│   └── plugin/           # Plugin project templates
├── production/           # Production template collection (example)
│   ├── server/           # Production server templates
│   ├── module/           # Production module templates
│   └── ...
└── custom/               # Custom template collection (example)
    ├── server/           # Custom server templates
    └── ...
```

**Template Structure:**
- **Template Collections**: Each directory under `scaffold/` is a template collection
- **Project Types**: Each template collection contains subdirectories for project types
- **Template Resolution**: `--template production` uses `scaffold/production/[project-type]/`
- **Default Template**: If no template specified, uses `scaffold/default/[project-type]/`

**Template Files:**
- **`.tpl` files**: Processed with Go templates, variables substituted
- **Regular files**: Copied as-is to the target project
- **Template variables**: `{{.AppName}}`, `{{.ModuleName}}` available in `.tpl` files

## Configuration

### Lokstra Configuration (`lokstra.yaml`)

The generated `lokstra.yaml` file defines your application structure:

```yaml
apps:
  - name: main
    address: :8080
    middleware:
      - name: lokstra.recovery

services: []
```

### Service URI Format

Lokstra uses a special URI format for service references:

```
lokstra://[package.]ServiceType/instance-name
```

**Examples:**
- `lokstra://UserService/main`
- `lokstra://auth.TokenValidator/jwt`
- `lokstra://db.Repository/users`

## Development

### Prerequisites

- Go 1.24 or later
- Git

### Building from Source

```bash
git clone https://github.com/primadi/lokstra_cli.git
cd lokstra_cli
go mod tidy
go build -o lokstra ./lokstra/
```

### Testing

```bash
# Run all tests
go test ./...

# Test the CLI locally
./lokstra --help
./lokstra init --help
./lokstra lint --help

# Test project creation
./lokstra init server test-app
cd test-app
go run cmd/main.go
```

### Project Structure

```
lokstra_cli/
├── cmd/                  # CLI commands implementation
│   ├── root.go          # Root command definition
│   └── init.go          # Init and lint commands
├── internal/            # Internal packages
│   ├── lint/            # Linting functionality
│   └── uri/             # URI validation
├── scaffold/            # Project templates
│   └── server/          # Server project template
├── lokstra/             # Main entry point
│   └── main.go          # CLI entry point
├── go.mod
├── go.sum
└── README.md
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Lokstra Framework](https://github.com/primadi/lokstra) - The core Lokstra framework for Go backend development

## Support

- 📖 [Documentation](https://lokstra.dev/docs)
- 💬 [GitHub Discussions](https://github.com/primadi/lokstra_cli/discussions)
- 🐛 [Issue Tracker](https://github.com/primadi/lokstra_cli/issues)

---

Built with ❤️ by the Lokstra team