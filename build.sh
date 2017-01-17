#!/bin/bash
go build github.com/kagalle/signetie/client_golang/gae/util && \
go build github.com/kagalle/signetie/client_golang/gae/authenticate && \
go build github.com/kagalle/signetie/client_golang/gae/login && \
go build main.go
