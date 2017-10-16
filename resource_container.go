package main

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thirdwavellc/go-proxmox/proxmox"
)

func resourceContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerCreate,
		Read:   resourceContainerRead,
		Update: resourceContainerUpdate,
		Delete: resourceContainerDelete,
		Schema: map[string]*schema.Schema{
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vmid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"os_template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"root_fs": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceContainerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)
	req := &proxmox.NewContainerRequest{}
	req.Node = d.Get("node").(string)
	req.VMID = d.Get("vmid").(string)
	req.OsTemplate = d.Get("os_template").(string)
	//req.Net0 = d.Get("net0").(string)
	//req.Storage = d.Get("storage").(string)
	req.RootFs = d.Get("root_fs").(string)
	//req.Cores = d.Get("cores").(int)
	//req.Memory = d.Get("memory").(int)
	//req.Swap = d.Get("swap").(int)
	//req.Hostname = d.Get("hostname").(string)
	//req.OnBoot = d.Get("on_boot").(bool)
	//req.Password = d.Get("root_password").(string)
	//req.SshPublicKeys = d.Get("ssh_public_keys").(string)
	//req.Unprivileged = d.Get("unprivileged").(bool)
	upid, err := client.CreateContainer(req)

	statusRequest := proxmox.NodeTaskStatusRequest{}
	statusRequest.Node = req.Node
	statusRequest.UPID = upid
	task, err := client.CheckNodeTaskStatus(statusRequest)

	if err != nil {
		return err
	}

	if task.ExitStatus == "OK" {
		d.SetId(req.VMID)
	}

	return nil
}

func resourceContainerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceContainerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceContainerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)
	req := &proxmox.ExistingContainerRequest{}
	req.Node = d.Get("node").(string)
	req.VMID = d.Get("vmid").(string)
	upid, err := client.DeleteContainer(req)

	if err != nil {
		return err
	}

	statusRequest := proxmox.NodeTaskStatusRequest{}
	statusRequest.Node = d.Get("node").(string)
	statusRequest.UPID = upid
	task, err := client.CheckNodeTaskStatus(statusRequest)

	if err != nil {
		return err
	}

	if task.ExitStatus != "OK" {
		return errors.New("Exit Status: " + task.ExitStatus)
	}

	return nil
}
