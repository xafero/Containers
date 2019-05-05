#!/bin/sh
NSPR_LOG_MODULES=timestamp,nsHttp:3,sync NSPR_LOG_FILE=/tmp/ff.log firefox -private
