#!/bin/bash

terraform state rm google_compute_network.vpc_network
terraform state rm google_compute_firewall.default
terraform state rm google_compute_instance.nakama_instance

terraform import google_compute_network.vpc_network nakama-instance-network
terraform import google_compute_firewall.default nakama-instance-firewall
terraform import google_compute_instance.nakama_instance nakama-instance

terraform refresh

