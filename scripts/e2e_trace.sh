#!/bin/bash

# Запускаем захват трафика
sudo tcpdump -i lo -w e2e_test.pcap port 8080
