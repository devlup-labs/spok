# SPoK - *Sans* Password or Key (v0.1.0)
An easier way for remote server authentication. Powered bv [OpenPubkey](https://github.com/openpubkey/openpubkey).

## Installation

### For Arch Users:

- Download the AUR package for SPoK:
   (You can install it with your favourite AUR helper)

```bash
yay -S spok-bin
```

### For Debian-based (Debian, Ubuntu, Linux Mint, etc.) Users:

- You can install by running the following commands
```bash
curl -s https://packagecloud.io/install/repositories/SaahilNotSahil/spok/script.deb.sh?any=true | sudo bash
sudo apt update
sudo apt install spok
```

### Coming soon for other Operation Systems :)

## Setup

1. Configure your server by typing the following commands

```bash
spok configure -s <user>@<server-ip> -e <email-id>
```

(Optional in case of key-pair authentication)

```bash
spok configure -i <pvt_key_path> -s <user>@<server-ip> -e <email-id>
```

2. Now you can login with your email account

```bash
spok login
```

3. Now you can SSH passwordless in your server

```bash
ssh <principal>@<server-ip>
```

### Note: Currently works only with Gmail accounts.
