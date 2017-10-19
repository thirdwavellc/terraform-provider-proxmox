package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/satori/go.uuid"
	"github.com/thirdwavellc/go-proxmox/proxmox"
)

func resourceBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackupCreate,
		Read:   resourceBackupRead,
		Update: resourceBackupUpdate,
		Delete: resourceBackupDelete,
		Schema: map[string]*schema.Schema{
			"start_time": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"all": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"compress": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mail_notification": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mail_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"node": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceBackupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*proxmox.ProxmoxClient)

	createReq := &proxmox.NewBackupRequest{
		StartTime:        d.Get("start_time").(string),
		All:              d.Get("all").(int),
		Compress:         d.Get("compress").(string),
		MailNotification: d.Get("mail_notification").(string),
		MailTo:           d.Get("mail_to").(string),
		Mode:             d.Get("mode").(string),
		Node:             d.Get("node").(string),
	}
	_, err := client.CreateBackup(createReq)

	if err != nil {
		return err
	}

	d.SetId(uuid.NewV4().String())

	return nil
}

func resourceBackupRead(d *schema.ResourceData, m interface{}) error {
	// Due to API limitations, we can't perform any actions after creation
	// This is because the API doesn't return the backup id.
	// TODO: work around?
	return nil
}

func resourceBackupUpdate(d *schema.ResourceData, m interface{}) error {
	// Due to API limitations, we can't perform any actions after creation
	// This is because the API doesn't return the backup id.
	// TODO: work around?
	return nil
}

func resourceBackupDelete(d *schema.ResourceData, m interface{}) error {
	// Due to API limitations, we can't perform any actions after creation
	// This is because the API doesn't return the backup id.
	// TODO: work around?
	return nil
}
