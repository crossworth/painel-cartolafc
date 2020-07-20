package vk

import (
	vkapi "github.com/himidori/golang-vk-api"
)

type VKClient struct {
	client *vkapi.VKClient
}

func NewVKClient(user string, password string) (*VKClient, error) {
	client, err := vkapi.NewVKClient(vkapi.DeviceIPhone, user, password)
	return &VKClient{
		client: client,
	}, err
}
