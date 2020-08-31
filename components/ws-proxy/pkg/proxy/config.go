// Copyright (c) 2020 TypeFox GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package proxy

import (
	"github.com/gitpod-io/gitpod/common-go/util"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/xerrors"
)

// Config is the configuration for a WorkspaceProxy
type Config struct {
	HTTPS struct {
		Enabled     bool   `json:"enabled"`
		Key         string `json:"key"`
		Certificate string `json:"crt"`
	} `json:"https,omitempty"`

	TransportConfig    *TransportConfig    `json:"transportConfig"`
	TheiaServer        *TheiaServer        `json:"theiaServer"`
	GitpodInstallation *GitpodInstallation `json:"gitpodInstallation"`
	WorkspacePodConfig *WorkspacePodConfig `json:"workspacePodConfig"`
}

// Validate validates the configuration to catch issues during startup and not at runtime
func (c *Config) Validate() error {
	err := c.TransportConfig.Validate()
	if err != nil {
		return err
	}
	err = c.TheiaServer.Validate()
	if err != nil {
		return err
	}
	err = c.GitpodInstallation.Validate()
	if err != nil {
		return err
	}
	err = c.WorkspacePodConfig.Validate()
	if err != nil {
		return err
	}
	return err
}

// WorkspacePodConfig contains config around the workspace pod
type WorkspacePodConfig struct {
	ServiceTemplate     string `json:"serviceTemplate"`
	PortServiceTemplate string `json:"portServiceTemplate"`
	TheiaPort           uint16 `json:"theiaPort"`
	SupervisorPort      uint16 `json:"supervisorPort"`
	StaticFrontendPort  uint16 `json:"staticFrontendPort"`
}

// Validate validates the configuration to catch issues during startup and not at runtime
func (c *WorkspacePodConfig) Validate() error {
	if c == nil {
		return xerrors.Errorf("WorkspacePodConfig not configured")
	}

	err := validation.ValidateStruct(c,
		validation.Field(&c.ServiceTemplate, validation.Required),
		validation.Field(&c.PortServiceTemplate, validation.Required),
		validation.Field(&c.TheiaPort, validation.Required),
		validation.Field(&c.SupervisorPort, validation.Required),
	)
	return err
}

// GitpodInstallation contains config regarding the Gitpod installation
type GitpodInstallation struct {
	Scheme              string `json:"scheme"`
	HostName            string `json:"hostName"`
	WorkspaceHostSuffix string `json:"workspaceHostSuffix"`
}

// Validate validates the configuration to catch issues during startup and not at runtime
func (c *GitpodInstallation) Validate() error {
	if c == nil {
		return xerrors.Errorf("GitpodInstallation not configured")
	}

	err := validation.ValidateStruct(c,
		validation.Field(&c.Scheme, validation.Required),
		validation.Field(&c.HostName, validation.Required),            // TODO IP ONLY: Check if there is any dependency. If yes, remove it.
		validation.Field(&c.WorkspaceHostSuffix, validation.Required), // TODO IP ONLY: Check if there is any dependency. If yes, remove it.
	)
	return err
}

// TheiaServer configures where to serve theia from
type TheiaServer struct {
	Scheme                  string `json:"scheme"`
	Host                    string `json:"host"`
	StaticVersionPathPrefix string `json:"staticVersionPathPrefix"`
}

// Validate validates the configuration to catch issues during startup and not at runtime
func (c *TheiaServer) Validate() error {
	if c == nil {
		return xerrors.Errorf("TheiaServer not configured")
	}

	err := validation.ValidateStruct(c,
		validation.Field(&c.Scheme, validation.Required),
		validation.Field(&c.Host, validation.Required),
		// StaticVersionPathPrefix might very well be ""
	)
	return err
}

// TransportConfig configures the way how ws-proxy connects to it's backend services
type TransportConfig struct {
	ConnectTimeout           util.Duration `json:"connectTimeout"`
	IdleConnTimeout          util.Duration `json:"idleConnTimeout"`
	WebsocketIdleConnTimeout util.Duration `json:"websocketIdleConnTimeout"`
	MaxIdleConns             int           `json:"maxIdleConns"`
}

// Validate validates the configuration to catch issues during startup and not at runtime
func (c *TransportConfig) Validate() error {
	if c == nil {
		return xerrors.Errorf("TransportConfig not configured")
	}

	err := validation.ValidateStruct(c,
		validation.Field(&c.ConnectTimeout, validation.Required),
		validation.Field(&c.IdleConnTimeout, validation.Required),
		validation.Field(&c.WebsocketIdleConnTimeout, validation.Required),
		validation.Field(&c.MaxIdleConns, validation.Required, validation.Min(1)),
	)
	return err
}
