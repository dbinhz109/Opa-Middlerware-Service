#!/bin/bash

# ifacemaker -f impl/account_service.go -s accountService -i IAccountService -p spec > spec/account_service.go
ifacemaker -f impl/"$1"_service.go -s "$1"Service -i I"$1"Service -p spec > spec/"$1"_service.go