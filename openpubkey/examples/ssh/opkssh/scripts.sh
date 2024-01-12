#!/usr/bin/sh
scp -i sos policy.yml $1:/etc/opk
scp -i sos opk-ssh $1:/root/
scp -i sos configure-opk-server.sh $1:/root/
