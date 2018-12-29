/**
 * Copyright 2018 Rafael Fernández López <ereslibre@ereslibre.es>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package server

import (
	"io"
	"net"
	"fmt"

	"github.com/spf13/cobra"

	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"github.com/ereslibre/cluster-apiserver/pkg/apiserver"
	"github.com/ereslibre/cluster-apiserver/pkg/apis/cluster/v1alpha1"
)

const defaultEtcdPathPrefix = "/registry/cluster.kubernetes.io"

type ServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
}

func NewServerOptions(out, errOut io.Writer) *ServerOptions {
	return &ServerOptions {
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			apiserver.Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion),
			genericoptions.NewProcessInfo("cluster-apiserver", "cluster"),
		),
	}
}

func NewCommandStartServer(serverOptions *ServerOptions, stopCh <-chan struct{}) *cobra.Command {
	o := *serverOptions
	cmd := &cobra.Command{
		Short: "Launch the API server",
		Long:  "Launch the API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)

	return cmd
}

func (o *ServerOptions) Complete() error {
	return nil
}

func (o *ServerOptions) Validate(args []string) error {
	return nil
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig, apiserver.Scheme); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
	}

	return config, nil
}

func (o *ServerOptions) RunServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
