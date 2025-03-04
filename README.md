# aws_ipadd

[![Actions Status](https://github.com/piyushsonigra/aws_ipadd/workflows/Build%20&%20Release/badge.svg)](https://github.com/piyushsonigra/aws_ipadd/actions)

> **Effortlessly manage AWS security group rules with a single command**

## 📖 About

`aws_ipadd` is a CLI tool that simplifies whitelisting and managing IP addresses in AWS security groups. It's designed specifically for scenarios where:

- You don't have a static IP address and your public IP changes frequently
- You need to maintain access to IP-restricted AWS resources
- You want to grant temporary access to specific users by whitelisting their IPs
- You need to maintain tight security by allowing only specific IPs to access particular ports

The tool automatically detects your current public IP and updates AWS security group rules accordingly. Alternatively, you can explicitly specify IPs to whitelist without fetching your current public IP—ideal for adding team members' addresses or other trusted sources.

`aws_ipadd` handles all the AWS security group rule management in the background, making IP whitelisting painless even with constantly changing IPs.

## ✨ Key Features

- **Automatic IP Detection** - Detects and adds your current public IP to security groups
- **Dynamic IP Management** - Updates rules when your public IP changes
- **Multi-Profile Support** - Manage rules across different AWS accounts and regions
- **Port Range Flexibility** - Configure single ports or port ranges
- **CLI Flexibility** - Override configuration with command-line arguments
- **Custom IP Support** - Specify any IP address for whitelisting instead of your current IP
- **Rule Management** - Automatically handles rule creation, updates, and identification

## 🖥️ Supported Operating Systems

- **macOS** (Intel x86_64 and Apple Silicon ARM64)
- **Linux** (x86_64 and ARM64)

## 🚀 Installation

### Linux (x86_64/AMD64)

```console
curl -s -L https://github.com/piyushsonigra/aws_ipadd/releases/latest/download/aws_ipadd_darwin_amd64.tar.gz | tar -xz -C /usr/local/bin
```

### Linux (ARM64)

```console
curl -s -L https://github.com/piyushsonigra/aws_ipadd/releases/latest/download/aws_ipadd_linux_arm64.tar.gz | tar -xz -C /usr/local/bin/
```

### macOS (Intel x86_64)

```console
curl -s -L https://github.com/piyushsonigra/aws_ipadd/releases/latest/download/aws_ipadd_darwin_amd64.tar.gz | tar -xz -C /usr/local/bin/
```

### macOS (Apple Silicon ARM64)

```console
curl -s -L https://github.com/piyushsonigra/aws_ipadd/releases/latest/download/aws_ipadd_darwin_arm64.tar.gz | tar -xz -C /usr/local/bin/
```

> **Note:** If you encounter permission errors, run the command with `sudo` for tar operation as shown example below.

```console
curl -s -L https://github.com/piyushsonigra/aws_ipadd/releases/latest/download/aws_ipadd_darwin_arm64.tar.gz | sudo tar -xz -C /usr/local/bin/
```

## ⚙️ Configuration

1. **Create configuration directory**

   ```console
   mkdir ~/.aws_ipadd
   ```

2. **Create configuration file**

   ```console
   touch ~/.aws_ipadd/aws_ipadd
   ```

3. **Edit the configuration file** with your security group details

### Configuration Parameters

| Parameter | Description |
|-----------|-------------|
| `aws_profile` | AWS CLI profile name |
| `region_name` | AWS region for the security group |
| `security_group_id` | Target security group ID |
| `rule_name` | Descriptive name for the security rule |
| `protocol` | Network protocol (TCP, UDP, or 'all') |
| `port` | Single port to whitelist (ignored if using port range) |
| `from_port` | Start of port range (used with `to_port`) |
| `to_port` | End of port range (used with `from_port`) |

### Sample Configuration

```ini
# Whitelist SSH port
[project-ssh]
aws_profile = aws_project_profile
security_group_id = sg-d26fdre9d
protocol = TCP
port = 22
rule_name = user_name_ssh
region_name = us-east-1

# Whitelist port range
[port-range]
aws_profile = my_project
security_group_id = sg-d26fdre9d
protocol = TCP
from_port = 3000
to_port = 3005
rule_name = office_ind
region_name = us-east-1

# Whitelist all traffic
[project-all-traffic]
aws_profile = project
security_group_id = sg-dfg9dwe
protocol = all
rule_name = all_traffic_from_home
region_name = us-west-2
```

## 🔧 Usage

### Basic Usage

```console
aws_ipadd --profile project-ssh
```

### Update When IP Changes

```console
$ aws_ipadd --profile project-ssh
---------------
project-ssh
---------------
Modifying existing rule...
Removing old whitelisted IP '12.10.1.14/32'.
Whitelisting new IP '131.4.10.16/32'.
Rule successfully updated!
```

### Command-Line Options

```console
Usage:
  aws_ipadd --profile <profile-name>
  aws_ipadd --profile <profile-name> --port <port> --current_ip <current_ip> [options]

Options:
  --profile <string>     aws_ipadd profile name (required)
  --port <int>           Port number (ignored if using port range)
  --from_port <int>      Start of port range (use with to_port)
  --to_port <int>        End of port range (use with from_port)
  --protocol <string>    Protocol e.g tcp, udp, all
  --ip <string>          Custom IP address e.g '10.10.19.1/32'
  --rule_name <string>   Security group rule name
```

### Specify Custom IP

```console
aws_ipadd --profile project-ssh --ip=10.10.10.10/32
```

### Automated Updates with Cron

```console
# Check and update IP every 3 hours
* */3 * * * /usr/local/bin/aws_ipadd --profile project-ssh
```

## 🚀 Upcoming Features

The following features are planned for future releases:

- **Security Group Rule Removal** - Remove specific rules with a simple command
- **Rule Listing** - View all security group rules across profiles in a clean, organized format
- **IPv6 Support** - Full support for IPv6 addresses and dual-stack environments

## 📋 Use Cases

- **Remote Development** - Securely access AWS resources while working from different locations
- **Infrastructure Management** - Simplify access control for DevOps teams with changing IPs
- **Cloud Security** - Maintain tight access controls to sensitive AWS resources
- **Home Office Setup** - Keep consistent access to cloud resources with dynamically assigned ISP IPs
- **Team Access Management** - Easily whitelist team members' IPs for specific resources

## 📜 License

- [MIT License](https://github.com/piyushsonigra/aws_ipadd/blob/master/LICENSE)

## 🙏 Acknowledgements

- [amazonaws_checkip](https://checkip.amazonaws.com) - For IP detection service