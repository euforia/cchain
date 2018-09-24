package compose

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/cli/cli/compose/loader"
	"github.com/docker/cli/cli/compose/types"
	"github.com/pkg/errors"
)

func parseDockerComposeFile(env map[string]string, workingDir, filename string) (*types.Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parsedComposeFile, err := loader.ParseYAML(b)
	if err != nil {
		return nil, err
	}

	// Config file
	configFile := types.ConfigFile{
		Filename: filename,
		Config:   parsedComposeFile,
	}

	// Config details
	configDetails := types.ConfigDetails{
		WorkingDir:  workingDir,
		ConfigFiles: []types.ConfigFile{configFile},
		Environment: env,
	}

	// Actual config
	// We load it in order to retrieve the parsed output configuration!
	// This will output a github.com/docker/cli ServiceConfig
	// Which is similar to our version of ServiceConfig
	return loader.Load(configDetails)
}

// converts os.Environ() ([]string) to map[string]string
// based on https://github.com/docker/cli/blob/5dd30732a23bbf14db1c64d084ae4a375f592cfa/cli/command/stack/deploy_composefile.go#L143
func buildEnvironment() (map[string]string, error) {
	env := os.Environ()
	result := make(map[string]string, len(env))
	for _, s := range env {
		// if value is empty, s is like "K=", not "K".
		if !strings.Contains(s, "=") {
			return result, errors.Errorf("unexpected environment %q", s)
		}
		kv := strings.SplitN(s, "=", 2)
		result[kv[0]] = kv[1]
	}
	return result, nil
}

func mergeComposeObject(oldCompose *types.Config, newCompose *types.Config) (*types.Config, error) {
	if oldCompose == nil || newCompose == nil {
		return nil, fmt.Errorf("Merge multiple compose error, compose config is nil")
	}
	oldComposeServiceNameMap := make(map[string]int, len(oldCompose.Services))
	for index, service := range oldCompose.Services {
		oldComposeServiceNameMap[service.Name] = index
	}

	for _, service := range newCompose.Services {
		index := 0
		if tmpIndex, ok := oldComposeServiceNameMap[service.Name]; !ok {
			oldCompose.Services = append(oldCompose.Services, service)
			continue
		} else {
			index = tmpIndex
		}
		tmpOldService := oldCompose.Services[index]
		if service.Build.Dockerfile != "" {
			tmpOldService.Build = service.Build
		}
		if len(service.CapAdd) != 0 {
			tmpOldService.CapAdd = service.CapAdd
		}
		if len(service.CapDrop) != 0 {
			tmpOldService.CapDrop = service.CapDrop
		}
		if service.CgroupParent != "" {
			tmpOldService.CgroupParent = service.CgroupParent
		}
		if len(service.Command) != 0 {
			tmpOldService.Command = service.Command
		}
		if len(service.Configs) != 0 {
			tmpOldService.Configs = service.Configs
		}
		if service.ContainerName != "" {
			tmpOldService.ContainerName = service.ContainerName
		}
		if service.CredentialSpec.File != "" || service.CredentialSpec.Registry != "" {
			tmpOldService.CredentialSpec = service.CredentialSpec
		}
		if len(service.DependsOn) != 0 {
			tmpOldService.DependsOn = service.DependsOn
		}
		if service.Deploy.Mode != "" {
			tmpOldService.Deploy = service.Deploy
		}
		if len(service.Devices) != 0 {
			tmpOldService.Devices = service.Devices
		}
		if len(service.DNS) != 0 {
			tmpOldService.DNS = service.DNS
		}
		if len(service.DNSSearch) != 0 {
			tmpOldService.DNSSearch = service.DNSSearch
		}
		if service.DomainName != "" {
			tmpOldService.DomainName = service.DomainName
		}
		if len(service.Entrypoint) != 0 {
			tmpOldService.Entrypoint = service.Entrypoint
		}
		if len(service.Environment) != 0 {
			tmpOldService.Environment = service.Environment
		}
		if len(service.EnvFile) != 0 {
			tmpOldService.EnvFile = service.EnvFile
		}
		if len(service.Expose) != 0 {
			tmpOldService.Expose = service.Expose
		}
		if len(service.ExternalLinks) != 0 {
			tmpOldService.ExternalLinks = service.ExternalLinks
		}
		if len(service.ExtraHosts) != 0 {
			tmpOldService.ExtraHosts = service.ExtraHosts
		}
		if service.Hostname != "" {
			tmpOldService.Hostname = service.Hostname
		}
		if service.HealthCheck != nil {
			tmpOldService.HealthCheck = service.HealthCheck
		}
		if service.Image != "" {
			tmpOldService.Image = service.Image
		}
		if service.Ipc != "" {
			tmpOldService.Ipc = service.Ipc
		}
		if len(service.Labels) != 0 {
			tmpOldService.Labels = service.Labels
		}
		if len(service.Links) != 0 {
			tmpOldService.Links = service.Links
		}
		if service.Logging != nil {
			tmpOldService.Logging = service.Logging
		}
		if service.MacAddress != "" {
			tmpOldService.MacAddress = service.MacAddress
		}
		if service.NetworkMode != "" {
			tmpOldService.NetworkMode = service.NetworkMode
		}
		if len(service.Networks) != 0 {
			tmpOldService.Networks = service.Networks
		}
		if service.Pid != "" {
			tmpOldService.Pid = service.Pid
		}
		if len(service.Ports) != 0 {
			tmpOldService.Ports = service.Ports
		}
		if service.Privileged != tmpOldService.Privileged {
			tmpOldService.Privileged = service.Privileged
		}
		if service.ReadOnly != tmpOldService.ReadOnly {
			tmpOldService.ReadOnly = service.ReadOnly
		}
		if service.Restart != "" {
			tmpOldService.Restart = service.Restart
		}
		if len(service.Secrets) != 0 {
			tmpOldService.Secrets = service.Secrets
		}
		if len(service.SecurityOpt) != 0 {
			tmpOldService.SecurityOpt = service.SecurityOpt
		}
		if service.StdinOpen != tmpOldService.StdinOpen {
			tmpOldService.StdinOpen = service.StdinOpen
		}
		if service.StopGracePeriod != nil {
			tmpOldService.StopGracePeriod = service.StopGracePeriod
		}
		if service.StopSignal != "" {
			tmpOldService.StopSignal = service.StopSignal
		}
		if len(service.Tmpfs) != 0 {
			tmpOldService.Tmpfs = service.Tmpfs
		}
		if service.Tty != tmpOldService.Tty {
			tmpOldService.Tty = service.Tty
		}
		if len(service.Ulimits) != 0 {
			tmpOldService.Ulimits = service.Ulimits
		}
		if service.User != "" {
			tmpOldService.User = service.User
		}
		if len(service.Volumes) != 0 {
			tmpOldService.Volumes = service.Volumes
		}
		if service.WorkingDir != "" {
			tmpOldService.WorkingDir = service.WorkingDir
		}
		oldCompose.Services[index] = tmpOldService
	}

	return oldCompose, nil
}
