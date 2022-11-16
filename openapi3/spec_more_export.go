package openapi3

func (sm *SpecMore) ExportByTags() (map[string]*Spec, error) {
	specs := map[string]*Spec{}
	if sm.Spec == nil {
		return specs, ErrSpecNotSet
	}
	tags := sm.Tags(false, true)

	for _, tag := range tags {
		tagSpec, err := sm.ExportByTag(tag)
		if err != nil {
			return specs, err
		}
		if tagSpec != nil {
			specs[tag] = tagSpec
		}
	}
	return specs, nil
}

func (sm *SpecMore) ExportByTag(tag string) (*Spec, error) {
	if (sm.Spec) == nil {
		return nil, ErrSpecNotSet
	}
	oms := sm.OperationMetas([]string{tag})
	tagSpec := &Spec{}
	if len(oms) == 0 {
		return nil, nil
	}
	for _, om := range oms {
		op, err := sm.OperationByPathMethod(om.Path, om.Method)
		if err != nil {
			return nil, err
		} else if op == nil {
			continue
		}
		tagSpec.AddOperation(om.Path, om.Method, op)
	}
	return Copy(tagSpec)
}
