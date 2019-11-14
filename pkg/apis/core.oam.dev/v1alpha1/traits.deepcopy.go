package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

func (in *Trait) DeepCopyInto(out *Trait) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

func (in *Trait) DeepCopy() *Trait {
	if in == nil {
		return nil
	}
	out := new(Trait)
	in.DeepCopyInto(out)
	return out
}

func (in *Trait) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *TraitList) DeepCopyInto(out *TraitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Trait, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *TraitList) DeepCopy() *TraitList {
	if in == nil {
		return nil
	}
	out := new(TraitList)
	in.DeepCopyInto(out)
	return out
}

func (in *TraitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *TraitSpec) DeepCopyInto(out *TraitSpec) {
	*out = *in
	if in.AppliesTo != nil {
		in, out := &in.AppliesTo, &out.AppliesTo
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

func (in *TraitSpec) DeepCopy() *TraitSpec {
	if in == nil {
		return nil
	}
	out := new(TraitSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *TraitStatus) DeepCopyInto(out *TraitStatus) {
	*out = *in
	return
}

func (in *TraitStatus) DeepCopy() *TraitStatus {
	if in == nil {
		return nil
	}
	out := new(TraitStatus)
	in.DeepCopyInto(out)
	return out
}
