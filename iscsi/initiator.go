package iscsi

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yasker/go-iscsi-helper/util"
)

const (
	iscsiBinary = "iscsiadm"

	retryInterval = 1 * time.Second
	retryMax      = 5
)

func DiscoverTarget(ip, target string, ne *util.NamespaceExecutor) error {
	opts := []string{
		"-m", "discovery",
		"-t", "sendtargets",
		"-p", ip,
	}
	output, err := ne.Execute(iscsiBinary, opts)
	if err != nil {
		return err
	}
	if !strings.Contains(output, target) {
		return fmt.Errorf("Cannot find target %s in discovered targets %s", target, output)
	}
	return nil
}

func DeleteDiscoveredTarget(ip, target string, ne *util.NamespaceExecutor) error {
	opts := []string{
		"-m", "node",
		"-o", "delete",
		"-p", ip,
		"-T", target,
	}
	_, err := ne.Execute(iscsiBinary, opts)
	if err != nil {
		return err
	}
	return nil
}

func LoginTarget(ip, target string, ne *util.NamespaceExecutor) error {
	opts := []string{
		"-m", "node",
		"-T", target,
		"-p", ip,
		"--login",
	}
	_, err := ne.Execute(iscsiBinary, opts)
	if err != nil {
		return err
	}
	return nil
}

func LogoutTarget(ip, target string, ne *util.NamespaceExecutor) error {
	opts := []string{
		"-m", "node",
		"-T", target,
		"-p", ip,
		"--logout",
	}
	_, err := ne.Execute(iscsiBinary, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetDevice(ip, target string, lun int, ne *util.NamespaceExecutor) (string, error) {
	var err error

	dev := ""
	for i := 0; i < retryMax; i++ {
		path := fmt.Sprintf("/dev/disk/by-path/ip-%s:3260-iscsi-%s-lun-%s", ip, target, strconv.Itoa(lun))
		opts := []string{
			"-fnve",
			path,
		}
		dev, err = ne.Execute("readlink", opts)
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(dev), nil
}