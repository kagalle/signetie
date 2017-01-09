#!/bin/bash
go build github.com/kagalle/signetie/client_golang/gae/login/util && \
go build github.com/kagalle/signetie/client_golang/gae/login/authenticate && \
go build github.com/kagalle/signetie/client_golang/gae/login/codetoken && \
go build github.com/kagalle/signetie/client_golang/gae/login && \
go build main.go
