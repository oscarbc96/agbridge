# agbridge

[![Latest GitHub release](https://img.shields.io/github/release/oscarbc96/agbridge.svg)](https://github.com/oscarbc96/agbridge/releases)
[![Tests](https://github.com/oscarbc96/agbridge/workflows/test/badge.svg)](https://github.com/oscarbc96/agbridge/actions?query=workflow%3A"test")
[![golangci-lint](https://github.com/oscarbc96/agbridge/workflows/golangci-lint/badge.svg)](https://github.com/oscarbc96/agbridge/actions?query=workflow%3A"golangci-lint")
[![Go Report Card](https://goreportcard.com/badge/github.com/oscarbc96/agbridge)](https://goreportcard.com/report/github.com/oscarbc96/agbridge)
[![GitHub License](https://img.shields.io/github/license/oscarbc96/agbridge)](https://github.com/oscarbc96/agbridge/blob/main/LICENSE)

`agbridge` is a lightweight CLI tool that acts as a local proxy to private AWS API Gateways, enabling you to securely forward HTTP requests to API Gateway endpoints that are **not publicly accessible**.

It’s designed to help developers and integration systems interact with private services in AWS—such as internal microservices—without the need to expose them to the internet or configure complex VPNs or tunnels.
Designed for testing, debugging, and automation, agbridge supports working with multiple API Gateways and is ideal for service integration workflows in isolated environments.

## 🚀 Features
- 🔒 Secure access to private API Gateways without exposing them publicly.
- 🧪 Simplifies testing and integration with internal AWS services from local environments or CI/CD pipelines.
- ⚙️ Flexible configuration, either through CLI flags or a YAML config file.
- 🌐 Supports multiple API Gateway definitions in a single run.
- 🐳 Docker-ready, perfect for ephemeral or automated environments.

Whether you’re building microservices, automating tests, or debugging internal APIs, agbridge gives you a safe and developer-friendly way to reach your private AWS resources.

## ⚙️ Usage

```bash
agbridge [flags]
```

### Flags

| Flag               | Description                                                                                                                | Default |
|--------------------|----------------------------------------------------------------------------------------------------------------------------|---------|
| `--version`        | Displays the application version and exits.                                                                                |         |
| `--config`         | Path to a configuration file for AGBridge. This flag cannot be used with `--profile-name`, `--rest-api-id`, or `--region`. |         |
| `--profile-name`   | Specifies the AWS profile name to access resources. Requires `--rest-api-id` and `--region` to be specified.               |         |
| `--rest-api-id`    | Specifies the Rest API ID of the AWS API gateway. Required if `--config` is not provided.                                  |         |
| `--region`         | Specifies the AWS region for the API gateway. Requires `--rest-api-id` and `--profile-name`.                               |         |
| `--log-level`      | Sets the logging level for output messages. Options: `debug`, `info`, `warn`, `error`, `fatal`.                            | `info`  |
| `--listen-address` | Address where AGBridge will listen for incoming requests. Format should be `host:port`.                                    | `:8080` |

### 🧪 Examples

#### Specify API GW with Profile
Specify a resource and profile to access a private API gateway:
```bash
agbridge --profile-name=myprofile --rest-api-id=12345
```

#### Load a Specific Configuration File
Run AGBridge with a configuration file:
```bash
agbridge --config=config.yaml
```
config.yaml
```yaml
gateways:
  - rest_api_id: xyz789ghi0
    profile_name: abc123def4
    region: eu-west-1

  - rest_api_id: 789ghi0xyz
    profile_name: myawsprofile
    region: eu-east-1
```

#### Change Listen Address
Set a custom port for AGBridge to listen on:
```bash
agbridge --listen-address=:9090
```

## 📦 Installation

### 🔧 Option 1: Using Homebrew (macOS & Linux)

1. Add the Homebrew tap:
   ```bash
    brew tap oscarbc96/agbridge git@github.com:oscarbc96/agbridge.git
   ```
2. Install:
   ```bash
   brew install agbridge
   ```

### 🧊Option 2: Download from Releases

1. Visit the [Releases page](https://github.com/oscarbc96/agbridge/releases) on GitHub.
2. Download the appropriate binary for your operating system.
3. Make the binary executable (if on Linux or macOS):
   ```bash
   chmod +x agbridge
   ```

### 🐳 Option 3: Using Docker

1. Pull the latest Docker image:
   ```bash
   docker pull ghcr.io/oscarbc96/agbridge:v0.0.10
   ```
2. Run the container with appropriate flags. For example:
   ```bash
   docker run --rm -it -p 8080:8080 ghcr.io/oscarbc96/agbridge:latest --profile-name=myprofile --rest-api-id=12345 --listen-address=:8080
   ```

### 🛠 Option 4: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/oscarbc96/agbridge
   cd agbridge
   ```
2. Build the CLI:
   ```bash
   make
   ```
3. Run the CLI:
   ```bash
   dist/agbridge
   ```
