#!/bin/bash
sudo mkdir -p /mnt/windows_docker
sudo mount -t drvfs '\\wsl.localhost\docker-desktop\mnt\docker-desktop-disk\data\docker' /mnt/windows_docker
