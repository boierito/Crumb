package rawsource

import "testing"

// Return values were independently measured on the exact daemon's native
// build host for all these class/relative-priority combinations.
func TestSourceQoSMatchesNativeValidation(t *testing.T) {
	for _, class := range []uint32{0, 5, 9, 17, 21, 25, 33, 127} {
		for _, relative := range []int32{0, -1, -15, -16, 1} {
			s := newState(&model{pages: map[uint64][]byte{}})
			s.reset()
			s.x[0] = uint64(class)
			s.x[1] = uint64(uint32(relative))
			r := &rawSourceRuntime{qosClass: 9, qosRelative: -3, qosSet: true}
			valid := (class == 5 || class == 9 || class == 17 || class == 21 || class == 25 || class == 33) && relative <= 0 && relative >= -15
			code := r.handleSourceQoS(s)
			if valid {
				if code != 0 || !r.qosSet || r.qosClass != class || r.qosRelative != relative {
					t.Fatal("valid native request not represented")
				}
			} else if code != 22 || r.qosClass != 9 || r.qosRelative != -3 || !r.qosSet {
				t.Fatal("invalid request modified virtual scheduling state")
			}
		}
	}
}
