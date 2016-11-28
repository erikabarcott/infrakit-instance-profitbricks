# InfraKit Instance Plugin for ProfitBricks

[![Build Status](https://travis-ci.org/StackPointCloud/infrakit-instance-profitbricks.svg?branch=master)](https://travis-ci.org/StackPointCloud/infrakit-instance-profitbricks)

This is an [InfraKit](https://github.com/docker/infrakit) instance plugin for creating and managing ProfitBricks servers.
The plugin development is work in progress. Infrakit is also still being developed and is constantly changing. 

## Usage

Start with building the InfraKit [binaries](https://github.com/docker/infrakit/blob/master/README.md#binaries).
Currently, you can use this plugin with a plain [vanilla flavor plugin](https://github.com/docker/infrakit/tree/master/pkg/example/flavor/vanilla) and the [default group plugin](https://github.com/docker/infrakit/blob/master/cmd/group/README.md).

To build the ProfitBricks Instance plugin, follow next steps:

1. Install [GO](https://golang.org/)  
2. Get the source code `go get github.com/profitbricks/infrakit-instance-profitbricks` 
3. Build the binaries. In root of the folder run `make` 

Use the help command to list the command line options available with the plugin.

```shell
$ ./infrakit-instance-profitbricks --help
ProfitBricks instance plugin

Usage:
  ./infrakit-instance-profitbricks [flags]
  ./infrakit-instance-profitbricks [command]

Available Commands:
  version     print build version information
  version     print build version information

Flags:
      --dir string        Existing directory for storing the plugin files (default "/var/folders/by/8dsj76cs15jgx_fvjltwhpg40000gn/T/")
      --log int           Logging level. 0 is least verbose. Max is 5 (default 4)
      --name string       Plugin name to advertise for discovery (default "infrakit-instance-profitbricks")
      --password string   ProfitBricks username
      --username string   ProfitBricks username

Use "./infrakit-instance-profitbricks [command] --help" for more information about a command.

```

Run the plugin:

```shell
$ ./infrakit-instance-profitbricks 
INFO[0000] Listening at: /Users/jasmingacic/.infrakit/plugins/infrakit-instance-profitbricks 
```

Note that `--password` and `--username` are required, if not provided environment variables `PROFITBRICKS_USERNAME` and `PROFITBRICKS_PASSWORD` are expected to be set.

From the InfraKit build directory run:

```shell
$ build/infrakit-group-default
```

```shell
$ build/infrakit-flavor-vanilla
```

Use the provided configuration example [pbexample.json](./pbexample.json) as a reference and feel free to change 
the values of the properties.

```shell
$ cat << EOF > pb.json
{
  "ID": "pb-example",
  "Properties": {
    "Allocation": {
      "Size": 1
    },
    "Instance": {
      "Plugin": "infrakit-instance-profitbricks",
      "Properties": {
        "ProfitBricksInstancesInput": {
          "DatacenterId": "aafa9ac4-e05b-4f84-b5cb-f69d0903f923",
          "Name": "Infrakit test01",
          "Image": "Ubuntu-16.04",
          "SSHKeyPath": "/Users/jasmingacic/.ssh/id_rsa",
          "Location": "us/las",
          "Cores": 1,
          "Ram": 1024,
          "DiskSize": 5,
          "DiskType": "SSD",
          "AvailabilityZone": "AUTO",
          "StaticIP": false,
          "Firewall": {
            "Name": "firewall name",
            "Protocol": "TCP"
          }
        },
        "Tags": {
          "Name": "infrakit-example"
        }
      }
    },
    "Flavor": {
      "Plugin": "flavor-vanilla",
      "Properties": {
        "Init": [
          "sh -c \"echo 'Hello, World!' > /hello\""
        ]
      }
    }
  }
}
EOF
```

Commit the configuration by running the InfraKit command:

```shell
$ build/infrakit group commit pb.json
Committed myGroup: Managing 1 instances
```

### Managing groups

You can use the set of InfraKit commands to manage the groups being monitored.

```
$ build/infrakit group --help
Access group plugin

Usage:
  build/infrakit group [command]

Available Commands:
  commit      commit a group configuration
  describe    describe the live instances that make up a group
  destroy     destroy a group
  free        free a group from active monitoring, nondestructive
  inspect     return the raw configuration associated with a group
  ls          list groups

Flags:
      --name string   Name of plugin (default "group")

Global Flags:
      --log int   Logging level. 0 is least verbose. Max is 5 (default 4)

Use "build/infrakit group [command] --help" for more information about a command.
```

Describe group command displays info about the instances in a tabular form.

```
$ ./infrakit group describe pb-example
ID                            	LOGICAL                       	TAGS
c1343579-d3d7-455e-be89-ebe114137231	162.254.27.219                	datacenter_id=aafa9ac4-e05b-4f84-b5cb-f69d0903f923,infrakit.config_sha=cwlm7Jm3K2oEZE-AhXCWs3PQESQ=,infrakit.group=pb-example
```

If you would like to increase or decrease the number of instances in the group, modify the allocation `Size` property and commit the modified configuration.

```
$ build/infrakit group commit pb.json
Committed pb-example: Adding 1 instances to increase the group size to 2
```
```
$ build/infrakit group commit pb.json
Committed pb-example: Terminating 2 instances to reduce the group size to 1
```

Run destroy command to terminate a group monitoring and delete all instances, i.e., servers in the group.

```
$ build/infrakit group destroy myGroup
destroy pb-example initiated
```

## Design notes

The plugin stores a basic info about the instances onto the provided location (`--dir`). You can stop and start a group monitoring without redeploying the servers.

The instance file names consist of the instance (server) UUID and `.pbstate` extension.

## Configuration parameters

Required parameters:

```
    DatacenterID     - (string) UUID of Virtual Data Center   
    Name             - (string) Name of the servers to be provisioned   
    Image            - (string) ProfitBricks volume image   
    SSHKeyPath       - (string) Path to private SSHkey.    
    Location         - (string) ProfitBricks location 
    DiskSize         - (int) Desired disk size  
    DiskType         - (string) Desired disk type  
    Cores            - (int) Number of cores per server      
    RAM              - (int) Size of RAM in MB
```

Optional parameters:
```
    StaticIP         - (bool) Flag indicating if the servers will be provisoned with a static IP addres
    AvailabilityZone - (string) Server and Volume availability zone
    Firewall         - (object) Firewall Rule object to be provisoned under the server.
```

Firewall parameters

Required Firewall Rule parameters:
```
    Protocol       - (string) The protocol for the rule: TCP, UDP, ICMP, ANY.
```

Optional Firewall Rule parameters:
```
    
    Name           - (string) Name of the Firewall rule
    SourceMac      - (string) Only traffic originating from the respective MAC address is allowed. 
    SourceIP       - (string) Only traffic originating from the respective IPv4 address is allowed.
    TargetIP       - (string) In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed.
    IcmpCode       - (int) Defines the allowed code (from 0 to 254) if protocol ICMP is chosen. 
    IcmpType       - (int) Defines the allowed type (from 0 to 254) if the protocol ICMP is chosen.   
    PortRangeStart - (int) Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen.   
    PortRangeEnd   - (int) Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen.   
```