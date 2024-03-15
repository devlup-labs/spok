# Secure Openpubkey Shell (SOS)

## Steps to setup SOS on your server:

1. Clone the repository:

```bash
git clone https://github.com/devlup-labs/sos.git
```

2. Add your Oauth credentials in cli/cli.go and verifier/cli.go in place of the dummy credentials.

3. Build your verifier.
```bash
go build -o verifier ./verifier/
```

4. To configure the server with your email id.(Note: you will need to enter the password approximately 4 times.)

```bash
go run cli/cli.go configure <email> <user>@<server-ip>
```

5. Login to your gmail to generate your open pubkey.

```bash
go run cli/cli.go login
```

6. (Optional in some cases) Run to add the ssh key.
```bash
ssh-add
```

Now you should be able to ssh into your server passwordless.
