module {{ .Module.Path }}

go 1.18

require (
    github.com/atomix/runtime-api {{ .Runtime.Version }}
)