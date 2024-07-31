package validations

import (
	errors "github.com/zgalor/weberr"
	v1 "k8s.io/api/core/v1"
)

func ValidateMachinePoolTaintEffect(effect v1.TaintEffect) error {
	switch effect {
	case v1.TaintEffectNoExecute, v1.TaintEffectNoSchedule, v1.TaintEffectPreferNoSchedule:
		return nil
	default:
		return errors.BadRequest.UserErrorf("Unrecognized taint effect '%s', "+
			"only the following effects are supported: '%s', '%s', '%s'.",
			effect,
			v1.TaintEffectNoExecute, v1.TaintEffectNoSchedule, v1.TaintEffectPreferNoSchedule)
	}
}
