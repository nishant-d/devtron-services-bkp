#!/bin/sh
#
# Copyright (c) 2024. Devtron Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

# no arguments passed
# or first arg is `-f` or `--some-option`
if [ "$#" -eq 0 ] || [ "${1#-}" != "$1" ]; then
	# add our default arguments
	set -- dockerd \
		--host=unix:///var/run/docker.sock \
		--host=tcp://0.0.0.0:2375 \
		"$@"
fi

if [ "$1" = 'dockerd' ]; then
	if [ -x '/usr/local/bin/dind' ]; then
		# if we have the (mostly defunct now) Docker-in-Docker wrapper script, use it
		set -- '/usr/local/bin/dind' "$@"
	fi

	# explicitly remove Docker's default PID file to ensure that it can start properly if it was stopped uncleanly (and thus didn't clean up the PID file)
	find /run /var/run -iname 'docker*.pid' -delete
fi

exec "$@"