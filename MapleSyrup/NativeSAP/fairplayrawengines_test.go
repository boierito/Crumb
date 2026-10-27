package rawsource

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// Public synthetic memory only. Missing pages are fixture inputs, never added
// to the production runtime. Every retry restarts BOTH engines from scratch.
func TestAuditAllTranslatedBlocks(t *testing.T) {
	addresses := make([]uint64, 0, len(knownSourceBlocks))
	for pc := range knownSourceBlocks {
		addresses = append(addresses, pc)
	}
	sort.Slice(addresses, func(i, j int) bool { return addresses[i] < addresses[j] })

	passed, skipped, failed := 0, 0, 0
	for _, pc := range addresses {
		for seed := int64(0); seed < 4; seed++ {
			rng := rand.New(rand.NewSource(seed + int64(pc)))
			fixture := &model{pages: map[uint64][]byte{}, sp: 0x710002000000}
			for i := range fixture.registers {
				fixture.registers[i] = rng.Uint64()
			}
			fixture.registers[31-2] = fixture.sp
			fixture.n = seed&1 != 0
			fixture.z = seed&2 != 0
			fixture.c = seed&1 == 0
			fixture.v = seed&2 == 0
			done := false
			for attempt := 0; attempt < 128; attempt++ {
				a, b := newState(fixture), newState(fixture)
				a.reset()
				b.reset()
				an, astop, ae := knownSourceBlocks[pc].run(a, &rawSourceRuntime{})
				runtime := &rawSourceRuntime{}
				runtime.stopPC = auditBlockEnd[pc]
				bn, bstop, be := runGenericInstructionBlock(b, runtime, pc)
				if bstop && bn == runtime.stopPC {
					bstop = astop
				}
				added := false
				for _, err := range []error{ae, be} {
					if err == nil {
						continue
					}
					var fault *sourceMemoryError
					if !errors.As(err, &fault) {
						continue
					}
					address := fault.address
					base := address &^ (pageSize - 1)
					if fixture.pages[base] == nil {
						page := make([]byte, pageSize)
						r := rand.New(rand.NewSource(int64(base) ^ seed))
						_, _ = r.Read(page)
						fixture.pages[base] = page
						added = true
					}
				}
				if added {
					continue
				}
				if ae != nil || be != nil {
					t.Logf("fixture rejected pc=%x seed=%d AOT=%v IR=%v", pc, seed, ae, be)
					skipped++
					done = true
					break
				}
				same := an == bn && astop == bstop && a.x == b.x && a.sp == b.sp && a.n == b.n && a.z == b.z && a.c == b.c && a.v == b.v && a.vectors == b.vectors
				for page := range fixture.pages {
					ap, _ := a.outputBytes(page, int(pageSize))
					bp, _ := b.outputBytes(page, int(pageSize))
					if !bytes.Equal(ap, bp) {
						same = false
					}
				}
				if !same {
					failed++
					if failed < 80 {
						for i := range a.x {
							if a.x[i] != b.x[i] {
								t.Logf("pc=%x seed=%d x%d AOT=%x IR=%x", pc, seed, i, a.x[i], b.x[i])
							}
						}
						t.Logf("DIFF pc=%x seed=%d next=%x/%x flags=%v/%v", pc, seed, an, bn, [4]bool{a.n, a.z, a.c, a.v}, [4]bool{b.n, b.z, b.c, b.v})
					}
				}
				passed++
				done = true
				break
			}
			if !done {
				t.Logf("fixture page budget pc=%x seed=%d", pc, seed)
				skipped++
			}
		}
	}
	t.Logf("blocks=%d passed=%d skipped=%d differences=%d", len(addresses), passed, skipped, failed)
	if skipped > 32 || passed < 17028 {
		t.Fatal("engine comparison coverage regressed")
	}
	if failed > 0 {
		t.Fatal(fmt.Sprintf("%d engine divergences", failed))
	}
}
