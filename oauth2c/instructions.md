## Instructions

First Install oauth2c(For linux):

```shell
curl -sSfL https://raw.githubusercontent.com/cloudentity/oauth2c/master/install.sh | \
  sudo sh -s -- -b /usr/local/bin latest
```

or 

```shell
go install github.com/cloudentity/oauth2c@latest
```

Now go to /oauth2c and run the following

```shell
chmod +x testing.sh
./testing.sh
```