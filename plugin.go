package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/docker/infrakit/pkg/spi/instance"
	"github.com/profitbricks/profitbricks-sdk-go"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	filesufix = ".pbstate"
)

//InstancePlugin is
func InstancePlugin(username, password, dir string) instance.Plugin {
	return &plugin{
		username: username,
		password: password,
		dir:      dir,
	}
}

type plugin struct {
	username string
	password string
	dir      string
}

type createInstance struct {
	Tags                       map[string]string
	ProfitBricksInstancesInput InstancesInput
}

func (p *plugin) Validate(req json.RawMessage) error {

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

	if properties.ProfitBricksInstancesInput.DatacenterID == "" {
		return fmt.Errorf("Name DatacenterId is required")

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

	imageID := getImageID(properties.ProfitBricksInstancesInput.Image, properties.ProfitBricksInstancesInput.DiskType, properties.ProfitBricksInstancesInput.Location)
	if imageID == "" {
		return fmt.Errorf("Image '%s' does not exist", properties.ProfitBricksInstancesInput.Image)
	}

	if properties.ProfitBricksInstancesInput.SSHKeyPath == "" {
		return fmt.Errorf("Path to ssh key is required")
	}

	if !pathExists(properties.ProfitBricksInstancesInput.SSHKeyPath) {
		return fmt.Errorf("Path '%s' does not exist", properties.ProfitBricksInstancesInput.SSHKeyPath)
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

func (p *plugin) Provision(spec instance.Spec) (*instance.ID, error) {
	pbdata := &createInstance{}

	err := json.Unmarshal(*spec.Properties, pbdata)
	if err != nil {
		return nil, fmt.Errorf("Invalid input formatting: %s", err)
	}
	temp, err := p.createPBMachine(pbdata.ProfitBricksInstancesInput)
	spec.Tags["datacenter_id"] = pbdata.ProfitBricksInstancesInput.DatacenterID
	if err == nil {
		err = p.recordInstanceState(spec, temp, pbdata)
		if err != nil {
			return nil, fmt.Errorf("Error while recoding state", err)
		}
	}

	var id instance.ID
	if temp != nil {
		id = instance.ID(temp.Id)
	}

	return &id, err
}

func (p plugin) DescribeInstances(tags map[string]string) ([]instance.Description, error) {

	return p.getExistingInstances(tags)
}

func (p plugin) Destroy(id instance.ID) error {

	profitbricks.SetAuth(p.username, p.password)
	resp := profitbricks.DeleteDatacenter(string(id))
	err := p.removeFile(string(id))
	if err != nil {
		return err
	}

	if resp.StatusCode != 202 {
		return errors.New("Error while deleting a Virtual DataCenter: " + string(resp.Body))
	}
	return nil
}

//InstancesInput  asdf
type InstancesInput struct {
	DatacenterID     string    `json:"DatacenterId,omitempty"`
	Name             string    `json:"Name,omitempty"`
	Image            string    `json:"Image,omitempty"`
	SSHKeyPath       string    `json:"SSHKeyPath,omitempty"`
	Location         string    `json:"Location,omitempty"`
	DiskSize         int       `json:"DiskSize,omitempty"`
	DiskType         string    `json:"DiskType,omitempty"`
	AvailabilityZone string    `json:"AvailabilityZone,omitempty"`
	StaticIP         bool      `json:"StaticIP,omitempty"`
	Cores            int       `json:"Cores,omitempty"`
	RAM              int       `json:"Ram,omitempty"`
	Firewall         *Firewall `json:"Firewall,omitempty"`
}

//Firewall rule object
type Firewall struct {
	Name           string `json:"Name,omitempty"`
	Protocol       string `json:"Protocol,omitempty"`
	SourceMac      string `json:"SourceMac,omitempty"`
	SourceIP       string `json:"SourceIp,omitempty"`
	TargetIP       string `json:"TargetIp,omitempty"`
	IcmpCode       int    `json:"IcmpCode,omitempty"`
	IcmpType       int    `json:"IcmpType,omitempty"`
	PortRangeStart int    `json:"PortRangeStart,omitempty"`
	PortRangeEnd   int    `json:"PortRangeEnd,omitempty"`
}

func (p *plugin) createPBMachine(input InstancesInput) (*profitbricks.Server, error) {
	profitbricks.SetAuth(p.username, p.password)
	profitbricks.SetDepth("5")
	SSHKey, err := getSSHKey(input.SSHKeyPath)
	if err != nil {
		return nil, err
	}

	imageID := getImageID(input.Image, input.DiskType, input.Location)
	dc := profitbricks.GetDatacenter(input.DatacenterID)

	if dc.StatusCode != 200 {
		return nil, fmt.Errorf("Error occurred while fetching datacenter %s", input.DatacenterID)
	}
	lan := profitbricks.Lan{

		Properties: profitbricks.LanProperties{
			Name:   input.Name,
			Public: true,
		},
	}

	lan = profitbricks.CreateLan(dc.Id, lan)

	if dc.StatusCode > 299 {
		return nil, fmt.Errorf("An error occurred while provisioning a Virtual DataCenter %s", dc.Response)
	}

	err = waitTillProvisioned(lan.Headers.Get("Location"))
	if err != nil {
		return nil, fmt.Errorf("An error occurred while provisioning Lan %s", lan.Response)
	}

	lan = profitbricks.GetLan(dc.Id, lan.Id)
	if lan.StatusCode != 200 {
		return nil, fmt.Errorf("An error occurred while retrieving LANs : %s, %d", lan.Response, lan.StatusCode)
	}

	lanID, _ := strconv.Atoi(lan.Id)
	server := profitbricks.Server{
		Properties: profitbricks.ServerProperties{
			Name:  input.Name,
			Cores: input.Cores,
			Ram:   input.RAM,
		},
		Entities: &profitbricks.ServerEntities{
			Nics: &profitbricks.Nics{
				Items: []profitbricks.Nic{
					{
						Properties: profitbricks.NicProperties{
							Dhcp: true,
							Lan:  lanID,
						},
					},
				},
			},
			Volumes: &profitbricks.Volumes{
				Items: []profitbricks.Volume{
					{
						Properties: profitbricks.VolumeProperties{
							Name:             input.Name,
							Size:             input.DiskSize,
							Type:             input.DiskType,
							SshKeys:          []string{SSHKey},
							AvailabilityZone: input.AvailabilityZone,
							Image:            imageID,
						},
					},
				},
			},
		},
	}

	if input.StaticIP == true {
		req := profitbricks.IpBlock{
			Properties: profitbricks.IpBlockProperties{
				Size:     1,
				Location: input.Location,
			},
		}

		resp := profitbricks.ReserveIpBlock(req)

		if resp.StatusCode != 202 {
			return nil, errors.New(resp.Response)
		}

		waitTillProvisioned(resp.Headers.Get("Location"))

		server.Entities.Nics.Items[0].Properties.Ips = resp.Properties.Ips
	}

	if input.Firewall != nil {
		server.Entities.Nics.Items[0].Properties.FirewallActive = true

		firewall := profitbricks.FirewallRule{
			Properties: profitbricks.FirewallruleProperties{
				Protocol: input.Firewall.Protocol,
			},
		}

		if input.Firewall.Name != "" {
			firewall.Properties.Name = input.Firewall.Name
		}

		if input.Firewall.IcmpCode != 0 {
			firewall.Properties.IcmpCode = input.Firewall.IcmpCode
		}

		if input.Firewall.IcmpType != 0 {
			firewall.Properties.IcmpType = input.Firewall.IcmpType
		}

		if input.Firewall.TargetIP != "" {
			firewall.Properties.TargetIp = input.Firewall.TargetIP
		}

		if input.Firewall.SourceMac != "" {
			firewall.Properties.SourceMac = input.Firewall.SourceMac
		}

		if input.Firewall.SourceIP != "" {
			firewall.Properties.SourceIp = input.Firewall.SourceIP
		}

		if input.Firewall.PortRangeStart != 0 {
			firewall.Properties.PortRangeStart = input.Firewall.PortRangeStart
		}

		if input.Firewall.PortRangeEnd != 0 {
			firewall.Properties.PortRangeEnd = input.Firewall.PortRangeEnd
		}

		server.Entities.Nics.Items[0].Entities = &profitbricks.NicEntities{
			Firewallrules: &profitbricks.FirewallRules{
				Items: []profitbricks.FirewallRule{
					firewall,
				},
			},
		}
	}

	resp := profitbricks.CreateServer(dc.Id, server)

	if resp.StatusCode != 202 {
		return nil, errors.New(resp.Response)
	}

	err = waitTillProvisioned(resp.Headers.Get("Location"))
	if err != nil {
		return nil, fmt.Errorf("An error occurred while provisioning Server %s", err)

	}
	return &resp, nil
}

func waitTillProvisioned(path string) error {
	waitCount := 120
	for i := 0; i < waitCount; i++ {
		request := profitbricks.GetRequestStatus(path)
		if request.Metadata.Status == "DONE" {
			break
		}
		if request.Metadata.Status == "FAILED" {
			return fmt.Errorf(request.Metadata.Message, request.Metadata.Status)
		}
		time.Sleep(1 * time.Second)
		i++
	}

	return nil
}

func getImageID(imageName string, imageType string, location string) string {
	images := profitbricks.ListImages()
	if images.StatusCode > 299 {
		log.Print(fmt.Errorf("Error while fetching the list of images %s", images.Response))
	}

	if len(images.Items) > 0 {
		for _, i := range images.Items {
			imgName := ""
			if i.Properties.Name != "" {
				imgName = i.Properties.Name
			}

			if imageType == "SSD" {
				imageType = "HDD"
			}
			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(imageName)) && i.Properties.ImageType == imageType && i.Properties.Location == location && i.Properties.Public == true {
				return i.Id
			}
		}
	}

	return ""
}

func getSSHKey(path string) (string, error) {
	pemBytes, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(pemBytes)

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return "", fmt.Errorf("Error getting ssh key: %s", err)
	}

	privblk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(priv),
	}

	_, err = ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&privblk)), err
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (p *plugin) recordInstanceState(spec instance.Spec, server *profitbricks.Server, pbdata *createInstance) error {

	if !pathExists(p.dir) {
		return fmt.Errorf("Directory %s doesn't exist", p.dir)
	}

	freshdc := profitbricks.GetDatacenter(spec.Tags["datacenter_id"])
	id := instance.ID(server.Id)

	servertmp := profitbricks.GetServer(freshdc.Id, server.Id)
	server = &servertmp

	if server.StatusCode != 200 {
		return fmt.Errorf("Error fetching server: %s", server.Response)
	}
	logicalID := instance.LogicalID(server.Entities.Nics.Items[0].Properties.Ips[0])
	description := instance.Description{
		Tags:      spec.Tags,
		ID:        id,
		LogicalID: &logicalID,
	}
	towrite, err := json.MarshalIndent(description, "", "\t")

	if err != nil {
		return fmt.Errorf("Error occurred while marshalling data into json %s", err)
	}

	filePath := path.Join(p.dir, server.Id+filesufix)
	err = ioutil.WriteFile(filePath, towrite, 0644)
	if err != nil {
		return fmt.Errorf("An error occurred while trying to write to file %s, %s", filePath, err.Error())
	}
	return nil
}

func (p *plugin) removeFile(file string) error {
	return os.Remove(path.Join(p.dir, file))
}

func (p *plugin) getExistingInstances(tags map[string]string) (descriptions []instance.Description, err error) {

	files, err := ioutil.ReadDir(p.dir)

	if err != nil {
		log.Printf("Error occurred while reading folder %s %s", p.dir, err.Error())
		return nil, err
	}
	profitbricks.SetAuth(p.username, p.password)
	for _, file := range files {
		if strings.Contains(file.Name(), filesufix) {
			substring := file.Name()[0 : len(file.Name())-len(filesufix)]
			b, err := ioutil.ReadFile(file.Name())
			if err != nil {
				return nil, fmt.Errorf("Error occurred while reading conent of '%s' %s", file.Name(), err)
			}
			var description instance.Description

			err = json.Unmarshal(b, &description)
			if err != nil {
				return nil, fmt.Errorf("Error occurred while parsing conent of '%s' %s", file.Name(), err)
			}

			server := profitbricks.GetServer(description.Tags["datacenter_id"], substring)
			if server.StatusCode != 200 {
				log.Printf("Instance %s seems to be removed. Skipping.", file.Name(), server.StatusCode)
				err := p.removeFile(file.Name())
				if err != nil {
					return descriptions, err
				}
				continue
			}

			ip := instance.LogicalID(server.Entities.Nics.Items[0].Properties.Ips[0])
			name := file.Name()

			descriptions = append(descriptions, instance.Description{
				ID:        instance.ID(name[0 : len(name)-len(filesufix)]),
				LogicalID: &ip,
				Tags:      description.Tags,
			})
		}

	}

	return descriptions, err
}
