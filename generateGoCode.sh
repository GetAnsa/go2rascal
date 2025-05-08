#!/bin/bash

cat tokenConversions | awk -f convertTokens.awk > tokens.go
