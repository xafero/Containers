#!/bin/sh
wget http://download.dragondisk.com/dragondisk_1.0.5-0ubuntu_amd64.deb
sudo apt-get install libqt4-dbus libqt4-network libqt4-xml libqtcore4 libqtgui4
sudo dpkg -i dragondisk_1.0.5-0ubuntu_amd64.deb
