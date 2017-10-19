package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thirdwavellc/go-proxmox/proxmox"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)
	groupId := d.Get("group_id").(string)

	createReq := &proxmox.NewGroupRequest{
		GroupId: groupId,
		Comment: d.Get("comment").(string),
	}
	_, err := client.CreateGroup(createReq)

	if err != nil {
		return err
	}

	d.SetId(groupId)

	return nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)

	configReq := &proxmox.GroupConfigRequest{
		GroupId: d.Get("group_id").(string),
	}
	groupConfig, err := client.GetGroupConfig(configReq)

	if err != nil {
		return err
	}

	d.Set("comment", groupConfig.Comment)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	client := m.(*proxmox.ProxmoxClient)

	updateReq := &proxmox.ExistingGroupRequest{
		GroupId: d.Get("group_id").(string),
	}
	if d.HasChange("comment") {
		updateReq.Comment = d.Get("comment").(string)
	}
	_, err := client.UpdateGroup(updateReq)

	if err != nil {
		return err
	}

	d.SetPartial("comment")

	return nil
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)

	deleteReq := &proxmox.ExistingGroupRequest{
		GroupId: d.Get("group_id").(string),
	}
	_, err := client.DeleteGroup(deleteReq)

	if err != nil {
		return err
	}

	return nil
}
