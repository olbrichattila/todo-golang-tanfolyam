#!/bin/bash

REMOTE_SERVER=root@167.71.134.25

scp todo.tar $REMOTE_SERVER:/root/todo
scp deploy.sh $REMOTE_SERVER:/root/todo
ssh $REMOTE_SERVER 'cd /root/todo && chmod +x ./deploy.sh && ./deploy.sh'
