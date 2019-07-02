#!/bin/sh
sudo dpkg --add-architecture armhf
sudo apt-get update
sudo curl -sS https://packages.gitlab.com/install/repositories/gitlab/raspberry-pi2/script.deb.sh | sudo os=raspbian dist=stretch bash
sudo apt install gitlab-ce
