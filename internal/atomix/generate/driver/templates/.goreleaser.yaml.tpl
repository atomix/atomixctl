project_name: {{ .Driver.Name | toSnake }}

before:
  hooks:
    - go mod tidy
    - go run github.com/atomix/cli/cmd/atomix-gen-deps@v0.2.1 --version {{ .Runtime.Version }} .
    - go mod tidy

builds:
  - id: plugin
    main: ./cmd/{{ .Driver.Name | toKebab }}
    binary: {{ .Driver.Name | toKebab }}-{{ "{{ .Version }}" }}.{{ .Runtime.Version }}.so
    goos:
      - linux
    goarch:
      - amd64
    env:
      - CC=gcc
      - CXX=g++
    flags:
      - -buildmode=plugin
      - -mod=readonly
      - -trimpath
    gcflags:
      - all=-N -l
    ldflags:
      - -s -w -X ./cmd/{{ .Driver.Name | toKebab }}.version={{ "{{ .Version }}" }} -X ./cmd/{{ .Driver.Name | toKebab }}.commit={{ "{{ .Commit }}" }}

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ "{{ incpatch .Version }}" }}-dev"

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'

{{- if .Repo.Name }}
release:
  github:
    owner: {{ .Repo.Owner }}
    name: {{ .Repo.Name }}
  prerelease: auto
  draft: true
{{- end }}
