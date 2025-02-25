
# Packetbeat

Packetbeat is an open source network packet analyzer that ships the data to
Elasticsearch. Think of it like a distributed real-time Wireshark with a lot
more analytics features.

The Packetbeat shippers sniff the traffic between your application processes,
parse on the fly protocols like HTTP, MySQL, PostgreSQL, Redis or Thrift and
correlate the messages into transactions.

For each transaction, the shipper inserts a JSON document into Elasticsearch,
where it is stored and indexed. You can then use Kibana to view key metrics and
do ad-hoc queries against the data.

To learn more about Packetbeat, check out <https://www.elastic.co/beats/packetbeat>.

## Getting started

Please follow the [getting started](https://www.elastic.co/guide/en/beats/packetbeat/current/packetbeat-installation-configuration.html)
guide from the docs.


```bash

dpdk >= DPDK 20.02.1

kernel >= 3.10.0

CentOS
#  yum install -y libpcap-devel gcc gcc-c++ make meson ninja  numactl-devel  numactl  net-tools pciutils
#  yum install -y kernel-devel-$(uname -r) kernel-headers-$(uname -r)

Debian + Ubuntu
# apt install -y libpcap-dev gcc g++ make meson ninja-build libnuma-dev numactl net-tools pciutils
# apt install -y linux-headers-$(uname -r)


#  wget http://fast.dpdk.org/rel/dpdk-20.11.10.tar.xz
#  tar -Jxvf dpdk-20.11.10.tar.xz
#  cd dpdk-stable-20.11.10 && meson build && cd build && ninja && ninja install
#  export LD_LIBRARY_PATH=/usr/local/lib64:$LD_LIBRARY_PATH
#  git clone git://dpdk.org/dpdk-kmods && cd  dpdk-kmods/linux/igb_uio
#  make
#  modprobe uio  &&  insmod igb_uio.ko
#  dpdk-devbind.py -b igb_uio 0000:03:00.0(pci-addr)
#  go clean -modcache && go mod tidy
#  CGO_CFLAGS="-msse4.2 -fno-strict-aliasing " CGO_LDFLAGS=" -lrte_eal -lrte_mbuf -lrte_mempool -lrte_ethdev -lpcap" go build
#  ./packetbeat8_dpdk --dpdk_status enable --dpdk_port 0 -c ~/go/packetbeat.dpdk.yml

```

