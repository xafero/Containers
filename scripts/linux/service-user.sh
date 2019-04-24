#!/bin/sh
sudo addgroup netwww
sudo useradd -M kestrel
sudo usermod -L kestrel
sudo adduser kestrel netwww

