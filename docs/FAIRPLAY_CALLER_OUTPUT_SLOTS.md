# Pinned account-entry caller output slots

The exact public account-entry snapshot supplies four scalar output pointers on
its stack, but omits their backing pages. In the owner's original capture the
native path reached the return sequence at text+0xe057c and failed while
storing the first scalar at text+0xe0580. The next three stores are at
text+0xe0588, text+0xe0594 and text+0xe059c; the pinned instructions write
respectively 8, 4, 8 and 4 bytes through those caller pointers.

Only when the pinned binary identity, exact saved caller pointers and absent
pages match the public snapshot, create each scratch page on its exact native
store. The model rejects reads before initialization, access outside the
scalar, writes of another size or writes from another instruction. The pages
hold results of the real emulated instructions, not assumed license material;
all model pages are cleared at attempt cleanup. Other snapshots retain their
existing memory behavior.

The captured pre-authorization state progressed beyond the four writes after
this change and subsequently reached a return outside the current profile.
The one SINF candidate checked during that attempt remained rejected. This
does not establish compatibility with the later authorized state on the phone,
that the returned target is valid, or that an IPA can be decrypted.
