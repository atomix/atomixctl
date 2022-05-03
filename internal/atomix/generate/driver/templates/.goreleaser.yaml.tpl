project_name: {{ .Driver.Name | toSnake }}
before:
  hooks:
    - go mod tidy
builds:
  - id: plugin
    hooks:
      pre:
        - atomix build plugin --name {{ .Driver.Name }} --version={{ "{{ .Version }}" }} --target {{ .Runtime.Version }} --output build/{{ .Driver.Name }}-{{ "{{ .Version }}" }}.{{ .Runtime.Version }}.so ./cmd/{{ .Driver.Name | toKebab }}
    builder: prebuilt
    goos:
      - linux
    goarch:
      - amd64
    prebuilt:
      path: build/{{ .Driver.Name }}-{{ "{{ .Version }}" }}.{{ .Runtime.Version }}.so
snapshot:
  name_template: "{{ "{{ incpatch .Version }}" }}-dev"
changelog:
  sort: asc
