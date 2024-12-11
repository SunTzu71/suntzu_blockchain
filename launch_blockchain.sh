#!/bin/bash

# Launch first node
go run main.go chain -port 8000 -miner suntzuchaine1237cd35892b554a49b04a39eb0c648f8fb4875 &

# Wait a bit for the first node to start
sleep 2

# Launch additional nodes
for port in {8001..8003}; do
    go run main.go chain \
        -port $port \
        -miner suntzuchaine1237cd35892b554a49b04a39eb0c648f8fb4875 \
        -remote_node http://127.0.0.1:8000 &
done

# Wait for all background processes
wait
