package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thirdwavellc/go-proxmox/proxmox"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"expire": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keys": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)
	userId := d.Get("user_id").(string)

	createReq := &proxmox.NewUserRequest{
		UserId:    userId,
		Comment:   d.Get("comment").(string),
		Email:     d.Get("email").(string),
		Enable:    d.Get("enable").(int),
		Expire:    d.Get("expire").(int),
		FirstName: d.Get("first_name").(string),
		Keys:      d.Get("keys").(string),
		LastName:  d.Get("last_name").(string),
		Password:  d.Get("password").(string),
	}

	// TODO: handle groups

	_, err := client.CreateUser(createReq)

	if err != nil {
		return err
	}

	d.SetId(userId)

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)

	configReq := &proxmox.UserConfigRequest{
		UserId: d.Get("user_id").(string),
	}
	userConfig, err := client.GetUserConfig(configReq)

	if err != nil {
		return err
	}

	d.Set("comment", userConfig.Comment)
	d.Set("email", userConfig.Email)
	d.Set("enable", userConfig.Enable)
	d.Set("expire", userConfig.Expire)
	d.Set("first_name", userConfig.FirstName)
	// TODO: handle groups
	d.Set("keys", userConfig.Keys)
	d.Set("last_name", userConfig.LastName)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	client := m.(*proxmox.ProxmoxClient)

	updateReq := &proxmox.ExistingUserRequest{
		UserId: d.Get("user_id").(string),
	}
	if d.HasChange("comment") {
		updateReq.Comment = d.Get("comment").(string)
	}
	if d.HasChange("email") {
		updateReq.Email = d.Get("email").(string)
	}
	if d.HasChange("enable") {
		updateReq.Enable = d.Get("enable").(int)
	}
	if d.HasChange("expire") {
		updateReq.Expire = d.Get("expire").(int)
	}
	if d.HasChange("first_name") {
		updateReq.FirstName = d.Get("first_name").(string)
	}
	if d.HasChange("keys") {
		updateReq.Keys = d.Get("keys").(string)
	}
	if d.HasChange("last_name") {
		updateReq.LastName = d.Get("last_name").(string)
	}
	_, err := client.UpdateUser(updateReq)

	if err != nil {
		return err
	}

	d.SetPartial("comment")
	d.SetPartial("email")
	d.SetPartial("enable")
	d.SetPartial("expire")
	d.SetPartial("first_name")
	d.SetPartial("keys")
	d.SetPartial("last_name")

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)

	deleteReq := &proxmox.ExistingUserRequest{
		UserId: d.Get("user_id").(string),
	}
	_, err := client.DeleteUser(deleteReq)

	if err != nil {
		return err
	}

	return nil
}
