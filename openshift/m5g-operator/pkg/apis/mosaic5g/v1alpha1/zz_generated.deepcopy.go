// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgCnGlobal) DeepCopyInto(out *CfgCnGlobal) {
	*out = *in
	if in.V1 != nil {
		in, out := &in.V1, &out.V1
		*out = make([]CfgCnV1, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.V2 != nil {
		in, out := &in.V2, &out.V2
		*out = make([]CfgCnV2, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgCnGlobal.
func (in *CfgCnGlobal) DeepCopy() *CfgCnGlobal {
	if in == nil {
		return nil
	}
	out := new(CfgCnGlobal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgCnV1) DeepCopyInto(out *CfgCnV1) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.Snap = in.Snap
	out.OaiHss = in.OaiHss
	out.OaiMme = in.OaiMme
	out.OaiSpgw = in.OaiSpgw
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgCnV1.
func (in *CfgCnV1) DeepCopy() *CfgCnV1 {
	if in == nil {
		return nil
	}
	out := new(CfgCnV1)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgCnV2) DeepCopyInto(out *CfgCnV2) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.ApnNi = in.ApnNi
	out.OaiHss = in.OaiHss
	out.OaiMme = in.OaiMme
	out.OaiSpgwc = in.OaiSpgwc
	out.OaiSpgwu = in.OaiSpgwu
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgCnV2.
func (in *CfgCnV2) DeepCopy() *CfgCnV2 {
	if in == nil {
		return nil
	}
	out := new(CfgCnV2)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgDatabase) DeepCopyInto(out *CfgDatabase) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgDatabase.
func (in *CfgDatabase) DeepCopy() *CfgDatabase {
	if in == nil {
		return nil
	}
	out := new(CfgDatabase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgFlexran) DeepCopyInto(out *CfgFlexran) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgFlexran.
func (in *CfgFlexran) DeepCopy() *CfgFlexran {
	if in == nil {
		return nil
	}
	out := new(CfgFlexran)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgGlobal) DeepCopyInto(out *CfgGlobal) {
	*out = *in
	out.Snap = in.Snap
	out.NodeFunction = in.NodeFunction
	out.MmeIpAddr = in.MmeIpAddr
	out.EutraBand = in.EutraBand
	out.DownlinkFrequency = in.DownlinkFrequency
	out.UplinkFrequencyOffset = in.UplinkFrequencyOffset
	out.NumberRbDl = in.NumberRbDl
	out.NbAntennasTx = in.NbAntennasTx
	out.NbAntennasRx = in.NbAntennasRx
	out.TxGain = in.TxGain
	out.RxGain = in.RxGain
	out.EnbName = in.EnbName
	out.EnbId = in.EnbId
	out.ParallelConfig = in.ParallelConfig
	out.MaxRxGain = in.MaxRxGain
	out.CuPortc = in.CuPortc
	out.DuPortc = in.DuPortc
	out.CuPortd = in.CuPortd
	out.DuPortd = in.DuPortd
	out.RruPortc = in.RruPortc
	out.RruPortd = in.RruPortd
	out.RccPortc = in.RccPortc
	out.RccPortd = in.RccPortd
	out.RccRruTrPreference = in.RccRruTrPreference
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgGlobal.
func (in *CfgGlobal) DeepCopy() *CfgGlobal {
	if in == nil {
		return nil
	}
	out := new(CfgGlobal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgHssGlobal) DeepCopyInto(out *CfgHssGlobal) {
	*out = *in
	if in.V1 != nil {
		in, out := &in.V1, &out.V1
		*out = make([]CfgHssV1, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.V2 != nil {
		in, out := &in.V2, &out.V2
		*out = make([]CfgHssV2, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgHssGlobal.
func (in *CfgHssGlobal) DeepCopy() *CfgHssGlobal {
	if in == nil {
		return nil
	}
	out := new(CfgHssGlobal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgHssV1) DeepCopyInto(out *CfgHssV1) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgHssV1.
func (in *CfgHssV1) DeepCopy() *CfgHssV1 {
	if in == nil {
		return nil
	}
	out := new(CfgHssV1)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgHssV2) DeepCopyInto(out *CfgHssV2) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.ApnNi = in.ApnNi
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgHssV2.
func (in *CfgHssV2) DeepCopy() *CfgHssV2 {
	if in == nil {
		return nil
	}
	out := new(CfgHssV2)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgLlMec) DeepCopyInto(out *CfgLlMec) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgLlMec.
func (in *CfgLlMec) DeepCopy() *CfgLlMec {
	if in == nil {
		return nil
	}
	out := new(CfgLlMec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgMmeGlobal) DeepCopyInto(out *CfgMmeGlobal) {
	*out = *in
	if in.V1 != nil {
		in, out := &in.V1, &out.V1
		*out = make([]CfgMmeV1, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.V2 != nil {
		in, out := &in.V2, &out.V2
		*out = make([]CfgMmeV2, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgMmeGlobal.
func (in *CfgMmeGlobal) DeepCopy() *CfgMmeGlobal {
	if in == nil {
		return nil
	}
	out := new(CfgMmeGlobal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgMmeV1) DeepCopyInto(out *CfgMmeV1) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgMmeV1.
func (in *CfgMmeV1) DeepCopy() *CfgMmeV1 {
	if in == nil {
		return nil
	}
	out := new(CfgMmeV1)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgMmeV2) DeepCopyInto(out *CfgMmeV2) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgMmeV2.
func (in *CfgMmeV2) DeepCopy() *CfgMmeV2 {
	if in == nil {
		return nil
	}
	out := new(CfgMmeV2)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgOaiEnb) DeepCopyInto(out *CfgOaiEnb) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.MmeService = in.MmeService
	out.Snap = in.Snap
	out.EutraBand = in.EutraBand
	out.DownlinkFrequency = in.DownlinkFrequency
	out.UplinkFrequencyOffset = in.UplinkFrequencyOffset
	out.NumberRbDl = in.NumberRbDl
	out.TxGain = in.TxGain
	out.RxGain = in.RxGain
	out.PuschP0Nominal = in.PuschP0Nominal
	out.PucchP0Nominal = in.PucchP0Nominal
	out.PdschReferenceSignalPower = in.PdschReferenceSignalPower
	out.PuSch10xSnr = in.PuSch10xSnr
	out.PuCch10xSnr = in.PuCch10xSnr
	out.ParallelConfig = in.ParallelConfig
	out.MaxRxGain = in.MaxRxGain
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgOaiEnb.
func (in *CfgOaiEnb) DeepCopy() *CfgOaiEnb {
	if in == nil {
		return nil
	}
	out := new(CfgOaiEnb)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgSpgwGlobal) DeepCopyInto(out *CfgSpgwGlobal) {
	*out = *in
	if in.V1 != nil {
		in, out := &in.V1, &out.V1
		*out = make([]CfgSpgwV1, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgSpgwGlobal.
func (in *CfgSpgwGlobal) DeepCopy() *CfgSpgwGlobal {
	if in == nil {
		return nil
	}
	out := new(CfgSpgwGlobal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgSpgwV1) DeepCopyInto(out *CfgSpgwV1) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgSpgwV1.
func (in *CfgSpgwV1) DeepCopy() *CfgSpgwV1 {
	if in == nil {
		return nil
	}
	out := new(CfgSpgwV1)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgSpgwcGlobal) DeepCopyInto(out *CfgSpgwcGlobal) {
	*out = *in
	if in.V2 != nil {
		in, out := &in.V2, &out.V2
		*out = make([]CfgSpgwcV2, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgSpgwcGlobal.
func (in *CfgSpgwcGlobal) DeepCopy() *CfgSpgwcGlobal {
	if in == nil {
		return nil
	}
	out := new(CfgSpgwcGlobal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgSpgwcV2) DeepCopyInto(out *CfgSpgwcV2) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.ApnNi = in.ApnNi
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgSpgwcV2.
func (in *CfgSpgwcV2) DeepCopy() *CfgSpgwcV2 {
	if in == nil {
		return nil
	}
	out := new(CfgSpgwcV2)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgSpgwuGlobal) DeepCopyInto(out *CfgSpgwuGlobal) {
	*out = *in
	if in.V2 != nil {
		in, out := &in.V2, &out.V2
		*out = make([]CfgSpgwuV2, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgSpgwuGlobal.
func (in *CfgSpgwuGlobal) DeepCopy() *CfgSpgwuGlobal {
	if in == nil {
		return nil
	}
	out := new(CfgSpgwuGlobal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CfgSpgwuV2) DeepCopyInto(out *CfgSpgwuV2) {
	*out = *in
	out.K8sPodResources = in.K8sPodResources
	if in.K8sLabelSelector != nil {
		in, out := &in.K8sLabelSelector, &out.K8sLabelSelector
		*out = make([]K8sLabelSelectorDescription, len(*in))
		copy(*out, *in)
	}
	if in.K8sNodeSelector != nil {
		in, out := &in.K8sNodeSelector, &out.K8sNodeSelector
		*out = make([]K8sNodeSelectorDescription, len(*in))
		copy(*out, *in)
	}
	out.Realm = in.Realm
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CfgSpgwuV2.
func (in *CfgSpgwuV2) DeepCopy() *CfgSpgwuV2 {
	if in == nil {
		return nil
	}
	out := new(CfgSpgwuV2)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnV1OaiHssDescription) DeepCopyInto(out *CnV1OaiHssDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnV1OaiHssDescription.
func (in *CnV1OaiHssDescription) DeepCopy() *CnV1OaiHssDescription {
	if in == nil {
		return nil
	}
	out := new(CnV1OaiHssDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnV1OaiMmeDescription) DeepCopyInto(out *CnV1OaiMmeDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnV1OaiMmeDescription.
func (in *CnV1OaiMmeDescription) DeepCopy() *CnV1OaiMmeDescription {
	if in == nil {
		return nil
	}
	out := new(CnV1OaiMmeDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnV1OaiSpgwDescription) DeepCopyInto(out *CnV1OaiSpgwDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnV1OaiSpgwDescription.
func (in *CnV1OaiSpgwDescription) DeepCopy() *CnV1OaiSpgwDescription {
	if in == nil {
		return nil
	}
	out := new(CnV1OaiSpgwDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnV2OaiHssDescription) DeepCopyInto(out *CnV2OaiHssDescription) {
	*out = *in
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnV2OaiHssDescription.
func (in *CnV2OaiHssDescription) DeepCopy() *CnV2OaiHssDescription {
	if in == nil {
		return nil
	}
	out := new(CnV2OaiHssDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnV2OaiMmeDescription) DeepCopyInto(out *CnV2OaiMmeDescription) {
	*out = *in
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnV2OaiMmeDescription.
func (in *CnV2OaiMmeDescription) DeepCopy() *CnV2OaiMmeDescription {
	if in == nil {
		return nil
	}
	out := new(CnV2OaiMmeDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnV2OaiSpgwcDescription) DeepCopyInto(out *CnV2OaiSpgwcDescription) {
	*out = *in
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnV2OaiSpgwcDescription.
func (in *CnV2OaiSpgwcDescription) DeepCopy() *CnV2OaiSpgwcDescription {
	if in == nil {
		return nil
	}
	out := new(CnV2OaiSpgwcDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CnV2OaiSpgwuDescription) DeepCopyInto(out *CnV2OaiSpgwuDescription) {
	*out = *in
	out.Snap = in.Snap
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CnV2OaiSpgwuDescription.
func (in *CnV2OaiSpgwuDescription) DeepCopy() *CnV2OaiSpgwuDescription {
	if in == nil {
		return nil
	}
	out := new(CnV2OaiSpgwuDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneralDescription) DeepCopyInto(out *GeneralDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneralDescription.
func (in *GeneralDescription) DeepCopy() *GeneralDescription {
	if in == nil {
		return nil
	}
	out := new(GeneralDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalConf) DeepCopyInto(out *GlobalConf) {
	*out = *in
	in.ConfYaml.DeepCopyInto(&out.ConfYaml)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalConf.
func (in *GlobalConf) DeepCopy() *GlobalConf {
	if in == nil {
		return nil
	}
	out := new(GlobalConf)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *K8sLabelSelectorDescription) DeepCopyInto(out *K8sLabelSelectorDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new K8sLabelSelectorDescription.
func (in *K8sLabelSelectorDescription) DeepCopy() *K8sLabelSelectorDescription {
	if in == nil {
		return nil
	}
	out := new(K8sLabelSelectorDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *K8sNodeSelectorDescription) DeepCopyInto(out *K8sNodeSelectorDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new K8sNodeSelectorDescription.
func (in *K8sNodeSelectorDescription) DeepCopy() *K8sNodeSelectorDescription {
	if in == nil {
		return nil
	}
	out := new(K8sNodeSelectorDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MmeServiceDescription) DeepCopyInto(out *MmeServiceDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MmeServiceDescription.
func (in *MmeServiceDescription) DeepCopy() *MmeServiceDescription {
	if in == nil {
		return nil
	}
	out := new(MmeServiceDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mosaic5g) DeepCopyInto(out *Mosaic5g) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mosaic5g.
func (in *Mosaic5g) DeepCopy() *Mosaic5g {
	if in == nil {
		return nil
	}
	out := new(Mosaic5g)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Mosaic5g) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mosaic5gList) DeepCopyInto(out *Mosaic5gList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Mosaic5g, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mosaic5gList.
func (in *Mosaic5gList) DeepCopy() *Mosaic5gList {
	if in == nil {
		return nil
	}
	out := new(Mosaic5gList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Mosaic5gList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mosaic5gSpec) DeepCopyInto(out *Mosaic5gSpec) {
	*out = *in
	if in.OaiEnb != nil {
		in, out := &in.OaiEnb, &out.OaiEnb
		*out = make([]CfgOaiEnb, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Flexran != nil {
		in, out := &in.Flexran, &out.Flexran
		*out = make([]CfgFlexran, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LlMec != nil {
		in, out := &in.LlMec, &out.LlMec
		*out = make([]CfgLlMec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Database != nil {
		in, out := &in.Database, &out.Database
		*out = make([]CfgDatabase, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.OaiCn.DeepCopyInto(&out.OaiCn)
	in.OaiHss.DeepCopyInto(&out.OaiHss)
	in.OaiMme.DeepCopyInto(&out.OaiMme)
	in.OaiSpgw.DeepCopyInto(&out.OaiSpgw)
	in.OaiSpgwc.DeepCopyInto(&out.OaiSpgwc)
	in.OaiSpgwu.DeepCopyInto(&out.OaiSpgwu)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mosaic5gSpec.
func (in *Mosaic5gSpec) DeepCopy() *Mosaic5gSpec {
	if in == nil {
		return nil
	}
	out := new(Mosaic5gSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mosaic5gStatus) DeepCopyInto(out *Mosaic5gStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mosaic5gStatus.
func (in *Mosaic5gStatus) DeepCopy() *Mosaic5gStatus {
	if in == nil {
		return nil
	}
	out := new(Mosaic5gStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourcesDescription) DeepCopyInto(out *ResourcesDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourcesDescription.
func (in *ResourcesDescription) DeepCopy() *ResourcesDescription {
	if in == nil {
		return nil
	}
	out := new(ResourcesDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SnapDescription) DeepCopyInto(out *SnapDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SnapDescription.
func (in *SnapDescription) DeepCopy() *SnapDescription {
	if in == nil {
		return nil
	}
	out := new(SnapDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SnapDescriptionFinal) DeepCopyInto(out *SnapDescriptionFinal) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SnapDescriptionFinal.
func (in *SnapDescriptionFinal) DeepCopy() *SnapDescriptionFinal {
	if in == nil {
		return nil
	}
	out := new(SnapDescriptionFinal)
	in.DeepCopyInto(out)
	return out
}
