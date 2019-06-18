#### Hp ilo exporter
Simple exporter for gathering some metrics from hp ilo interfaces through web api.
Have inside auto-update of cookies.

Listen port is `9600` by default. You can simply change this value by flags interface:
```
./hpilo-exporter --help
Usage of ./hpilo-exporter:
  -hosts-file string
        full path to hosts file (default "hosts.yaml")
  -listen-addr string
        Define your listen address if needed (default "0.0.0.0")
  -listen-port string
        Listen port (default "9600")
```

Hosts file is used for get list of ilo targets for gather information.

As writed in hosts.yaml you can define global login/pass for every ilo who have not own different user/pass:
```
global:
  login: root
  password: toor
hosts:
  - address: 192.168.0.2
    name: ilo1
    password: root
  - address: 192.168.0.3
    name: ilo2
    login: test
    password: test
  - address: 192.168.0.4
    name: ilo3
  - address: 192.168.0.5
    name: ilo5
    login: admin
```