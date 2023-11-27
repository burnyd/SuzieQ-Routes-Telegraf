### Telegraf Suzie integration

This repo exists to create a layer between SuzieQ's API and Influx to further graf within Grafana.

##### The following are needed
```bash
- My fork of [Telegraf.](https://github.com/burnyd/telegraf/tree/suzieq)
- [Containerlab](https://containerlab.dev/) Currently using v0.48.3.
- Container EOS.  Currently using 4.29.2F But I dont think it matters as long as its somewhat modern-ish.
- go 1.18 or better.
```

##### Clone the telegraf repo

Clone the repo somewhere this is only going to be needed to build the telegraf binary

```bash
git clone -b suzieq https://github.com/burnyd/telegraf && cd telegraf
```

Run the make file by simply typing in make.

```bash
➜  telegraf git:(suzieq) make
go mod download -x
env -u GOOS -u GOARCH -u GOARM -- go build -o ./tools/readme_config_includer/generator ./tools/readme_config_includer/generator.go
go generate -run="readme_config_includer/generator$" ./plugins/inputs/...
go generate -run="readme_config_includer/generator$" ./plugins/outputs/...
go generate -run="readme_config_includer/generator$" ./plugins/processors/...
go generate -run="readme_config_includer/generator$" ./plugins/aggregators/...
go build -ldflags " -X main.commit=65b6c961 -X main.branch=suzieq -X main.goos=linux -X main.goarch=amd64 -X main.version=1.24.0-65b6c961" ./cmd/telegraf
```

Now there should be a ./telegraf binary within the system.  Simply move it to the current directory of this repo.

```bash
cp telegraf /pathtoSuzieQ-Routes-Telegraf/
```

##### Bring up containerlab / start the lab fun.

```bash
sudo containerlab -t clab.yml deploy
```

Output

```bash

+---+-----------------+--------------+------------------------------------+-------+---------+------------------+----------------------+
| # |      Name       | Container ID |               Image                | Kind  |  State  |   IPv4 Address   |     IPv6 Address     |
+---+-----------------+--------------+------------------------------------+-------+---------+------------------+----------------------+
| 1 | clab-sq-ceos1   | 57fcadae4f87 | gcr.io/act-sandbox/ceoslab:4.29.2F | ceos  | running | 172.20.20.100/24 | 2001:172:20:20::2/64 |
| 2 | clab-sq-ceos2   | 4f8ee2f6379d | gcr.io/act-sandbox/ceoslab:4.29.2F | ceos  | running | 172.20.20.101/24 | 2001:172:20:20::3/64 |
| 3 | clab-sq-grafana | 901dae1b98bf | grafana/grafana:7.4.3              | linux | running | 172.20.20.103/24 | 2001:172:20:20::4/64 |
| 4 | clab-sq-influx  | 7cd2dd006916 | influxdb:2.0                       | linux | running | 172.20.20.104/24 | 2001:172:20:20::6/64 |
| 5 | clab-sq-sq      | 3cfbef17a2d2 | netenglabs/suzieq                  | linux | running | 172.20.20.102/24 | 2001:172:20:20::5/64 |
+---+-----------------+--------------+------------------------------------+-------+---------+------------------+----------------------+
```

##### Run the rest service and poller within SuzieQ

```bash
docker exec -it clab-sq-sq bash
```

Run the poller command and kick off the rest api within sq.

```bash
sq-poller -I /tmp/inventory.yaml & sq-rest-server &
```

Before running telegraf this is a small testing binary I have setup. Which does the same thing as what telegraf would do.

```bash
➜  SuzieQ-Routes-Telegraf git:(main) ✗ go run main.go
Prefix  10.0.0.0/24     connected       On Host         ceos2
Prefix  2.2.2.2/32      connected       On Host         ceos2
Prefix  1.1.1.1/32      ebgp    On Host         ceos2
Prefix  20.0.0.0/24     connected       On Host         ceos2
Prefix  172.20.20.0/24  connected       On Host         ceos2
Prefix  30.0.0.0/24     connected       On Host         ceos2
Prefix  100.100.100.100/32      connected       On Host         ceos2
Prefix  9.9.9.99/32     ebgp    On Host         ceos2
Prefix  0.0.0.0/0       static  On Host         ceos1
Prefix  30.0.0.0/24     connected       On Host         ceos1
Prefix  20.0.0.0/24     connected       On Host         ceos1
Prefix  10.0.0.0/24     connected       On Host         ceos1
Prefix  9.9.9.99/32     connected       On Host         ceos1
Prefix  2.2.2.2/32      ebgp    On Host         ceos1
Prefix  1.1.1.1/32      connected       On Host         ceos1
Prefix  100.100.100.100/32      ebgp    On Host         ceos1
Prefix  172.20.20.0/24  connected       On Host         ceos1
Prefix  0.0.0.0/0       ospfexternaltype2       On Host         ceos2
Prefix  11.22.33.44/32  droproute       On Host         ceos1
Prefix  11.22.33.44/32  ospfexternaltype2       On Host         ceos2
Prefix  22.11.44.33/32  ospfexternaltype2       On Host         ceos1
Prefix  22.11.44.33/32  droproute       On Host         ceos2
```

##### Run Telegraf

Since we moved the binary to this directory and the telegraf config file is within the configs/ portion.

```bash
./telegraf --config configs/telegraf.conf
```

I have this output to stdout just for verbosity / testing.

```
NetworkRoutes,Device=ceos2,OutgoingInterfaces=Management0,RoutingProtocol=connected,Timestamp=1701099314841,VRF=default,host=burnyd Prefix="172.20.20.0/24" 1701099830000000000
```

##### Connect to grafana

Should be 127.0.0.1:3000 unless changed to something else.

Click the four square icon off to the left

![1](/images/1.png)

New dashboard
Where it says "Import via panel json" paste in the contents from configs - > dashboard.json and press load.

![2](/images/2.png)