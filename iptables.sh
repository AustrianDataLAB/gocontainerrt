#!/bin/bash
export pid=3129
iptables -tnat -N hi
iptables -tnat -A PREROUTING -m addrtype --dst-type LOCAL -j hi 
iptables -tnat -A OUTPUT ! -d 127.0.0.0/8 -m addrtype --dst-type LOCAL -j hi
iptables -tnat -A POSTROUTING -s 10.1.1.1/24 ! -o brg0 -j MASQUERADE
iptables -tnat -A hi -i brg0 -j RETURN

# and in the container run echo "nameserver 8.8.8.8" >> /etc/resolv.conf