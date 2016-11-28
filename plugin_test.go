package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/docker/infrakit/pkg/spi/instance"
	"github.com/profitbricks/profitbricks-sdk-go"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
)

var (
	dir           = "./"
	tags          = map[string]string{"group": "test"}
	instanceProps = json.RawMessage(`
	 {
        "ProfitBricksInstancesInput": {
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
          "Firewall" :{
            "Name" : "firewall name",
            "Protocol" : "TCP"
          }
        },
        "Tags": {
          "Name": "infrakit-example"
        }
      }`)
	serverIn = `
{
  "id" : "ebb0018f-cd19-4d53-8f33-d219055c0d4b",
  "type" : "datacenter",
  "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b",
  "metadata" : {
    "createdDate" : "2016-11-22T23:24:06Z",
    "createdBy" : "jasmin@stackpointcloud.com",
    "etag" : "18e75b53fa94763a744ed3916a2d6b5b",
    "lastModifiedDate" : "2016-11-22T23:29:55Z",
    "lastModifiedBy" : "jasmin@stackpointcloud.com",
    "state" : "AVAILABLE"
  },
  "properties" : {
    "name" : "datacenter 01",
    "description" : "description of the datacenter",
    "location" : "us/las",
    "version" : 4,
    "features" : [ "SSD", "MULTIPLE_CPU" ]
  },
  "entities" : {
    "servers" : {
      "id" : "ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers",
      "type" : "collection",
      "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers",
      "items" : [ {
        "id" : "74d05f49-03f8-4564-b822-ff9835a5026d",
        "type" : "server",
        "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d",
        "metadata" : {
          "createdDate" : "2016-11-22T23:24:23Z",
          "createdBy" : "jasmin@stackpointcloud.com",
          "etag" : "18e75b53fa94763a744ed3916a2d6b5b",
          "lastModifiedDate" : "2016-11-22T23:29:55Z",
          "lastModifiedBy" : "jasmin@stackpointcloud.com",
          "state" : "AVAILABLE"
        },
        "properties" : {
          "name" : "webserver",
          "cores" : 1,
          "ram" : 1024,
          "availabilityZone" : "ZONE_1",
          "vmState" : "RUNNING",
          "bootCdrom" : null,
          "bootVolume" : {
            "id" : "f9635d5c-cbf1-4fa2-baf9-f07f2a161ffe",
            "type" : "volume",
            "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/volumes/f9635d5c-cbf1-4fa2-baf9-f07f2a161ffe",
            "metadata" : {
              "createdDate" : "2016-11-22T23:24:23Z",
              "createdBy" : "jasmin@stackpointcloud.com",
              "etag" : "61e099a72c9d99a69cc141407e94241e",
              "lastModifiedDate" : "2016-11-22T23:24:23Z",
              "lastModifiedBy" : "jasmin@stackpointcloud.com",
              "state" : "AVAILABLE"
            },
            "properties" : {
              "name" : "system",
              "type" : "SSD",
              "size" : 5,
              "image" : "deffa7ac-9ff1-11e6-a389-52540005ab80",
              "imagePassword" : null,
              "bus" : "VIRTIO",
              "licenceType" : "LINUX",
              "cpuHotPlug" : true,
              "cpuHotUnplug" : false,
              "ramHotPlug" : true,
              "ramHotUnplug" : false,
              "nicHotPlug" : true,
              "nicHotUnplug" : true,
              "discVirtioHotPlug" : true,
              "discVirtioHotUnplug" : true,
              "discScsiHotPlug" : false,
              "discScsiHotUnplug" : false,
              "deviceNumber" : 1
            }
          },
          "cpuFamily" : "AMD_OPTERON"
        },
        "entities" : {
          "cdroms" : {
            "id" : "74d05f49-03f8-4564-b822-ff9835a5026d/cdroms",
            "type" : "collection",
            "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/cdroms",
            "items" : [ ]
          },
          "volumes" : {
            "id" : "74d05f49-03f8-4564-b822-ff9835a5026d/volumes",
            "type" : "collection",
            "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/volumes",
            "items" : [ {
              "id" : "f9635d5c-cbf1-4fa2-baf9-f07f2a161ffe",
              "type" : "volume",
              "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/volumes/f9635d5c-cbf1-4fa2-baf9-f07f2a161ffe",
              "metadata" : {
                "createdDate" : "2016-11-22T23:24:23Z",
                "createdBy" : "jasmin@stackpointcloud.com",
                "etag" : "61e099a72c9d99a69cc141407e94241e",
                "lastModifiedDate" : "2016-11-22T23:24:23Z",
                "lastModifiedBy" : "jasmin@stackpointcloud.com",
                "state" : "AVAILABLE"
              },
              "properties" : {
                "name" : "system",
                "type" : "SSD",
                "size" : 5,
                "availabilityZone" : "AUTO",
                "image" : "deffa7ac-9ff1-11e6-a389-52540005ab80",
                "imagePassword" : null,
                "sshKeys" : null,
                "bus" : "VIRTIO",
                "licenceType" : "LINUX",
                "cpuHotPlug" : true,
                "cpuHotUnplug" : false,
                "ramHotPlug" : true,
                "ramHotUnplug" : false,
                "nicHotPlug" : true,
                "nicHotUnplug" : true,
                "discVirtioHotPlug" : true,
                "discVirtioHotUnplug" : true,
                "discScsiHotPlug" : false,
                "discScsiHotUnplug" : false,
                "deviceNumber" : 1
              }
            } ]
          },
          "nics" : {
            "id" : "74d05f49-03f8-4564-b822-ff9835a5026d/nics",
            "type" : "collection",
            "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics",
            "items" : [ {
              "id" : "99756cf4-b512-4fb6-a278-97a0b0254f00",
              "type" : "nic",
              "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/99756cf4-b512-4fb6-a278-97a0b0254f00",
              "metadata" : {
                "createdDate" : "2016-11-22T23:29:55Z",
                "createdBy" : "jasmin@stackpointcloud.com",
                "etag" : "18e75b53fa94763a744ed3916a2d6b5b",
                "lastModifiedDate" : "2016-11-22T23:29:55Z",
                "lastModifiedBy" : "jasmin@stackpointcloud.com",
                "state" : "AVAILABLE"
              },
              "properties" : {
                "name" : null,
                "mac" : "02:01:d2:b1:1f:54",
                "ips" : [ "192.152.28.140" ],
                "dhcp" : true,
                "lan" : 1,
                "firewallActive" : false,
                "nat" : false
              },
              "entities" : {
                "firewallrules" : {
                  "id" : "99756cf4-b512-4fb6-a278-97a0b0254f00/firewallrules",
                  "type" : "collection",
                  "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/99756cf4-b512-4fb6-a278-97a0b0254f00/firewallrules",
                  "items" : [ ]
                }
              }
            }, {
              "id" : "92a2de68-dabb-4f48-b689-2c83ea907a5c",
              "type" : "nic",
              "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/92a2de68-dabb-4f48-b689-2c83ea907a5c",
              "metadata" : {
                "createdDate" : "2016-11-22T23:24:23Z",
                "createdBy" : "jasmin@stackpointcloud.com",
                "etag" : "61e099a72c9d99a69cc141407e94241e",
                "lastModifiedDate" : "2016-11-22T23:24:23Z",
                "lastModifiedBy" : "jasmin@stackpointcloud.com",
                "state" : "AVAILABLE"
              },
              "properties" : {
                "name" : null,
                "mac" : "02:01:aa:9b:94:7e",
                "ips" : [ "162.254.25.254" ],
                "dhcp" : true,
                "lan" : 1,
                "firewallActive" : true,
                "nat" : false
              },
              "entities" : {
                "firewallrules" : {
                  "id" : "92a2de68-dabb-4f48-b689-2c83ea907a5c/firewallrules",
                  "type" : "collection",
                  "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/92a2de68-dabb-4f48-b689-2c83ea907a5c/firewallrules",
                  "items" : [ {
                    "id" : "a94f50fa-0fca-4ca3-bd29-085af0ece49e",
                    "type" : "firewall-rule",
                    "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/92a2de68-dabb-4f48-b689-2c83ea907a5c/firewallrules/a94f50fa-0fca-4ca3-bd29-085af0ece49e"
                  } ]
                }
              }
            } ]
          }
        }
      } ]
    },
    "volumes" : {
      "id" : "ebb0018f-cd19-4d53-8f33-d219055c0d4b/volumes",
      "type" : "collection",
      "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/volumes",
      "items" : [ {
        "id" : "f9635d5c-cbf1-4fa2-baf9-f07f2a161ffe",
        "type" : "volume",
        "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/volumes/f9635d5c-cbf1-4fa2-baf9-f07f2a161ffe",
        "metadata" : {
          "createdDate" : "2016-11-22T23:24:23Z",
          "createdBy" : "jasmin@stackpointcloud.com",
          "etag" : "61e099a72c9d99a69cc141407e94241e",
          "lastModifiedDate" : "2016-11-22T23:24:23Z",
          "lastModifiedBy" : "jasmin@stackpointcloud.com",
          "state" : "AVAILABLE"
        },
        "properties" : {
          "name" : "system",
          "type" : "SSD",
          "size" : 5,
          "availabilityZone" : "AUTO",
          "image" : "deffa7ac-9ff1-11e6-a389-52540005ab80",
          "imagePassword" : null,
          "sshKeys" : null,
          "bus" : "VIRTIO",
          "licenceType" : "LINUX",
          "cpuHotPlug" : true,
          "cpuHotUnplug" : false,
          "ramHotPlug" : true,
          "ramHotUnplug" : false,
          "nicHotPlug" : true,
          "nicHotUnplug" : true,
          "discVirtioHotPlug" : true,
          "discVirtioHotUnplug" : true,
          "discScsiHotPlug" : false,
          "discScsiHotUnplug" : false,
          "deviceNumber" : 1
        }
      } ]
    },
    "loadbalancers" : {
      "id" : "ebb0018f-cd19-4d53-8f33-d219055c0d4b/loadbalancers",
      "type" : "collection",
      "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/loadbalancers",
      "items" : [ ]
    },
    "lans" : {
      "id" : "ebb0018f-cd19-4d53-8f33-d219055c0d4b/lans",
      "type" : "collection",
      "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/lans",
      "items" : [ {
        "id" : "1",
        "type" : "lan",
        "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/lans/1",
        "metadata" : {
          "createdDate" : "2016-11-22T23:24:09Z",
          "createdBy" : "jasmin@stackpointcloud.com",
          "etag" : "b189f9e5d2014cfb9c44d565cc5a6a56",
          "lastModifiedDate" : "2016-11-22T23:24:09Z",
          "lastModifiedBy" : "jasmin@stackpointcloud.com",
          "state" : "AVAILABLE"
        },
        "properties" : {
          "name" : "public",
          "public" : true
        },
        "entities" : {
          "nics" : {
            "id" : "1/nics",
            "type" : "collection",
            "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/lans/1/nics",
            "items" : [ {
              "id" : "99756cf4-b512-4fb6-a278-97a0b0254f00",
              "type" : "nic",
              "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/99756cf4-b512-4fb6-a278-97a0b0254f00",
              "metadata" : {
                "createdDate" : "2016-11-22T23:29:55Z",
                "createdBy" : "jasmin@stackpointcloud.com",
                "etag" : "18e75b53fa94763a744ed3916a2d6b5b",
                "lastModifiedDate" : "2016-11-22T23:29:55Z",
                "lastModifiedBy" : "jasmin@stackpointcloud.com",
                "state" : "AVAILABLE"
              },
              "properties" : {
                "name" : null,
                "mac" : "02:01:d2:b1:1f:54",
                "ips" : [ "192.152.28.140" ],
                "dhcp" : true,
                "lan" : 1,
                "firewallActive" : false,
                "nat" : false
              },
              "entities" : {
                "firewallrules" : {
                  "id" : "99756cf4-b512-4fb6-a278-97a0b0254f00/firewallrules",
                  "type" : "collection",
                  "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/99756cf4-b512-4fb6-a278-97a0b0254f00/firewallrules",
                  "items" : [ ]
                }
              }
            }, {
              "id" : "92a2de68-dabb-4f48-b689-2c83ea907a5c",
              "type" : "nic",
              "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/92a2de68-dabb-4f48-b689-2c83ea907a5c",
              "metadata" : {
                "createdDate" : "2016-11-22T23:24:23Z",
                "createdBy" : "jasmin@stackpointcloud.com",
                "etag" : "61e099a72c9d99a69cc141407e94241e",
                "lastModifiedDate" : "2016-11-22T23:24:23Z",
                "lastModifiedBy" : "jasmin@stackpointcloud.com",
                "state" : "AVAILABLE"
              },
              "properties" : {
                "name" : null,
                "mac" : "02:01:aa:9b:94:7e",
                "ips" : [ "162.254.25.254" ],
                "dhcp" : true,
                "lan" : 1,
                "firewallActive" : true,
                "nat" : false
              },
              "entities" : {
                "firewallrules" : {
                  "id" : "92a2de68-dabb-4f48-b689-2c83ea907a5c/firewallrules",
                  "type" : "collection",
                  "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/92a2de68-dabb-4f48-b689-2c83ea907a5c/firewallrules",
                  "items" : [ {
                    "id" : "a94f50fa-0fca-4ca3-bd29-085af0ece49e",
                    "type" : "firewall-rule",
                    "href" : "https://api.profitbricks.com/cloudapi/v3/datacenters/ebb0018f-cd19-4d53-8f33-d219055c0d4b/servers/74d05f49-03f8-4564-b822-ff9835a5026d/nics/92a2de68-dabb-4f48-b689-2c83ea907a5c/firewallrules/a94f50fa-0fca-4ca3-bd29-085af0ece49e"
                  } ]
                }
              }
            } ]
          }
        }
      } ]
    }
  }
}`
)

func TestInstanceLifecycle(t *testing.T) {

	prop := createInstance{}

	require.NoError(t, json.Unmarshal([]byte(instanceProps), &prop))

	p := &fakeplugin{
		username: "u",
		password: "p",
		dir:      dir,
	}

	instanceID, err := p.Provision(instance.Spec{Properties: &instanceProps, Tags: tags})
	require.NoError(t, err)

	var buff []byte
	buff, err = ioutil.ReadFile(path.Join(p.dir, string(*instanceID)+filesufix))
	require.NoError(t, err)

	desc := instance.Description{}
	err = json.Unmarshal(buff, &desc)
	require.NoError(t, err)
	require.Equal(t, string(*instanceID), string(desc.ID))

	require.NoError(t, p.Destroy(*instanceID))
}

func TestPlugin_Destroy(t *testing.T) {
	p := &fakeplugin{
		username: "u",
		password: "p",
		dir:      dir,
	}
	var dc profitbricks.Datacenter
	err := json.Unmarshal([]byte(serverIn), &dc)

	var id instance.ID
	id = instance.ID(dc.Id)

	err = p.Destroy(id)

	require.Error(t, err)
}

func TestPlugin_Validate(t *testing.T) {
	p := &fakeplugin{
		username: "u",
		password: "p",
		dir:      dir,
	}

	require.NoError(t, p.Validate(instanceProps))

}

type fakeplugin struct {
	username string
	password string
	dir      string
}

func (p *fakeplugin) Validate(req json.RawMessage) error {
	properties := createInstance{}
	if err := json.Unmarshal([]byte(req), &properties); err != nil {
		return fmt.Errorf("Invalid instance properties: %s", err)
	}

	if !pathExists(p.dir) {
		return fmt.Errorf("Directoy % doesn't exist", p.dir)
	}
	profitbricks.SetAuth(p.username, p.password)

	locations := profitbricks.ListLocations()

	if locations.StatusCode == 403 {
		return fmt.Errorf("ProfitBricks credentials you providedare incorrect")
	}

	if properties.ProfitBricksInstancesInput.Name == "" {
		return fmt.Errorf("Name parameter is required")
	}

	if properties.ProfitBricksInstancesInput.Location == "" {
		return fmt.Errorf("Location parameter is required")
	}

	if properties.ProfitBricksInstancesInput.DiskType == "" {
		return fmt.Errorf("DiskType is required")
	}

	if properties.ProfitBricksInstancesInput.Image == "" {
		return fmt.Errorf("Image parameter is required")
	}

	if properties.ProfitBricksInstancesInput.SSHKeyPath == "" {
		return fmt.Errorf("Path to ssh key is required")
	}

	if properties.ProfitBricksInstancesInput.Location == "" {
		return fmt.Errorf("Location is required")
	}

	if properties.ProfitBricksInstancesInput.DiskSize == 0 {
		return fmt.Errorf("DiskSize is required")
	}

	if properties.ProfitBricksInstancesInput.Cores == 0 {
		return fmt.Errorf("Cores is required")
	}

	if properties.ProfitBricksInstancesInput.RAM == 0 {
		return fmt.Errorf("RAM is required")
	}

	if properties.ProfitBricksInstancesInput.Firewall != nil {
		if properties.ProfitBricksInstancesInput.Firewall.Protocol == "" {
			return fmt.Errorf("Firewall protocol is required")
		}
	}

	return nil
}

func (p *fakeplugin) Provision(spec instance.Spec) (*instance.ID, error) {
	var dc profitbricks.Datacenter

	pbdata := &createInstance{}

	err := json.Unmarshal(*spec.Properties, pbdata)
	//instead of create
	err = json.Unmarshal([]byte(serverIn), &dc)

	p.recordInstanceState(spec, &dc, pbdata)
	if err != nil {
		return nil, fmt.Errorf("Invalid input formatting: %s", err)
	}

	var id instance.ID
	id = instance.ID(dc.Id)

	return &id, nil
}

func (p fakeplugin) DescribeInstances(tags map[string]string) ([]instance.Description, error) {
	return nil, nil
}

func (p fakeplugin) Destroy(id instance.ID) error {
	err := p.removeFile(string(id) + filesufix)
	if err != nil {
		return err
	}

	return nil
}

func (p *fakeplugin) recordInstanceState(spec instance.Spec, dc *profitbricks.Datacenter, pbdata *createInstance) error {
	if !pathExists(p.dir) {
		return fmt.Errorf("Directory %s doesn't exist", p.dir)
	}

	id := instance.ID(dc.Id)
	logicalID := instance.LogicalID(dc.Entities.Servers.Items[0].Entities.Nics.Items[0].Properties.Ips[0])
	description := instance.Description{
		Tags:      spec.Tags,
		ID:        id,
		LogicalID: &logicalID,
	}
	towrite, _ := json.MarshalIndent(description, "", "\t")

	filePath := path.Join(p.dir, dc.Id+filesufix)

	err := ioutil.WriteFile(filePath, towrite, 0644)
	if err != nil {
		return fmt.Errorf("An error occurred while trying to write to file %s, %s", filePath, err.Error())
	}
	return nil
}

func (p *fakeplugin) removeFile(file string) error {
	return os.Remove(path.Join(p.dir, file))
}
