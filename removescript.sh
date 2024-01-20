#!/usr/bin/env bash
export pid=3129
ip netns exec hi ifconfig veth1_container down
ip link delete veth1_container netns hi
ifconfig veth1_root down
ifconfig veth1_container down
ip link delete veth1_container 
rm /var/run/netns/hi