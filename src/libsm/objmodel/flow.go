package objmodel

import (
	"errors"
	"libsm/yamlmodel"
)

type Flow struct {
	CoreObject
	protocol map[string]FlowSpec
	sender CoreSpec
	receiver CoreSpec
}

func (f *Flow) Init(fl *yamlmodel.Flow, r Resolver) []error {
	var errs []error

	if fl == nil {
		return []error{errors.New("cannot convert nil yaml to flow specification")}
	}
	
	err := f.SetID(fl.Id)
	if err != nil { return []error{err} }
	err = f.SetName(fl.Name)
	if err != nil { return []error{err} }
	err = f.SetDescription(fl.Description)
	if err != nil { return []error{err} }

	if fl.AdmDir != "" {
		for _, adm := range fl.ADM {
			adm = fl.AdmDir + "/" + adm
			f.AddADM(adm)
		}
	} else {
		f.SetADM(fl.ADM)
	}

	
	for _, proto := range fl.Protocol {
		obj, protoErrs := r(proto)
		if len(protoErrs) != 0 {
			errs = append(errs, protoErrs...)
		}
		if p, ok := obj.(FlowSpec); ok {
			f.AddProtocol(proto, p)
		} else {
			errs = append(errs, errors.New("error in resolving protocol '" + proto + "' for flow '" + f.id + "'"))
		}
	}

	if fl.Sender != "" {
		obj, senderErrs := r(fl.Sender)
		if len(senderErrs) != 0 {
			errs = append(errs, senderErrs...)
		}
		if send, ok := obj.(CoreSpec); ok {
			f.SetSender(send)
		} else {
			errs = append(errs, errors.New("error in resolving sender '" + fl.Sender + "' for flow '" + f.id + "'"))
		}
	}

	if fl.Receiver != "" {
		obj, recvErrs := r(fl.Receiver)
		if len(recvErrs) != 0 {
			errs = append(errs, recvErrs...)
		}
		if recv, ok := obj.(CoreSpec); ok {
			f.SetReceiver(recv)
		} else {
			errs = append(errs, errors.New("error in resolving receiver '" + fl.Receiver + "' for flow '" + f.id + "'"))
		}
	}

	f.mitigations = append(f.mitigations, fl.Mitigations...)
	f.recommendations = append(f.recommendations, fl.Recommendations...)
	
	return errs
}

func (f *Flow) AddADM(adm string) error {
	f.adm = append(f.adm, adm)
	return nil
}

// Collect all ADMs from flow
func (f *Flow) GetADM() (allADM map[string][]string) {
	allADM = make(map[string][]string)
	allADM[f.id] = f.adm
	if f.protocol != nil {
		for _, proto := range f.protocol {
			allADM = merge(allADM, f.id + ".protocol", proto.GetADM())
		}
	}
	return
}

func (f *Flow) GetMitigations() map[string][]string {
	allMitigations := make(map[string][]string)
	if len(f.mitigations) > 0 {
		allMitigations[f.GetName()] = append(allMitigations[f.GetName()], f.mitigations...)
	}
	
	if f.protocol != nil {
		for _, proto := range f.protocol {
			for key, mitigations := range proto.GetMitigations() {
				if len(mitigations) > 0 {
					allMitigations[f.GetName() + " -> Protocol:" + key] = append(allMitigations[f.GetName() + " ->  Protocol:" + key], mitigations...)
				}
			}
		}
	}
	return allMitigations
}

func (f *Flow) GetRecommendations() map[string][]string {
	allRecommendations := make(map[string][]string)
	if len(f.recommendations) > 0 {
		allRecommendations[f.GetName()] = append(allRecommendations[f.GetName()], f.recommendations...)
	}
	
	if f.protocol != nil {
		for _, proto := range f.protocol {
			for key, recos := range proto.GetRecommendations() {
				if len(recos) > 0 {
					allRecommendations[f.GetName() + " -> Protocol:" + key] = append(allRecommendations[f.GetName() + " ->  Protocol:" + key], recos...)
				}
			}
		}
	}
	return allRecommendations
}

func (f *Flow) AddMitigation(mitigation string) error {
	f.mitigations = append(f.mitigations, mitigation)
	return nil
}

func (f *Flow) AddRecommendation(reco string) error {
	f.recommendations = append(f.recommendations, reco)
	return nil
}

func (f *Flow) GetProtocol() map[string]FlowSpec {
	return f.protocol
}

func (f *Flow) AddProtocol(id string, protocol FlowSpec) error {
	if f.protocol == nil { f.protocol = make(map[string]FlowSpec) }
	if _, present := f.protocol[id]; present {
		return errors.New("'" + id + "' is already part of protocols list for '" + f.id + "'")
	} else {
		f.protocol[id] = protocol
	}
	return nil
}

func (f *Flow) GetSender() CoreSpec {
	return f.sender
}

func (f *Flow) SetSender(s CoreSpec) error {
	f.sender = s
	return nil
}

func (f *Flow) GetReceiver() CoreSpec {
	return f.receiver
}

func (f *Flow) SetReceiver(r CoreSpec) error {
	f.receiver = r
	return nil
}