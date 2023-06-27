package Util

const MEMORY = "free -m |awk  ' NR == 2 {print $2 \" \" $3/$2*100 \" \" $4/$2*100 \" \"$7/$2*100}'\n"

const CPU = "mpstat"

const DISK = "df . -m"

const SYSTEM = "uname -a |awk '{print $1 \" \" $2 \" \" $4}'"

const SYSTEM_UP_SECONDS = "awk '{print $1}' /proc/uptime"

const Ifconfig = "cat /proc/net/dev | awk 'NR>2 {gsub(/:/,\"\"); print $1, $2, $3, $10, $11, $2+$10, $3+$11}'"

const SSH_TIMEOUT_SECONDS = 10

const PING_TIMEOUT_SECONDS = 10
