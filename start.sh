#! /bin/bash

/etc/init.d/mysql start
/usr/local/go/bin/go run main.go
