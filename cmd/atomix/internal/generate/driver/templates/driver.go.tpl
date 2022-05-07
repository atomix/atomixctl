{{ $api := printf "%s%s" ( .Driver.Name | toSnake ) .Driver.APIVersion }}

package main

import (
    {{ $api }} "{{ .Module.Path }}/{{ .Driver.Name | toKebab }}/{{ .Driver.APIVersion }}"
	"github.com/atomix/sdk/pkg/config"
	"github.com/atomix/sdk/pkg/runtime/driver"
	"context"
)

var (
    version string
    commit string
)

var Codec = config.NewCodec[*{{ $api }}.Config](&{{ $api }}.Config{})

var Driver = driver.New[*{{ $api }}.Config](func(ctx context.Context, config *{{ $api }}.Config) (driver.Client, error) {
    return newMultiRaftClient(config), nil // TODO
})
