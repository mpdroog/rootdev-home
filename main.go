// main contains a super lightweight HTTP-server
// so we can let this 'daemon' work with a really
// small attach vector (keeping maintenance super low)
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/schema"
	"gopkg.in/mailgun/mailgun-go.v1"
	"gopkg.in/validator.v2"
	"net/http"
	"os"
	"strings"
)

var (
	verbose       bool
	mailgunDomain string
	mailgunApi    string
	mailgunApiPub string

	decode = schema.NewDecoder()
)

func emailDecode(r *http.Request) (Email, error) {
	msg := Email{}

	if r.Body == nil {
		return msg, fmt.Errorf("Empty body")
	}
	if e := r.ParseForm(); e != nil {
		return msg, e
	}

	if e := decode.Decode(&msg, r.PostForm); e != nil {
		return msg, e
	}
	if e := validator.Validate(msg); e != nil {
		return msg, e
	}
	return msg, nil
}

func email(w http.ResponseWriter, r *http.Request) {
	msg, e := emailDecode(r)
	if e != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Sorry invalid input."))
		fmt.Printf("ERR:email:validate: %s\n", e)
		return
	}

	mg := mailgun.NewMailgun(mailgunDomain, mailgunApi, mailgunApiPub)
	message := mailgun.NewMessage("noreply@rootdev.nl", "Contact request", "From "+msg.Email+"\n\n"+msg.Body, "rootdev@gmail.com")

	if _, idx, e := mg.Send(message); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry failed sending email."))
		fmt.Printf("ERR:email:send(%s): %s\n", idx, e)
		return
	}
	w.Write([]byte("Sent email."))
}

func ipxe(w http.ResponseWriter, r *http.Request) {
	raw := `#!/bin/bash
cat > "cloud-config.yaml" <<EOF
#cloud-config

users:
  - name: core
    passwd: $1$E7Ges1fL$My4ZqJcA70J7ER.ERORqq.
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCxGmKrzJdo6u5PsweqpIpuWF7NSwhaj+uaUnL5Nc0n7ZqIssE800chJE96qwSDp3Vepz0omV4V3rUK7Qa5/TZtEEgnhFs0uShOc01qwBfIRqW4VPPm3Z69TNjjVKDKKygZIYgG8hIOWNRG9pM4LcmMGO1NkKvfjG0pW+3yF+WkoicdVCwK7h+haKjOppoey/aNSrIoZdrgveDVuIJpZLHo2PiFvKrMZuEHf2rrX7PpCJzwuPpaESFZtsf2flz+OwM1KmlSz2W5tJdwA/PgEuiYygNeljD2JpAkpuiQxxpkvoI7K3zu9WqOBOqmCmOyGJO2CKqgP4eAwWnO4pIjDm6E0IDKCKWUdqIojbjUTfmE4KL0iqF+C3QIfTe6g/y7efl66K4Rz4MwECb6iHO2/OHswlmTlUJ5VbOnt7uI9+iTmNIf0RhydlVBpzFKGBQ8d59ntgTU4EzwzXZ6bFu3IQW1qgAU1KT5AZ7KfJoLKGVgtql/l2IrGG3OTY8L4t11pFR82bue3lJ8fS/LzElLaIQXw1TFh3Sx/IzXy+jtu78Jft9V4Q2yWGcJ+jlh7CNJ+e0iy1I7AgDdXuqsYc2goGHNsGBhhw4nVDIXFXqRuO47zOVMIyBBAyPpGjNLoZNIfrdY95TcdflqU0WG0Mb7X9GeWK6L/PlBgk/d3QMSE8P7Pw== mpdroog@icloud.com

locksmith:
  window_start: "10:00"
  window_length: "1h"
update:
  reboot-strategy: "best-effort"

write_files:
  - path: /etc/ssh/sshd_config
    permissions: 0600
    owner: root:root
    content: |
      # Use most defaults for sshd configuration.
      UsePrivilegeSeparation sandbox
      Subsystem sftp internal-sftp

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
      -A INPUT -i eth1 -j ACCEPT
      -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
      -A INPUT -p tcp -m tcp --dport 2017 -j ACCEPT
      -A INPUT -p tcp -m tcp --dport 80 -j ACCEPT
      -A INPUT -p tcp -m tcp --dport 443 -j ACCEPT
      -A INPUT -p icmp -m icmp --icmp-type 0 -j ACCEPT
      -A INPUT -p icmp -m icmp --icmp-type 3 -j ACCEPT
      -A INPUT -p icmp -m icmp --icmp-type 11 -j ACCEPT
      COMMIT

coreos:
  units:
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
    enable: true
EOF

sudo coreos-install -d /dev/vda -c cloud-config.yaml
sudo reboot`

	ip2Host := map[string]string{
		"127.0.0.1": "test",

		"45.77.53.209": "fra1",
		"82.196.0.212": "ams1",
		"104.236.123.169": "nyc1",
	}
	idx := strings.LastIndex(r.RemoteAddr, ":")
	if idx == -1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Unsupported IP."))
		fmt.Printf("cloudinit: unparsable IP=%s\n", r.RemoteAddr)
		return
	}

	hostname := ip2Host[r.RemoteAddr[:idx]]
	if hostname == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Unsupported IP."))
		fmt.Printf("cloudinit: unsupported IP=%s\n", r.RemoteAddr)
		return
	}
	hostname = hostname + ".rootdev.nl"

	raw = strings.Replace(raw, "{hostname}", hostname, 1)
	w.Write([]byte(raw))
}

func main() {
	listen := ":80"
	flag.BoolVar(&verbose, "v", false, "Verbose-mode (log more)")
	flag.StringVar(&listen, "l", ":80", "Listener port")
	flag.Parse()

	mailgunDomain = os.Getenv("MAILGUN_DOMAIN")
	mailgunApi = os.Getenv("MAILGUN_APIKEY")
	mailgunApiPub = os.Getenv("MAILGUN_PUBLICAPIKEY")

	if mailgunDomain == "" {
		fmt.Printf("ERR:main:env: Missing MAILGUN_DOMAIN\n")
		return
	}
	if mailgunApi == "" {
		fmt.Printf("ERR:main:env: Missing MAILGUN_APIKEY\n")
		return
	}
	if mailgunApiPub == "" {
		fmt.Printf("ERR:main:env: Missing MAILGUN_PUBLICAPIKEY\n")
		return
	}

	fs := http.FileServer(http.Dir("build"))
	http.Handle("/", fs)
	http.HandleFunc("/action/ipxe", ipxe)
	http.HandleFunc("/action/email", email)

	if verbose {
		fmt.Printf("Listening on %s\n", listen)
	}
	if e := http.ListenAndServe(listen, nil); e != nil {
		panic(e)
	}
}
