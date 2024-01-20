#!/bin/bash
export pid=3129
mkdir -p /var/run/netns
ln -sf /proc/$pid/ns/net /var/run/netns/hi
ip link add veth1_container type veth peer name veth1_root
ifconfig veth1_container up 
ifconfig veth1_root up
ip link set veth1_container netns hi
ip netns exec hi ifconfig veth1_container up
ip netns exec hi ip addr add 10.1.1.1/24 dev veth1_container
ip netns exec hi ip link set dev veth1_container up