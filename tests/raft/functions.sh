#!/bin/bash


function get {
    ./rafter data get --key="${1}" --peer=@tests/raft/node1.json --peer=@tests/raft/node2.json --peer=@tests/raft/node3.json --peer=@tests/raft/node4.json --peer=@tests/raft/node5.json
}

function put {
    ./rafter data set --key="${1}" --value="${2}" --peer=@tests/raft/node1.json --peer=@tests/raft/node2.json --peer=@tests/raft/node3.json --peer=@tests/raft/node4.json --peer=@tests/raft/node5.json
}

function list {
    ./rafter data list --peer=@tests/raft/node1.json --peer=@tests/raft/node2.json --peer=@tests/raft/node3.json --peer=@tests/raft/node4.json --peer=@tests/raft/node5.json
}

function loop {
    for i in $(seq 1 ${1}); do 
        put "key_$i" "value_$i" 
        get "key_$i" 
    done
}

function fill {
    for i in $(seq 1 ${1}); do 
        put "key_$i" "value_$i" 
        get "key_$i" 
    done
}


