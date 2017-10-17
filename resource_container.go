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
			"net0": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cores": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"swap": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"root_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"on_boot": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"unprivileged": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
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
	req.Net0 = d.Get("net0").(string)
	req.Storage = d.Get("storage").(string)
	req.RootFs = d.Get("root_fs").(string)
	req.Cores = d.Get("cores").(int)
	req.Memory = d.Get("memory").(int)
	req.Swap = d.Get("swap").(int)
	req.Hostname = d.Get("hostname").(string)
	req.Password = d.Get("root_password").(string)
	req.OnBoot = d.Get("on_boot").(int)
	req.Unprivileged = d.Get("unprivileged").(int)
	//req.SshPublicKeys = d.Get("ssh_public_keys").(string)
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
	client := m.(*proxmox.ProxmoxClient)
	req := &proxmox.ExistingContainerRequest{}
	req.Node = d.Get("node").(string)
	req.VMID = d.Get("vmid").(string)
	container, err := client.GetContainerConfig(req)

	if err != nil {
		return err
	}

	d.Set("hostname", container.Hostname)
	d.Set("cores", container.Cores)
	d.Set("memory", container.Memory)
	d.Set("swap", container.Swap)
	d.Set("net0", container.Net0)

	return nil
}

func resourceContainerUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	if d.HasChange("vmid") {
		return errors.New("You cannot change the VMID of an already created machine")
	}

	client := m.(*proxmox.ProxmoxClient)
	req := &proxmox.ExistingContainerRequest{}
	req.Node = d.Get("node").(string)
	req.VMID = d.Get("vmid").(string)
	if d.HasChange("os_template") {
		req.OsTemplate = d.Get("os_template").(string)
	}
	if d.HasChange("net0") {
		req.Net0 = d.Get("net0").(string)
	}
	if d.HasChange("storage") {
		req.Storage = d.Get("storage").(string)
	}
	if d.HasChange("root_fs") {
		req.RootFs = d.Get("root_fs").(string)
	}
	if d.HasChange("cores") {
		req.Cores = d.Get("cores").(int)
	}
	if d.HasChange("memory") {
		req.Memory = d.Get("memory").(int)
	}
	if d.HasChange("swap") {
		req.Swap = d.Get("swap").(int)
	}
	if d.HasChange("hostname") {
		req.Hostname = d.Get("hostname").(string)
	}
	if d.HasChange("root_password") {
		req.Password = d.Get("root_password").(string)
	}
	if d.HasChange("on_boot") {
		req.OnBoot = d.Get("on_boot").(int)
	}
	if d.HasChange("unprivileged") {
		req.Unprivileged = d.Get("unprivileged").(int)
	}
	//req.SshPublicKeys = d.Get("ssh_public_keys").(string)
	_, err := client.UpdateContainer(req)

	if err != nil {
		return err
	}

	d.SetPartial("os_template")
	d.SetPartial("net0")
	d.SetPartial("storage")
	d.SetPartial("root_fs")
	d.SetPartial("cores")
	d.SetPartial("memory")
	d.SetPartial("swap")
	d.SetPartial("hostname")
	d.SetPartial("root_password")
	d.SetPartial("on_boot")
	d.SetPartial("unprivileged")

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
