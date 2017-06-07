package main

import (
  "net/http"
  "strings"
  "fmt"
)

func config(node string) string {
raw := `#cloud-config

users:
  - name: core
    passwd: $1$E7Ges1fL$My4ZqJcA70J7ER.ERORqq.
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxGmKrzJdo6u5PsweqpIpuWF7NSwhaj+uaUnL5Nc0n7ZqIssE800chJE96qwSDp3Vepz0omV4V3rUK7Qa5/TZtEEgnhFs0uShOc01qwBfIRqW4VPPm3Z69TNjjVKDKKygZIYgG8hIOWNRG9pM4LcmMGO1NkKvfjG0pW+3yF+WkoicdVCwK7h+haKjOppoey/aNSrIoZdrgveDVuIJpZLHo2PiFvKrMZuEHf2rrX7PpCJzwuPpaESFZtsf2flz+OwM1KmlSz2W5tJdwA/PgEuiYygNeljD2JpAkpuiQxxpkvoI7K3zu9WqOBOqmCmOyGJO2CKqgP4eAwWnO4pIjDm6E0IDKCKWUdqIojbjUTfmE4KL0iqF+C3QIfTe6g/y7efl66K4Rz4MwECb6iHO2/OHswlmTlUJ5VbOnt7uI9+iTmNIf0RhydlVBpzFKGBQ8d59ntgTU4EzwzXZ6bFu3IQW1qgAU1KT5AZ7KfJoLKGVgtql/l2IrGG3OTY8L4t11pFR82bue3lJ8fS/LzElLaIQXw1TFh3Sx/IzXy+jtu78Jft9V4Q2yWGcJ+jlh7CNJ+e0iy1I7AgDdXuqsYc2goGHNsGBhhw4nVDIXFXqRuO47zOVMIyBBAyPpGjNLoZNIfrdY95TcdflqU0WG0Mb7X9GeWK6L/PlBgk/d3QMSE8P7Pw== mpdroog@icloud.com
      - ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEAlM0FRHTjdyqZffgyaEixdhQ6taBJUJXPlA38APbBMCrQZ5z30T+IqJ78C/ydAGfou9/+rTvdr4pxl10GqfkJoj4oSKoxfh6sz9BaNWVBfOaIwS5W3JuAMaTblCrW/pWcnhhkpAWu+fDS6rJSB34shpD5R0N7seXbY1a3ZzptKEOySv19tfo/UizAoj9EVhuV15kWgBHFDctN5L2DOTh5ZngJ/jXWhYxHGcLuY/4xq4ouy4RLYcYNg32UywX82tiXCCw9btunpUQfWrGcawkCo4pMfdWb3sP7PemM5hToeD8ddOIx6Ez5Kty5kcTnWNxQcigg0KAAT4um5HvJiiPGdw== rsa-key-20170517

write_files:
  - path: /etc/ssh/sshd_config
    permissions: 0600
    owner: root:root
    content: |
      # Use most defaults for sshd configuration.
      UsePrivilegeSeparation sandbox
      Subsystem sftp internal-sftp
      Protocol 2

      PermitRootLogin no
      IgnoreRhosts yes
      StrictModes yes
      X11Forwarding no
      AllowUsers core
      PasswordAuthentication no
      ChallengeResponseAuthentication no

  # https://www.jimmycuadra.com/posts/securing-coreos-with-iptables/
  - path: /var/lib/iptables/rules-save
    permissions: 0644
    owner: root:root
    content: |
      *filter
      :INPUT DROP [0:0]
      :FORWARD DROP [0:0]
      :OUTPUT ACCEPT [0:0]
      -A INPUT -i lo -j ACCEPT
      -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
      -A INPUT -p tcp -m tcp --dport 2017 -j ACCEPT
      -A INPUT -p tcp -m tcp --dport 80 -j ACCEPT

      -A INPUT -p tcp -m tcp --dport 8989 -j ACCEPT
      -A INPUT -p tcp -m tcp --dport 26257 -j ACCEPT

      -A INPUT -p icmp -m icmp --icmp-type 0 -j ACCEPT
      -A INPUT -p icmp -m icmp --icmp-type 3 -j ACCEPT
      -A INPUT -p icmp -m icmp --icmp-type 11 -j ACCEPT
      COMMIT

coreos:
  etcd2:
    discovery: https://discovery.etcd.io/c0ee0a7bccfa83c349330891be0e3004
    # multi-region and multi-cloud deployments need to use $public_ipv4
    advertise-client-urls: "http://$public_ipv4:2379"
    initial-advertise-peer-urls: "http://$private_ipv4:2380"
  locksmith:
    window_start: "10:00"
    window_length: "1h"
  update:
    reboot-strategy: "etcd-lock"

  units:
  - name: etcd2.service
    command: start

  - name: sethostname.service
    command: start
    content: |
    [Unit]
      Description=Set Hostname Workaround https://github.com/coreos/bugs/issues/1272

    [Service]
      Type=oneshot
      ExecStart=/bin/sh -c "/usr/bin/hostnamectl set-hostname {hostname}"

      [Install]
        WantedBy=local.target

  - name: sshd.socket
    command: restart
    runtime: true
    content: |
      [Socket]
      ListenStream=2017
      FreeBind=true
      Accept=yes
  - name: iptables-restore.service
    enable: true`

  hostname := node + ".rootdev.nl"
  raw = strings.Replace(raw, "{hostname}", hostname, 1)
  return raw
}

func nodeName(r *http.Request) (string) {
  urls := map[string]bool{
    "ams1": true,
    "nyc1": true,
    "fra1": true,
  }
  name := r.URL.Query().Get("node")
  if name == "" {
    return ""
  }
  if _, ok := urls[name]; !ok {
    return "";
  }
  return name
}

func cloudinit(w http.ResponseWriter, r *http.Request) {
  node := nodeName(r)
  if node == "" {
    w.WriteHeader(http.StatusUnprocessableEntity)
    w.Write([]byte("Unsupported ?node"))
    fmt.Printf("cloudinit: unsupported ?node=%s\n", r.URL.Query().Get("node"))
    return
  }

  raw := config(node)
  if r.URL.Path == "/cloud/init" {
    w.Write([]byte(raw))
  } else if r.URL.Path == "/cloud/ipxe" {
    w.Write([]byte(
      `#!/bin/bash
      cat > "cloud-config.yaml" <<EOF
      ` + raw + `EOF
      sudo coreos-install -d /dev/vda -c cloud-config.yaml
      sudo reboot`,
    ))
  } else {
    w.WriteHeader(http.StatusUnprocessableEntity)
    w.Write([]byte("Unimplemented path"))
    fmt.Printf("cloudinit: invalid path=%s\n", r.URL.Path)
  }
}
