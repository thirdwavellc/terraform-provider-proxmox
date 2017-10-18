package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thirdwavellc/go-proxmox/proxmox"
	"strings"
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
			"ssh_keys": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceContainerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)
	req := &proxmox.NewContainerRequest{
		Node:         d.Get("node").(string),
		VMID:         d.Get("vmid").(string),
		OsTemplate:   d.Get("os_template").(string),
		Net0:         d.Get("net0").(string),
		Storage:      d.Get("storage").(string),
		RootFs:       d.Get("root_fs").(string),
		Cores:        d.Get("cores").(int),
		Memory:       d.Get("memory").(int),
		Swap:         d.Get("swap").(int),
		Hostname:     d.Get("hostname").(string),
		Password:     d.Get("root_password").(string),
		OnBoot:       d.Get("on_boot").(int),
		Unprivileged: d.Get("unprivileged").(int),
	}
	if ssh_keys_len := d.Get("ssh_keys.#").(int); ssh_keys_len > 0 {
		req.SshPublicKeys = formatSshKeys(d, ssh_keys_len)
	}
	createUpid, err := client.CreateContainer(req)

	createStatusReq := &proxmox.NodeTaskStatusRequest{
		Node: req.Node,
		UPID: createUpid,
	}
	createTask, err := client.CheckNodeTaskStatus(createStatusReq)

	if err != nil {
		return err
	}

	if createTask.ExitStatus != "OK" {
		return err
	}

	startReq := &proxmox.ExistingContainerRequest{
		Node: req.Node,
		VMID: req.VMID,
	}
	startUpid, err := client.StartContainer(startReq)

	if err != nil {
		return err
	}

	startStatusReq := &proxmox.NodeTaskStatusRequest{
		Node: req.Node,
		UPID: startUpid,
	}
	startTask, err := client.CheckNodeTaskStatus(startStatusReq)

	if startTask.ExitStatus != "OK" {
		return errors.New("Exit Status: " + startTask.ExitStatus)
	}

	d.SetId(req.VMID)

	return nil
}

func resourceContainerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)
	req := &proxmox.ContainerConfigRequest{
		Node: d.Get("node").(string),
		VMID: d.Get("vmid").(string),
	}
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
		return errors.New("You cannot change the vmid of an already created machine!")
	}

	client := m.(*proxmox.ProxmoxClient)
	req := &proxmox.ExistingContainerRequest{
		Node: d.Get("node").(string),
		VMID: d.Get("vmid").(string),
	}
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
		// TODO: handle this with separate resizing call?
		return errors.New("You cannot change the root_fs of an already created machine!")
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
	if d.HasChange("ssh_keys") {
		return errors.New("You cannot change ssh_keys of an already created machine!")
	}
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
	node := d.Get("node").(string)
	vmid := d.Get("vmid").(string)

	containerStatusReq := &proxmox.ContainerStatusRequest{
		Node: node,
		VMID: vmid,
	}
	containerStatus, err := client.GetContainerStatus(containerStatusReq)

	if err != nil {
		return err
	}

	containerRunning := containerStatus.Status == "running"

	if containerRunning {
		shutdownReq := &proxmox.ExistingContainerRequest{
			Node: node,
			VMID: vmid,
		}
		shutdownUpid, err := client.ShutdownContainer(shutdownReq)

		if err != nil {
			return err
		}

		shutdownStatusReq := &proxmox.NodeTaskStatusRequest{
			Node: node,
			UPID: shutdownUpid,
		}
		shutdownTask, err := client.CheckNodeTaskStatus(shutdownStatusReq)

		if shutdownTask.ExitStatus != "OK" {
			return err
		}
	}

	req := &proxmox.ExistingContainerRequest{
		Node: d.Get("node").(string),
		VMID: d.Get("vmid").(string),
	}
	deleteUpid, err := client.DeleteContainer(req)

	if err != nil {
		return err
	}

	deleteStatusRequest := &proxmox.NodeTaskStatusRequest{
		Node: d.Get("node").(string),
		UPID: deleteUpid,
	}
	deleteTask, err := client.CheckNodeTaskStatus(deleteStatusRequest)

	if err != nil {
		return err
	}

	if deleteTask.ExitStatus != "OK" {
		return errors.New("Exit Status: " + deleteTask.ExitStatus)
	}

	return nil
}

func formatSshKeys(d *schema.ResourceData, ssh_keys_len int) string {
	keys := make([]string, ssh_keys_len)
	for i := 0; i < ssh_keys_len; i++ {
		keys = append(keys, d.Get(fmt.Sprintf("ssh_keys.%d", i)).(string))
	}
	return strings.Join(keys, "\n")
}
