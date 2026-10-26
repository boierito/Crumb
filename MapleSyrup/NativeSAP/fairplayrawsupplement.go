// Additional deterministic 26A428/26A434 bootstrap block, decoded from the
// sha256-pinned CoreFP/fairplayd system binary (no JIT, no guest Apple code).
// Source build: 26A434, sha256 1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784.
// Original public CFG omitted this authenticated indirect-dispatch landing.
// This function does not claim that subsequent dynamic targets are covered.
package rawsource

import "encoding/binary"

func putBootstrapStackWord(s *state, offset uint64, register int) error {
    var value [4]byte
    binary.LittleEndian.PutUint32(value[:], uint32(s.x[register]))
    return s.store(s.sp+offset, value[:])
}

func rawSourceBootstrapBlock1008fff2c(s *state, runtime *rawSourceRuntime) (next uint64, stopped bool, err error) {
    if runtime != nil && runtime.stopPC == 0x1008fff2c { return 0x1008fff2c, true, nil }
    if err = putBootstrapStackWord(s, 0x18c, 4); err != nil { return 0, false, err }
    if runtime != nil && runtime.stopPC == 0x1008fff30 { return 0x1008fff30, true, nil }
    if err = putBootstrapStackWord(s, 0x190, 27); err != nil { return 0, false, err }
    if runtime != nil && runtime.stopPC == 0x1008fff34 { return 0x1008fff34, true, nil }
    if err = putBootstrapStackWord(s, 0x194, 25); err != nil { return 0, false, err }
    if runtime != nil && runtime.stopPC == 0x1008fff38 { return 0x1008fff38, true, nil }
    s.x[27] = s.x[21]
    if runtime != nil && runtime.stopPC == 0x1008fff3c { return 0x1008fff3c, true, nil }
    s.x[25] = s.x[19]
    if runtime != nil && runtime.stopPC == 0x1008fff40 { return 0x1008fff40, true, nil }
    if err = putBootstrapStackWord(s, 0x118, 20); err != nil { return 0, false, err }
    if runtime != nil && runtime.stopPC == 0x1008fff44 { return 0x1008fff44, true, nil }
    s.x[2] = 0
    if runtime != nil && runtime.stopPC == 0x1008fff48 { return 0x1008fff48, true, nil }
    s.x[9] = 0x64e4
    if runtime != nil && runtime.stopPC == 0x1008fff4c { return 0x1008fff4c, true, nil }
    s.x[9] = uint64(uint32(uint32(s.x[12]) + uint32(s.x[9])))
    if runtime != nil && runtime.stopPC == 0x1008fff50 { return 0x1008fff50, true, nil }
    value, loadErr := s.load(s.x[30] + (uint64(int64(int32(uint32(s.x[9]))) * 8)), 8); if loadErr != nil { return 0, false, loadErr }; s.x[6] = value
    if runtime != nil && runtime.stopPC == 0x1008fff54 { return 0x1008fff54, true, nil }
    s.x[9] = 0x64d2
    if runtime != nil && runtime.stopPC == 0x1008fff58 { return 0x1008fff58, true, nil }
    s.x[9] = uint64(uint32(uint32(s.x[12]) + uint32(s.x[9])))
    if runtime != nil && runtime.stopPC == 0x1008fff5c { return 0x1008fff5c, true, nil }
    s.x[10] = 0x8f47
    if runtime != nil && runtime.stopPC == 0x1008fff60 { return 0x1008fff60, true, nil }
    s.x[10] = uint64(uint32(uint32(s.x[10]) | (0xaf47 << 16)))
    if runtime != nil && runtime.stopPC == 0x1008fff64 { return 0x1008fff64, true, nil }
    s.x[11] = 0xd
    if runtime != nil && runtime.stopPC == 0x1008fff68 { return 0x1008fff68, true, nil }
    s.x[14] = 0x6f
    if runtime != nil && runtime.stopPC == 0x1008fff6c { return 0x1008fff6c, true, nil }
    s.x[15] = 0x4307
    if runtime != nil && runtime.stopPC == 0x1008fff70 { return 0x1008fff70, true, nil }
    s.x[15] = uint64(uint32(uint32(s.x[15]) | (1 << 16)))
    if runtime != nil && runtime.stopPC == 0x1008fff74 { return 0x1008fff74, true, nil }
    s.x[16] = 0xb0
    if runtime != nil && runtime.stopPC == 0x1008fff78 { return 0x1008fff78, true, nil }
    s.x[0] = 0x25
    if runtime != nil && runtime.stopPC == 0x1008fff7c { return 0x1008fff7c, true, nil }
    s.x[1] = 0x383d
    if runtime != nil && runtime.stopPC == 0x1008fff80 { return 0x1008fff80, true, nil }
    s.x[1] = uint64(uint32(uint32(s.x[1]) | (0x10 << 16)))
    if runtime != nil && runtime.stopPC == 0x1008fff84 { return 0x1008fff84, true, nil }
    s.x[3] = 0x49
    if runtime != nil && runtime.stopPC == 0x1008fff88 { return 0x1008fff88, true, nil }
    s.x[4] = 0xa5
    if runtime != nil && runtime.stopPC == 0x1008fff8c { return 0x1008fff8c, true, nil }
    s.x[5] = 0x6e
    if runtime != nil && runtime.stopPC == 0x1008fff90 { return 0x1008fff90, true, nil }
    value, loadErr := s.load(s.sp+0x1b8, 8); if loadErr != nil { return 0, false, loadErr }; s.x[13] = value
    if runtime != nil && runtime.stopPC == 0x1008fff94 { return 0x1008fff94, true, nil }
    s.x[17] = 0x62bb
    if runtime != nil && runtime.stopPC == 0x1008fff98 { return 0x1008fff98, true, nil }
    next = s.x[6]; return next, false, nil
}
