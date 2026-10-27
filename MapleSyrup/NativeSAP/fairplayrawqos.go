package rawsource

// pthread_set_qos_class_self_np in the one-thread guest. The virtual class is
// retained; no host scheduling or code execution permission is requested.
func (r *rawSourceRuntime) handleSourceQoS(s *state) uint64 {
	class, relative := uint32(s.x[0]), int32(s.x[1])
	if class != 5 && class != 9 && class != 17 && class != 21 && class != 25 && class != 33 {
		return 22
	}
	if relative > 0 || relative < -15 {
		return 22
	}
	r.qosClass, r.qosRelative, r.qosSet = class, relative, true
	return 0
}
