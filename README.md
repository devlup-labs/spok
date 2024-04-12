<div align="center">


<img src="https://i.imgur.com/6vd8QF2.png" width=500>

# SPoK - _Sans_ Password or Key

<br>

[![License](https://img.shields.io/badge/License-MIT-blue)](#license)
![Github Release](https://img.shields.io/github/v/release/devlup-labs/spok)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/devlup-labs/spok)

An easier way for remote server authentication. Powered by [OpenPubkey](https://github.com/openpubkey/openpubkey).

</div>

## Features

- **Extensibility**: Eliminate the need of using public keys (RSA, etc.) to add people to your server; you can simply add them using their Email addresses.
- **Scalability**: Add as many people as necessary to your server via their emails.
- **Security**: This project integrates [OpenPubkey](https://github.com/openpubkey/openpubkey), leveraging the OpenID Connect (OIDC) Protocol for enhanced SSH authentication security.
- **Single Command**: Configure your server for SPoK with just one command.
- **Runs Everywhere**: Set up SPoK on any machine—whether it's local, remote, cloud-based, physical server, or a VM—and on any architecture, including x86 or ARM

## Why SPok?

The motivation behind SPoK (Sans Password or Key) for SSH is to revolutionize SSH authentication by addressing security concerns and simplifying key management.

Traditional SSH authentication relies on the manual distribution and management of public keys, which can be error-prone and cumbersome, particularly in environments with numerous users or frequent key rotations. SPoK aims to streamline this process by introducing a modern approach to SSH authentication.

SPoK utilizes a combination of cryptographic techniques, including secure key exchange protocols and cryptographic signatures, to ensure secure and efficient authentication. By integrating SPoK with SSH, users can authenticate using their email addresses , eliminating the need for managing SSH keys separately.

This approach offers several advantages:

1. **Simplified Management**: SPoK eliminates the need for manually distributing and managing SSH keys, reducing administrative overhead and potential errors.

2. **Enhanced Security**: SPoK leverages modern cryptographic techniques to ensure secure authentication, mitigating common security risks associated with SSH key management.

3. **Scalability**: SPoK is designed to scale effectively, allowing organizations to manage authentication for large numbers of users and devices with ease.

4. **Compatibility**: SPoK is compatible with existing SSH infrastructure, making it easy to integrate into existing systems without major modifications.

Overall, SPoK aims to modernize SSH authentication, making it more secure, convenient, and scalable for organizations of all sizes. By eliminating the complexities associated with traditional SSH key management, SPoK offers a streamlined solution that meets the security needs of today's dynamic computing environments.

## Installation

### Linux:

#### Arch-based Distros (Arch Linux, EndeavourOS, Manjaro, etc.):

- Download the AUR package for SPoK:
  (You can install it with your favourite AUR helper)

```bash
yay -S spok-bin
```

#### Debian-based Distros (Debian, Ubuntu, Linux Mint, etc.):

- You can install by running the following commands

```bash
curl -s https://packagecloud.io/install/repositories/SaahilNotSahil/spok/script.deb.sh?any=true | sudo bash
sudo apt update
sudo apt install spok
```

#### RHEL-based Distros (RHEL, Fedora, CentOS, etc.):

- You can install by running the following commands

```shell
curl -s https://packagecloud.io/install/repositories/SaahilNotSahil/spok/script.rpm.sh?any=true | sudo bash
sudo rpm install spok
```

#### From archive:

- Download the latest release (`spok_<version>_linux_<amd64/arm64>.tar.gz`) from [here](https://github.com/devlup-labs/spok/releases).
- Extract the `tar.gz` file and run the installer script
```shell
tar zxvf spok_<version>_linux_<amd64/arm64>.tar.gz
chmod +x install.sh
./install.sh
```
- SPoK is now installed on your system in the `/usr/bin` directory, which is already in the PATH.

### Mac:

#### Homebrew

- You can install by running the following commands

```shell
brew tap devlup-labs/spok
brew install spok
```

- To upgrade the package:

```shell
brew upgrade spok
```

### Windows:

#### Scoop:

- First, you need to install [Scoop](https://scoop.sh/).
- Next, run the following commands in PowerShell

```shell
scoop bucket add org https://github.com/devlup-labs/scoop-spok.git
scoop install spok
```

#### From archive:

- Download the latest release (`spok_<version>_windows_amd64.zip`) from [here](https://github.com/devlup-labs/spok/releases).
- Extract the zip file.
- Open Powershell as administrator and run the following commands

```shell
cd <path-to-extracted-folder>
.\install.ps1
```

- SPoK is now installed on your system in the `C:\Program Files\SPoK` directory, and is added to the PATH.

## Usage

SPoK consists of two parts: the `spok` client CLI tool, and the `verifier` server-side tool, which is downloaded automatically while configuring your server to use SPoK.  
You must have access to the `root` user on the server, or any other user with `sudo` privileges, to configure the server to use SPoK.

### Client side:

- Configure your server by running the following command

```shell
spok configure -s <user>@<server-ip-or-hostname> -e <email-id>
```

(Optional in case of key-pair authentication)

```shell
spok configure -i <pvt_key_path> -s <user>@<server-ip> -e <email-id>
```

- Login to your email account (the one you provided while configuring the server)

```shell
spok login
```

#### Note: Currently works only with Google (Gmail + Google Workspace) accounts.

- Now you can SSH into your server, and it won't ask for a password or key

```shell
ssh <user>@<server-ip>
```

#### Note 2: The server must have an active internet connection for configuring SPoK, as well as every time you SSH into the server. If it ever loses internet connectivity, you can always fall back to using a password or key.

#### Note 3: Currently the validity of the token is 1 hour. After that, you will have to re-login to your email account. Just use the `spok login` command again.

### Server side:

Once the server is successfully configured to use with SPoK, you'll find a new directory `/etc/spok` on the server, which contains two files - a `policy.yml` file and the the `verifier` program. If any of these files are missing or empty, make sure the server has an active internet connection, the user is root or has sudo privileges, and run the `configure` command again.

#### Policy.yml:

The `policy.yml` file contains the information regarding which email addresses can access which users, or `principals`, on the server, using SPoK.  
This file can be edited directly using a text editor, or preferably using the `verifier` tool.

#### Verifier:

The `verifier` tool primarily serves two purposes:-

1. Once configured, the `verifier` becomes the default authentication provider for `sshd` on the server. When you ssh into the server using a certificate that is generated by `spok login`, the verifier verifies the certificate for authenticity, and also checks it against the policy stored in the `policy.yml` file. Once verified, you are automatically logged into the server. If it fails to verify, it'll fall back on other configured modes of authentication.
2. It can also be used to add/remove `principals` for different email addresses in the `policy.yml` file.

- To add a new principal called `user` for the email address `someone@example.com`, run the following command:

```shell
/etc/spok/verifier add someone@example.com user
```

- Similarly, to remove the principal, use the `remove` command:

```shell
/etc/spok/verifier remove someone@example.com user
```

## License

This repository contains SPoK, covered under the [MIT License](LICENSE), except where noted.

It is distributed under the terms of the MIT License.

Third parties are permitted to distribute the software independently, but they are restricted from utilizing any SPoK trademarks, proprietary cloud services, etc.

We expressly authorize you to incorporate our trademarks while developing SPoK itself. However, you are prohibited from publishing or sharing the resulting build, and you may not employ that build to operate SPoK for any other purpose.
