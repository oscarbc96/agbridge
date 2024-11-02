# `agbridge` CLI

![Latest GitHub release](https://img.shields.io/github/release/oscarbc96/agbridge.svg)
[![Tests](https://github.com/oscarbc96/agbridge/workflows/test/badge.svg)](https://github.com/oscarbc96/agbridge/actions?query=workflow%3A"test")
[![golangci-lint](https://github.com/oscarbc96/agbridge/workflows/golangci-lint/badge.svg)](https://github.com/oscarbc96/agbridge/actions?query=workflow%3A"golangci-lint")

AGBridge is a command-line tool that acts as a proxy, forwarding requests to private API gateways running in AWS. It’s designed to handle API gateways that are not publicly accessible, allowing secure and efficient access to private resources. Ideal for debugging, testing, or integration scenarios, AGBridge prevents the need to expose sensitive resources while enabling smooth access.

AGBridge supports multiple API gateways simultaneously, making it suitable for integration testing across different services.

## Usage

```bash
agbridge [flags]
```

### Flags

| Flag               | Description                                                                                                      | Default       |
|--------------------|------------------------------------------------------------------------------------------------------------------|---------------|
| `--version`        | Displays the application version and exits.                                                                      |               |
| `--config`         | Path to a configuration file for AGBridge. This flag cannot be used with `--profile-name` or `--resource-id`.    |               |
| `--profile-name`   | Specifies the AWS profile name to access resources. Requires `--resource-id` to be specified.                    |               |
| `--resource-id`    | Specifies the resource ID of the AWS API gateway. Required if `--config` is not provided.                        |               |
| `--log-level`      | Sets the logging level for output messages. Options: `debug`, `info`, `warn`, `error`, `fatal`.                  | `info`        |
| `--listen-address` | Address where AGBridge will listen for incoming requests. Format should be `host:port`.                          | `:8080`       |

### Examples

#### Specify Resource with Profile
Specify a resource and profile to access a private API gateway:
```bash
agbridge --profile-name=myprofile --resource-id=12345
```

#### Load a Specific Configuration File
Run AGBridge with a configuration file:
```bash
agbridge --config=config.yaml
```

#### Change Listen Address
Set a custom port for AGBridge to listen on:
```bash
agbridge --listen-address=:9090
```

## Installationq

### Option 1: Using Homebrew

1. Add the Homebrew tap:
   ```bash
    brew tap oscarbc96/agbridge git@github.com:oscarbc96/agbridge.git
   ```
2. Install:
   ```bash
   brew install agbridge
   ```

### Option 2: Download from Releases

1. Visit the [Releases page](https://github.com/oscarbc96/agbridge/releases) on GitHub.
2. Download the appropriate binary for your operating system.
3. Make the binary executable (if on Linux or macOS):
   ```bash
   chmod +x agbridge
   ```

### Option 3: Using Docker

1. Pull the latest Docker image:
   ```bash
   docker pull ghcr.io/oscarbc96/agbridge:latest
   ```
2. Run the container with appropriate flags. For example:
   ```bash
   docker run --rm -it -p 8080:8080 ghcr.io/oscarbc96/agbridge:latest --profile-name=myprofile --resource-id=12345 --listen-address=:8080
   ```

### Option 4: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/oscarbc96/agbridge
   cd agbridge
   ```
2. Build the CLI:
   ```bash
   make snapshot
   ```
3. Run the CLI:
   ```bash
   dist/agbridge
   ```
