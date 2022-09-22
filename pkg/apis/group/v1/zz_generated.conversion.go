//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by conversion-gen. DO NOT EDIT.

package v1

import (
	unsafe "unsafe"

	group "github.com/solid-wang/covid/pkg/apis/group"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Demo)(nil), (*group.Demo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Demo_To_group_Demo(a.(*Demo), b.(*group.Demo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*group.Demo)(nil), (*Demo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_group_Demo_To_v1_Demo(a.(*group.Demo), b.(*Demo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DemoList)(nil), (*group.DemoList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DemoList_To_group_DemoList(a.(*DemoList), b.(*group.DemoList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*group.DemoList)(nil), (*DemoList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_group_DemoList_To_v1_DemoList(a.(*group.DemoList), b.(*DemoList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DemoSpec)(nil), (*group.DemoSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DemoSpec_To_group_DemoSpec(a.(*DemoSpec), b.(*group.DemoSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*group.DemoSpec)(nil), (*DemoSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_group_DemoSpec_To_v1_DemoSpec(a.(*group.DemoSpec), b.(*DemoSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DemoStatus)(nil), (*group.DemoStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DemoStatus_To_group_DemoStatus(a.(*DemoStatus), b.(*group.DemoStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*group.DemoStatus)(nil), (*DemoStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_group_DemoStatus_To_v1_DemoStatus(a.(*group.DemoStatus), b.(*DemoStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_Demo_To_group_Demo(in *Demo, out *group.Demo, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_DemoSpec_To_group_DemoSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1_DemoStatus_To_group_DemoStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1_Demo_To_group_Demo is an autogenerated conversion function.
func Convert_v1_Demo_To_group_Demo(in *Demo, out *group.Demo, s conversion.Scope) error {
	return autoConvert_v1_Demo_To_group_Demo(in, out, s)
}

func autoConvert_group_Demo_To_v1_Demo(in *group.Demo, out *Demo, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_group_DemoSpec_To_v1_DemoSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_group_DemoStatus_To_v1_DemoStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_group_Demo_To_v1_Demo is an autogenerated conversion function.
func Convert_group_Demo_To_v1_Demo(in *group.Demo, out *Demo, s conversion.Scope) error {
	return autoConvert_group_Demo_To_v1_Demo(in, out, s)
}

func autoConvert_v1_DemoList_To_group_DemoList(in *DemoList, out *group.DemoList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]group.Demo)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1_DemoList_To_group_DemoList is an autogenerated conversion function.
func Convert_v1_DemoList_To_group_DemoList(in *DemoList, out *group.DemoList, s conversion.Scope) error {
	return autoConvert_v1_DemoList_To_group_DemoList(in, out, s)
}

func autoConvert_group_DemoList_To_v1_DemoList(in *group.DemoList, out *DemoList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Demo)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_group_DemoList_To_v1_DemoList is an autogenerated conversion function.
func Convert_group_DemoList_To_v1_DemoList(in *group.DemoList, out *DemoList, s conversion.Scope) error {
	return autoConvert_group_DemoList_To_v1_DemoList(in, out, s)
}

func autoConvert_v1_DemoSpec_To_group_DemoSpec(in *DemoSpec, out *group.DemoSpec, s conversion.Scope) error {
	out.V1 = in.V1
	return nil
}

// Convert_v1_DemoSpec_To_group_DemoSpec is an autogenerated conversion function.
func Convert_v1_DemoSpec_To_group_DemoSpec(in *DemoSpec, out *group.DemoSpec, s conversion.Scope) error {
	return autoConvert_v1_DemoSpec_To_group_DemoSpec(in, out, s)
}

func autoConvert_group_DemoSpec_To_v1_DemoSpec(in *group.DemoSpec, out *DemoSpec, s conversion.Scope) error {
	out.V1 = in.V1
	return nil
}

// Convert_group_DemoSpec_To_v1_DemoSpec is an autogenerated conversion function.
func Convert_group_DemoSpec_To_v1_DemoSpec(in *group.DemoSpec, out *DemoSpec, s conversion.Scope) error {
	return autoConvert_group_DemoSpec_To_v1_DemoSpec(in, out, s)
}

func autoConvert_v1_DemoStatus_To_group_DemoStatus(in *DemoStatus, out *group.DemoStatus, s conversion.Scope) error {
	return nil
}

// Convert_v1_DemoStatus_To_group_DemoStatus is an autogenerated conversion function.
func Convert_v1_DemoStatus_To_group_DemoStatus(in *DemoStatus, out *group.DemoStatus, s conversion.Scope) error {
	return autoConvert_v1_DemoStatus_To_group_DemoStatus(in, out, s)
}

func autoConvert_group_DemoStatus_To_v1_DemoStatus(in *group.DemoStatus, out *DemoStatus, s conversion.Scope) error {
	return nil
}

// Convert_group_DemoStatus_To_v1_DemoStatus is an autogenerated conversion function.
func Convert_group_DemoStatus_To_v1_DemoStatus(in *group.DemoStatus, out *DemoStatus, s conversion.Scope) error {
	return autoConvert_group_DemoStatus_To_v1_DemoStatus(in, out, s)
}
