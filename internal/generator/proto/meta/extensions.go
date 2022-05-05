// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package meta

import (
	runtimev1 "github.com/atomix/api/pkg/atomix/runtime/v1"
	"github.com/atomix/sdk/pkg/errors"
	"github.com/gogo/protobuf/gogoproto"
	gogoprotobuf "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
)

// GetPrimitiveType gets the name extension from the given service
func GetPrimitiveType(service pgs.Service) (string, error) {
	var primitiveType string
	ok, err := service.Extension(getExtensionDesc(runtimev1.E_Atom), &primitiveType)
	if err != nil {
		return "", err
	} else if !ok {
		return "", errors.NewInvalid("extension '%s' is not set", runtimev1.E_Atom.Name)
	}
	return primitiveType, nil
}

// GetAsync gets the async extension from the given service
func GetAsync(method pgs.Method) (bool, error) {
	var async bool
	ok, err := method.Extension(getExtensionDesc(runtimev1.E_Atom), &async)
	if err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}
	return async, nil
}

// GetOperationName gets the name extension from the given method
func GetOperationName(method pgs.Method) (string, error) {
	var opName string
	ok, err := method.Extension(getExtensionDesc(runtimev1.E_OperationName), &opName)
	if err != nil {
		return "", err
	} else if !ok {
		return method.Name().String(), nil
	}
	return opName, nil
}

// GetOperationID gets the id extension from the given method
func GetOperationID(method pgs.Method) (uint32, error) {
	var operationID uint32
	ok, err := method.Extension(getExtensionDesc(runtimev1.E_OperationId), &operationID)
	if err != nil {
		return 0, err
	} else if !ok {
		return 0, errors.NewInvalid("extension '%s' is not set", runtimev1.E_OperationId.Name)
	}
	return operationID, nil
}

// GetOperationType gets the optype extension from the given method
func GetOperationType(method pgs.Method) (runtimev1.OperationType, error) {
	var operationType runtimev1.OperationType
	ok, err := method.Extension(getExtensionDesc(runtimev1.E_OperationType), &operationType)
	if err != nil {
		return 0, err
	} else if !ok {
		return 0, errors.NewInvalid("extension '%s' is not set", runtimev1.E_OperationType.Name)
	}
	return operationType, nil
}

// GetHeaders gets the headers extension from the given field
func GetHeaders(field pgs.Field) (bool, error) {
	var headers bool
	ok, err := field.Extension(getExtensionDesc(runtimev1.E_Headers), &headers)
	if err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}
	return headers, nil
}

func getExtensionDesc(extension *gogoprotobuf.ExtensionDesc) *proto.ExtensionDesc {
	return &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: extension.ExtensionType,
		Field:         extension.Field,
		Name:          extension.Name,
		Tag:           extension.Tag,
		Filename:      extension.Filename,
	}
}

// GetEmbed gets the embed extension from the given field
func GetEmbed(field pgs.Field) (*bool, error) {
	var embed bool
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Embed), &embed)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &embed, nil
}

// GetCastType gets the casttype extension from the given field
func GetCastType(field pgs.Field) (*string, error) {
	var castType string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Casttype), &castType)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &castType, nil
}

// GetCastKey gets the castkey extension from the given field
func GetCastKey(field pgs.Field) (*string, error) {
	var castKey string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Castkey), &castKey)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &castKey, nil
}

// GetCastValue gets the castvalue extension from the given field
func GetCastValue(field pgs.Field) (*string, error) {
	var castValue string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Castvalue), &castValue)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &castValue, nil
}

// GetCustomName gets the customname extension from the given field
func GetCustomName(field pgs.Field) (*string, error) {
	var customName string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Customname), &customName)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &customName, nil
}

// GetCustomType gets the customtype extension from the given field
func GetCustomType(field pgs.Field) (*string, error) {
	var customType string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Customtype), &customType)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &customType, nil
}

// GetNullable gets the nullable extension from the given field
func GetNullable(field pgs.Field) (*bool, error) {
	var nullable bool
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Nullable), &nullable)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &nullable, nil
}
