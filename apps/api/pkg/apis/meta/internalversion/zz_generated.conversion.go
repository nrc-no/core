// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by conversion-gen. DO NOT EDIT.

package internalversion

import (
  v1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
  "github.com/nrc-no/core/apps/api/pkg/conversion"
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  unsafe "unsafe"

)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*List)(nil), (*v1.List)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_List_To_v1_List(a.(*List), b.(*v1.List), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.List)(nil), (*List)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_List_To_List(a.(*v1.List), b.(*List), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ListOptions)(nil), (*v1.ListOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ListOptions_To_v1_ListOptions(a.(*ListOptions), b.(*v1.ListOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ListOptions)(nil), (*ListOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ListOptions_To_ListOptions(a.(*v1.ListOptions), b.(*ListOptions), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_List_To_v1_List(in *List, out *v1.List, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]runtime.RawExtension, len(*in))
		for i := range *in {
			if err := runtime.Convert_runtime_Object_To_runtime_RawExtension(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_List_To_v1_List is an autogenerated conversion function.
func Convert_List_To_v1_List(in *List, out *v1.List, s conversion.Scope) error {
	return autoConvert_List_To_v1_List(in, out, s)
}

func autoConvert_v1_List_To_List(in *v1.List, out *List, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]runtime.Object, len(*in))
		for i := range *in {
			if err := runtime.Convert_runtime_RawExtension_To_runtime_Object(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1_List_To_List is an autogenerated conversion function.
func Convert_v1_List_To_List(in *v1.List, out *List, s conversion.Scope) error {
	return autoConvert_v1_List_To_List(in, out, s)
}

func autoConvert_ListOptions_To_v1_ListOptions(in *ListOptions, out *v1.ListOptions, s conversion.Scope) error {
	if err := v1.Convert_labels_Selector_To_string(&in.LabelSelector, &out.LabelSelector, s); err != nil {
		return err
	}
	if err := v1.Convert_fields_Selector_To_string(&in.FieldSelector, &out.FieldSelector, s); err != nil {
		return err
	}
	out.Watch = in.Watch
	out.AllowWatchBookmarks = in.AllowWatchBookmarks
	out.ResourceVersion = in.ResourceVersion
	out.ResourceVersionMatch = v1.ResourceVersionMatch(in.ResourceVersionMatch)
	out.TimeoutSeconds = (*int64)(unsafe.Pointer(in.TimeoutSeconds))
	out.Limit = in.Limit
	out.Continue = in.Continue
	return nil
}

// Convert_ListOptions_To_v1_ListOptions is an autogenerated conversion function.
func Convert_ListOptions_To_v1_ListOptions(in *ListOptions, out *v1.ListOptions, s conversion.Scope) error {
	return autoConvert_ListOptions_To_v1_ListOptions(in, out, s)
}

func autoConvert_v1_ListOptions_To_ListOptions(in *v1.ListOptions, out *ListOptions, s conversion.Scope) error {
	if err := v1.Convert_string_To_labels_Selector(&in.LabelSelector, &out.LabelSelector, s); err != nil {
		return err
	}
	if err := v1.Convert_string_To_fields_Selector(&in.FieldSelector, &out.FieldSelector, s); err != nil {
		return err
	}
	out.Watch = in.Watch
	out.AllowWatchBookmarks = in.AllowWatchBookmarks
	out.ResourceVersion = in.ResourceVersion
	out.ResourceVersionMatch = v1.ResourceVersionMatch(in.ResourceVersionMatch)
	out.TimeoutSeconds = (*int64)(unsafe.Pointer(in.TimeoutSeconds))
	out.Limit = in.Limit
	out.Continue = in.Continue
	return nil
}

// Convert_v1_ListOptions_To_ListOptions is an autogenerated conversion function.
func Convert_v1_ListOptions_To_ListOptions(in *v1.ListOptions, out *ListOptions, s conversion.Scope) error {
	return autoConvert_v1_ListOptions_To_ListOptions(in, out, s)
}
