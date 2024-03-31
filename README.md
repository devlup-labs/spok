# Secure Openpubkey Shell (SOS) (v0.1.0)

## Steps to setup SOS on your server:

## For Arch Users:

1. Download the AUR package for SOS:
(You can install with your favourite AUR helper)

```bash
yay -S sos-bin
```

2. Now, you can configure your server by typing the following commands

```bash
sos configure -s <user>@<server-ip> -e <email-id>
```
(Optional in case of private keys)
```bash
sos configure -i <pvt_key_path> -s <user>@<server-ip> -e <email-id> 
```

3. Now you can login with your email account

```bash
sos login
```

4. Now you can SSH passwordless in your server
```bash
ssh <principal>@<server-ip> 
```