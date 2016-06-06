package main

import (
	"fmt"
	"os"
	"github.com/dghubble/sling"
	"net/http"
	"github.com/mscrypto/mscryptotest/data_types"
	"errors"
)

type Client struct {
	LoadBalancerService *LoadBalancerService
	VirtualGuestService *VirtualGuestService
}

func NewClient(username string, apiKey string, service string) (*Client, error) {

	if service == "LoadBalancer" {
		return &Client{
			LoadBalancerService: NewLoadBalancerService(username, apiKey),
		}, nil
	}

	if service == "VirtualServer" {
		return &Client{
			VirtualGuestService: NewVirtualGuestService(username, apiKey),
		}, nil
	}

	return nil, errors.New(service + "is not available.")
}

const baseURL = "api.softlayer.com/rest/v3/"

type LoadBalancerService struct {
	sling *sling.Sling
}

func NewLoadBalancerService(username string, apiKey string) *LoadBalancerService {
	return &LoadBalancerService{
		sling: sling.New().Client(nil).Base("https://"+username+":"+apiKey+"@"+baseURL),
	}
}

func (s *LoadBalancerService) getIpAddress() (*data_types.ResLoadBalancer, *http.Response, error) {
	resLoadBalancer := new(data_types.ResLoadBalancer)
	resLoadBalancerError := new(data_types.ResLoadBalancer)
	resp, err := s.sling.New().Get("SoftLayer_Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress/146505/getIpAddress.json").Receive(resLoadBalancer,resLoadBalancerError)
	return resLoadBalancer, resp, err
}


type VirtualGuestService struct {
	sling *sling.Sling
}

func NewVirtualGuestService(username string, apiKey string) *VirtualGuestService {
	return &VirtualGuestService{
		sling: sling.New().Client(nil).Base("https://"+username+":"+apiKey+"@"+baseURL),
	}
}

func (s *VirtualGuestService) createObject(virtualGuestTemplate *data_types.SoftLayer_Virtual_Guest_Template) (*data_types.SoftLayer_Virtual_Guest, *http.Response, error) {
	resVirtualGuest := new(data_types.SoftLayer_Virtual_Guest)
	resVirtualGuestError := new(data_types.SoftLayer_Virtual_Guest)
	resp, err := s.sling.New().Post("SoftLayer_Virtual_Guest/generateOrderTemplate.json").BodyJSON(virtualGuestTemplate).Receive(resVirtualGuest, resVirtualGuestError)
	return resVirtualGuest, resp, err
}

func main() {
	username := os.Getenv("SL_USERNAME")
	apiKey := os.Getenv("SL_API_KEY")
	myClientLoadBalancer,_ := NewClient(username, apiKey, "LoadBalancer")
	loadBalancerService := myClientLoadBalancer.LoadBalancerService
	loadBalancer,_,_ := loadBalancerService.getIpAddress()
	fmt.Println("ID : ", loadBalancer.ID)
	fmt.Println("IPAddress : ", loadBalancer.IpAddress)

	myClientVirtualServer,_ := NewClient(username, apiKey, "VirtualServer")
	virtualGuestTemplate := data_types.SoftLayer_Virtual_Guest_Template{
		Hostname:  "testgo",
		Domain:    "ms.com",
		StartCpus: 1,
		MaxMemory: 1024,
		Datacenter: data_types.Datacenter{
			Name: "ams01",
		},
//		SshKeys:                      []SshKey{},  //or get the necessary keys and add here
		HourlyBillingFlag:            true,
		LocalDiskFlag:                true,
		OperatingSystemReferenceCode: "UBUNTU_LATEST",
	}
	virtualGuestService := myClientVirtualServer.VirtualGuestService
	virtualGuest,res,err := virtualGuestService.createObject(&virtualGuestTemplate)
	fmt.Println(res)
	fmt.Println(err)
	fmt.Println("-----Virtual Guest----")
	fmt.Println("HostnName : ", virtualGuest.Hostname)
}