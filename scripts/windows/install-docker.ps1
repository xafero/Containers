Install-Module DockerMsftProvider -Force
Install-Package Docker -ProviderName DockerMsftProvider -Force
(Install-WindowsFeature Containers).RestartNeeded
Restart-Computer

copy daemon.json c:\ProgramData\docker\config\daemon.json
& netsh advfirewall firewall add rule name="Docker Remote Port" dir=in action=allow protocol=TCP localport=2375
docker run hello-world:nanoserver
