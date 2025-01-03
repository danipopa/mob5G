#!/bin/bash
cd "$(dirname "$0")"
rm -f /tmp/clixon.sock  # Remove any leftover socket from previous runs
export CLICON_YANG_DIRS=/usr/local/share/clixon
clixon_backend -f clixon.xml -s
