// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package proto

import (
	"fmt"
	runtimev1 "github.com/atomix/api/pkg/atomix/runtime/v1"
	"github.com/atomix/cli/internal/generator/proto/meta"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

const moduleName = "atomix"

// NewModule creates a new proto module
func NewModule(name string, template string) pgs.Module {
	return &Module{
		ModuleBase: &pgs.ModuleBase{},
		file:       name,
		template:   template,
	}
}

// Module is the code generation module
type Module struct {
	*pgs.ModuleBase
	ctx      *meta.Context
	file     string
	template string
}

// Name returns the module name
func (m *Module) Name() string {
	return moduleName
}

// InitContext initializes the module context
func (m *Module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = meta.NewContext(pgsgo.InitContext(c.Parameters()))
}

// Execute executes the code generator
func (m *Module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, target := range targets {
		m.executeTarget(target)
	}
	return m.Artifacts()
}

func (m *Module) executeTarget(target pgs.File) {
	println(target.File().InputPath())
	for _, service := range target.Services() {
		m.executeService(service)
	}
}

// executeService generates a store from a Protobuf service
//nolint:gocyclo
func (m *Module) executeService(service pgs.Service) {
	primitiveType, err := meta.GetPrimitiveType(service)
	if err != nil {
		return
	}

	importsSet := make(map[string]meta.PackageMeta)
	var addImport = func(t *meta.TypeMeta) {
		if t.Package.Import {
			baseAlias := t.Package.Alias
			i := 0
			for {
				importPackage, ok := importsSet[t.Package.Alias]
				if ok {
					if importPackage.Path != t.Package.Path {
						t.Package.Alias = fmt.Sprintf("%s%d", baseAlias, i)
					} else {
						break
					}
				} else {
					importsSet[t.Package.Alias] = t.Package
				}
				i++
			}
		}
	}

	// Iterate through the methods on the service and construct method metadata for the template.
	methods := make([]meta.MethodMeta, 0)
	for _, method := range service.Methods() {
		operationID, err := meta.GetOperationID(method)
		if err != nil {
			panic(err)
		}

		// Get the operation type for the method.
		operationType, err := meta.GetOperationType(method)
		if err != nil {
			panic(err)
		}

		async, err := meta.GetAsync(method)
		if err != nil {
			panic(err)
		}

		methodTypeMeta := meta.MethodTypeMeta{
			IsCommand: operationType == runtimev1.OperationType_COMMAND,
			IsQuery:   operationType == runtimev1.OperationType_QUERY,
			IsSync:    !async,
			IsAsync:   async,
		}

		requestHeaders, err := m.ctx.GetHeadersFieldMeta(method.Input())
		if err != nil {
			panic(err)
		} else if requestHeaders == nil {
			panic("no request headers found on method input " + method.Input().FullyQualifiedName())
		}

		requestMeta := meta.RequestMeta{
			MessageMeta: meta.MessageMeta{
				Type: m.ctx.GetMessageTypeMeta(method.Input()),
			},
			Headers:  *requestHeaders,
			IsUnary:  !method.ClientStreaming(),
			IsStream: method.ClientStreaming(),
		}
		addImport(&requestMeta.Type)

		responseHeaders, err := m.ctx.GetHeadersFieldMeta(method.Output())
		if err != nil {
			panic(err)
		} else if responseHeaders == nil {
			panic("no request headers found on method input " + method.Output().FullyQualifiedName())
		}

		// Generate output metadata from the output type.
		responseMeta := meta.ResponseMeta{
			MessageMeta: meta.MessageMeta{
				Type: m.ctx.GetMessageTypeMeta(method.Output()),
			},
			Headers:  *responseHeaders,
			IsUnary:  !method.ServerStreaming(),
			IsStream: method.ServerStreaming(),
		}
		addImport(&responseMeta.Type)

		methodMeta := meta.MethodMeta{
			ID:       operationID,
			Name:     method.Name().UpperCamelCase().String(),
			Comment:  method.SourceCodeInfo().LeadingComments(),
			Type:     methodTypeMeta,
			Request:  requestMeta,
			Response: responseMeta,
		}

		methods = append(methods, methodMeta)
	}

	// Generate a list of imports from the deduplicated package metadata set.
	imports := make([]meta.PackageMeta, 0, len(importsSet))
	for _, importPkg := range importsSet {
		imports = append(imports, importPkg)
	}

	atomMeta := meta.AtomMeta{
		Name: primitiveType,
		ServiceMeta: meta.ServiceMeta{
			Type: meta.ServiceTypeMeta{
				Name:    pgsgo.PGGUpperCamelCase(service.Name()).String(),
				Package: m.ctx.GetPackageMeta(service),
			},
			Comment: service.SourceCodeInfo().LeadingComments(),
			Methods: methods,
		},
	}

	// Generate the store metadata.
	meta := meta.CodegenMeta{
		Generator: meta.GeneratorMeta{
			Prefix: m.BuildContext.Parameters().Str("prefix"),
		},
		Location: meta.LocationMeta{},
		Package:  m.ctx.GetPackageMeta(service),
		Imports:  imports,
		Atom:     atomMeta,
	}

	m.OverwriteGeneratorTemplateFile(m.ctx.GetFilePath(service, m.file), NewTemplate(m.ctx.GetTemplatePath(m.template), importsSet), meta)
}
