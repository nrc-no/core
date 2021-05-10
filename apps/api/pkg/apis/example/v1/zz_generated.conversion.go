// +build !ignore_autogenerated

// Hello!

// Code generated by conversion-gen. DO NOT EDIT.

package v1

import (
	url "net/url"

	example "github.com/nrc-no/core/apps/api/pkg/apis/example"

	conversion "github.com/nrc-no/core/apps/api/pkg/conversion"
	runtime "github.com/nrc-no/core/apps/api/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*TestModel)(nil), (*example.TestModel)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_TestModel_To_example_TestModel(a.(*TestModel), b.(*example.TestModel), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*example.TestModel)(nil), (*TestModel)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_example_TestModel_To_TestModel(a.(*example.TestModel), b.(*TestModel), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TestModel2)(nil), (*example.TestModel2)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_TestModel2_To_example_TestModel2(a.(*TestModel2), b.(*example.TestModel2), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*example.TestModel2)(nil), (*TestModel2)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_example_TestModel2_To_TestModel2(a.(*example.TestModel2), b.(*TestModel2), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TestModel2Spec)(nil), (*example.TestModel2Spec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_TestModel2Spec_To_example_TestModel2Spec(a.(*TestModel2Spec), b.(*example.TestModel2Spec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*example.TestModel2Spec)(nil), (*TestModel2Spec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_example_TestModel2Spec_To_TestModel2Spec(a.(*example.TestModel2Spec), b.(*TestModel2Spec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TestModelSpec)(nil), (*example.TestModelSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_TestModelSpec_To_example_TestModelSpec(a.(*TestModelSpec), b.(*example.TestModelSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*example.TestModelSpec)(nil), (*TestModelSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_example_TestModelSpec_To_TestModelSpec(a.(*example.TestModelSpec), b.(*TestModelSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TestModelUrlValues)(nil), (*example.TestModelUrlValues)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_TestModelUrlValues_To_example_TestModelUrlValues(a.(*TestModelUrlValues), b.(*example.TestModelUrlValues), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*example.TestModelUrlValues)(nil), (*TestModelUrlValues)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_example_TestModelUrlValues_To_TestModelUrlValues(a.(*example.TestModelUrlValues), b.(*TestModelUrlValues), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*url.Values)(nil), (*TestModelUrlValues)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_url_Values_To_TestModelUrlValues(a.(*url.Values), b.(*TestModelUrlValues), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_TestModel_To_example_TestModel(in *TestModel, out *example.TestModel, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_TestModelSpec_To_example_TestModelSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_TestModel_To_example_TestModel is an autogenerated conversion function.
func Convert_TestModel_To_example_TestModel(in *TestModel, out *example.TestModel, s conversion.Scope) error {
	return autoConvert_TestModel_To_example_TestModel(in, out, s)
}

func autoConvert_example_TestModel_To_TestModel(in *example.TestModel, out *TestModel, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_example_TestModelSpec_To_TestModelSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_example_TestModel_To_TestModel is an autogenerated conversion function.
func Convert_example_TestModel_To_TestModel(in *example.TestModel, out *TestModel, s conversion.Scope) error {
	return autoConvert_example_TestModel_To_TestModel(in, out, s)
}

func autoConvert_TestModel2_To_example_TestModel2(in *TestModel2, out *example.TestModel2, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_TestModelSpec_To_example_TestModelSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_TestModel2_To_example_TestModel2 is an autogenerated conversion function.
func Convert_TestModel2_To_example_TestModel2(in *TestModel2, out *example.TestModel2, s conversion.Scope) error {
	return autoConvert_TestModel2_To_example_TestModel2(in, out, s)
}

func autoConvert_example_TestModel2_To_TestModel2(in *example.TestModel2, out *TestModel2, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_example_TestModelSpec_To_TestModelSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_example_TestModel2_To_TestModel2 is an autogenerated conversion function.
func Convert_example_TestModel2_To_TestModel2(in *example.TestModel2, out *TestModel2, s conversion.Scope) error {
	return autoConvert_example_TestModel2_To_TestModel2(in, out, s)
}

func autoConvert_TestModel2Spec_To_example_TestModel2Spec(in *TestModel2Spec, out *example.TestModel2Spec, s conversion.Scope) error {
	out.SomeProperty = in.SomeProperty
	return nil
}

// Convert_TestModel2Spec_To_example_TestModel2Spec is an autogenerated conversion function.
func Convert_TestModel2Spec_To_example_TestModel2Spec(in *TestModel2Spec, out *example.TestModel2Spec, s conversion.Scope) error {
	return autoConvert_TestModel2Spec_To_example_TestModel2Spec(in, out, s)
}

func autoConvert_example_TestModel2Spec_To_TestModel2Spec(in *example.TestModel2Spec, out *TestModel2Spec, s conversion.Scope) error {
	out.SomeProperty = in.SomeProperty
	return nil
}

// Convert_example_TestModel2Spec_To_TestModel2Spec is an autogenerated conversion function.
func Convert_example_TestModel2Spec_To_TestModel2Spec(in *example.TestModel2Spec, out *TestModel2Spec, s conversion.Scope) error {
	return autoConvert_example_TestModel2Spec_To_TestModel2Spec(in, out, s)
}

func autoConvert_TestModelSpec_To_example_TestModelSpec(in *TestModelSpec, out *example.TestModelSpec, s conversion.Scope) error {
	out.SomeProperty = in.SomeProperty
	return nil
}

// Convert_TestModelSpec_To_example_TestModelSpec is an autogenerated conversion function.
func Convert_TestModelSpec_To_example_TestModelSpec(in *TestModelSpec, out *example.TestModelSpec, s conversion.Scope) error {
	return autoConvert_TestModelSpec_To_example_TestModelSpec(in, out, s)
}

func autoConvert_example_TestModelSpec_To_TestModelSpec(in *example.TestModelSpec, out *TestModelSpec, s conversion.Scope) error {
	out.SomeProperty = in.SomeProperty
	return nil
}

// Convert_example_TestModelSpec_To_TestModelSpec is an autogenerated conversion function.
func Convert_example_TestModelSpec_To_TestModelSpec(in *example.TestModelSpec, out *TestModelSpec, s conversion.Scope) error {
	return autoConvert_example_TestModelSpec_To_TestModelSpec(in, out, s)
}

func autoConvert_TestModelUrlValues_To_example_TestModelUrlValues(in *TestModelUrlValues, out *example.TestModelUrlValues, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.Abc = in.Abc
	return nil
}

// Convert_TestModelUrlValues_To_example_TestModelUrlValues is an autogenerated conversion function.
func Convert_TestModelUrlValues_To_example_TestModelUrlValues(in *TestModelUrlValues, out *example.TestModelUrlValues, s conversion.Scope) error {
	return autoConvert_TestModelUrlValues_To_example_TestModelUrlValues(in, out, s)
}

func autoConvert_example_TestModelUrlValues_To_TestModelUrlValues(in *example.TestModelUrlValues, out *TestModelUrlValues, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.Abc = in.Abc
	return nil
}

// Convert_example_TestModelUrlValues_To_TestModelUrlValues is an autogenerated conversion function.
func Convert_example_TestModelUrlValues_To_TestModelUrlValues(in *example.TestModelUrlValues, out *TestModelUrlValues, s conversion.Scope) error {
	return autoConvert_example_TestModelUrlValues_To_TestModelUrlValues(in, out, s)
}

func autoConvert_url_Values_To_TestModelUrlValues(in *url.Values, out *TestModelUrlValues, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	if values, ok := map[string][]string(*in)["abc"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_string(&values, &out.Abc, s); err != nil {
			return err
		}
	} else {
		out.Abc = ""
	}
	return nil
}

// Convert_url_Values_To_TestModelUrlValues is an autogenerated conversion function.
func Convert_url_Values_To_TestModelUrlValues(in *url.Values, out *TestModelUrlValues, s conversion.Scope) error {
	return autoConvert_url_Values_To_TestModelUrlValues(in, out, s)
}
