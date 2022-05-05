// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package meta

// CodegenMeta is the metadata for the code generator
type CodegenMeta struct {
	Generator GeneratorMeta
	Location  LocationMeta
	Package   PackageMeta
	Imports   []PackageMeta
	Atom      AtomMeta
}

// GeneratorMeta is the metadata for the code generator
type GeneratorMeta struct {
	Prefix string
}

// LocationMeta is the location of a code file
type LocationMeta struct {
	Path string
	File string
}

// PackageMeta is the package for a code file
type PackageMeta struct {
	Name   string
	Path   string
	Alias  string
	Import bool
}

// TypeMeta is the metadata for a store type
type TypeMeta struct {
	Name        string
	Package     PackageMeta
	IsPointer   bool
	IsScalar    bool
	IsCast      bool
	IsMessage   bool
	IsMap       bool
	IsRepeated  bool
	IsEnum      bool
	IsEnumValue bool
	IsBytes     bool
	IsString    bool
	IsInt32     bool
	IsInt64     bool
	IsUint32    bool
	IsUint64    bool
	IsFloat     bool
	IsDouble    bool
	IsBool      bool
	KeyType     *TypeMeta
	ValueType   *TypeMeta
	Values      []TypeMeta
}

// AtomMeta is the metadata for an atom
type AtomMeta struct {
	ServiceMeta
	Name string
}

// ServiceMeta is the metadata for a service
type ServiceMeta struct {
	Type    ServiceTypeMeta
	Comment string
	Methods []MethodMeta
}

// ServiceTypeMeta is metadata for a service type
type ServiceTypeMeta struct {
	Name    string
	Package PackageMeta
}

// FieldRefMeta is metadata for a field reference
type FieldRefMeta struct {
	Field FieldMeta
}

// FieldMeta is metadata for a field
type FieldMeta struct {
	Type TypeMeta
	Path []PathMeta
}

// PathMeta is metadata for a field path
type PathMeta struct {
	Name string
	Type TypeMeta
}

// MethodMeta is the metadata for a primitive method
type MethodMeta struct {
	ID       uint32
	Name     string
	Type     MethodTypeMeta
	Comment  string
	Request  RequestMeta
	Response ResponseMeta
}

// MessageMeta is the metadata for a message
type MessageMeta struct {
	Type TypeMeta
}

// RequestMeta is the type metadata for a message
type RequestMeta struct {
	MessageMeta
	Headers        FieldRefMeta
	PartitionKey   *FieldRefMeta
	PartitionRange *FieldRefMeta
	IsUnary        bool
	IsStream       bool
}

// ResponseMeta is the type metadata for a message
type ResponseMeta struct {
	MessageMeta
	Headers  FieldRefMeta
	IsUnary  bool
	IsStream bool
}

// MethodTypeMeta is the metadata for a store method type
type MethodTypeMeta struct {
	IsCommand bool
	IsQuery   bool
	IsSync    bool
	IsAsync   bool
}
