#!/bin/bash

mkdir logs
mkdir sessions

tarantool init_tarantool.lua 2>./logs/tarantool-stderr.log  >./logs/tarantool-stdout.log &
sleep 3
./crypto 2>./logs/server-stderr.log >./logs/server-stdout.log &
