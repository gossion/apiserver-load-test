#!/bin/bash

#nodeips=("10.224.0.5" "10.224.0.6" "10.224.0.7" "10.224.0.8" "10.224.0.9" "10.224.0.10" "10.224.0.11" "10.224.0.12" "10.224.0.13" "10.224.0.14") # "10.224.0.15" "10.224.0.716"
nodeips=(
"10.224.0.5" "10.224.0.6" "10.224.0.7" "10.224.0.8" "10.224.0.9"
"10.224.0.10" "10.224.0.11" "10.224.0.12" "10.224.0.13" "10.224.0.14" "10.224.0.15" "10.224.0.16" "10.224.0.17" "10.224.0.18" "10.224.0.19"
"10.224.0.20" "10.224.0.21" "10.224.0.22" "10.224.0.23" "10.224.0.24" "10.224.0.25" "10.224.0.26" "10.224.0.27" "10.224.0.28" "10.224.0.29"
"10.224.0.30" "10.224.0.31" "10.224.0.32" "10.224.0.33" "10.224.0.34")
privatekey_path="ssh.key"
user="azureuser"
apiserverip=10.224.0.4

run_ssh() {
  local privatekey_path=$1
  local user=$2
  local ip=$3
  local port=$4
  local command=$5

  sshCommand="ssh -i $privatekey_path -A -p $port $user@$ip -2 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -o PreferredAuthentications=publickey -o PasswordAuthentication=no -o ConnectTimeout=5 -o GSSAPIAuthentication=no -o ServerAliveInterval=30 -o ServerAliveCountMax=10 $command"
  $sshCommand
}

get_established_connections() {
  local privatekey_path=$1
  local user=$2
  local ip=$3

  port=22
  run_ssh $privatekey_path $user $ip $port "sudo ss"
}


echo "IP,Established,Connecting"
for ip in "${nodeips[@]}"
do
  get_established_connections $privatekey_path $user $ip > ${ip}.log 2>/dev/null
  established=0
  connecting=0
  established=$(cat ${ip}.log | grep $apiserverip | grep "ESTAB" | wc -l)
  connecting=$(cat ${ip}.log | grep $apiserverip | grep -v "ESTAB" | wc -l)
  echo $ip,$established,$connecting
done

