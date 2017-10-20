package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thirdwavellc/go-proxmox/proxmox"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROXMOX_HOST", nil),
				Description: "Proxmox host for authentication",
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROXMOX_USER", nil),
				Description: "Proxmox user for authentication",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROXMOX_PASSWORD", nil),
				Description: "Proxmox password for authentication",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"proxmox_container": resourceContainer(),
			"proxmox_group":     resourceGroup(),
			"proxmox_backup":    resourceBackup(),
			"proxmox_user":      resourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := proxmox.ProxmoxClient{
		Host:     d.Get("host").(string),
		User:     d.Get("user").(string),
		Password: d.Get("password").(string),
	}

	ticketReq := &proxmox.TicketRequest{
		Username: d.Get("user").(string),
		Password: d.Get("password").(string),
	}

	auth, err := config.GetAuth(ticketReq)
	config.Auth = auth

	return &config, err
}
