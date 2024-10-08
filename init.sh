#!/bin/sh

set -o errexit
set -o nounset

# install CNIs
cp -f /usr/bin/terway /opt/cni/bin/
chmod +x /opt/cni/bin/terway

if [ "$TERWAY_DAEMON_MODE" != "VPC" ]; then
  cp -f /usr/bin/cilium-cni /opt/cni/bin/
  chmod +x /opt/cni/bin/cilium-cni
fi

# init cni config
cp /tmp/eni/eni_conf /etc/eni/eni.json

terway-cli cni /tmp/eni/10-terway.conflist /tmp/eni/10-terway.conf --output /etc/cni/net.d/10-terway.conflist
terway-cli nodeconfig

node_capabilities=/var/run/eni/node_capabilities
if [ ! -f "$node_capabilities" ]; then
  echo "Init node capabilities"
  mkdir -p /var/run/eni
  touch "$node_capabilities"
fi

require_erdma=$(jq '.enable_erdma' -r </etc/eni/eni.json)
if [ "$require_erdma" = "true" ]; then
  echo "Init erdma driver"
  if modprobe erdma; then
    echo "node support erdma"
    echo "erdma = true" >>"$node_capabilities"
    if ! grep -q "erdma *= *true" "$node_capabilities"; then
      sed -i '/erdma *=/d' "$node_capabilities"
      echo "erdma = true" >> "$node_capabilities"
    fi
  else
    sed -i '/erdma *= *true/d' "$node_capabilities"
    echo "node not support erdma, pls install the latest erdma driver"
  fi
fi

# copy node capabilities to tmpfs so policy container can read it
cp $node_capabilities /var-run-eni/node_capabilities

sysctl -w net.ipv4.conf.eth0.rp_filter=0
modprobe sch_htb || true

if [ "$TERWAY_DAEMON_MODE" != "VPC" ]; then
  chroot /host sh -c "systemctl disable eni.service; rm -f /etc/udev/rules.d/75-persistent-net-generator.rules /lib/udev/rules.d/60-net.rules /lib/udev/rules.d/61-eni.rules /lib/udev/write_net_rules"
fi
