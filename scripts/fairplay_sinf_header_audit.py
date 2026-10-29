"""Read-only audit of the pinned 26A428 native SINF header decision.

No bytes, addresses, keys, hashes of secrets or register values are reported.
This observes a native predicate; it does not supply a credential or validate a
license. Callers must separately complete the independent native replay.
"""

MASK = 0xFFFFFFFF
DECODE_PC = 0x1004CCDF4
WORD_PC = 0x1004CCE50
COMPARE_PC = 0x1004CB590
ERROR_PC = 0x1008C7ED8


def expected_word(scalar):
    return ((((scalar + 0x5ABE3E8C) & MASK) & 0x164D6F1F) + 0xA700BF94) & MASK


def encoded_word(word):
    # Pinned instructions at 4cce3c..4cce4c, simplified without changing bits.
    return word ^ 0xCE757BC5


class HeaderAudit:
    def __init__(self):
        self.pending = None
        self.rows = []

    def observe(self, pc, registers, sp, read):
        if pc == DECODE_PC:
            # Read exactly the four bytes used by 4cce00/10/1c/30. Do not
            # interpret arbitrary memory as a key or change guest state.
            offsets = [(registers[8] + registers[i]) & MASK for i in (20, 3, 4, 5)]
            word = 0
            for offset in offsets:
                word = (word << 8) | ((read(registers[25] + offset, 1)[0] ^ registers[26]) & 255)
            self.pending = {"word": word, "word_checked": False}
        elif pc == WORD_PC and self.pending is not None:
            self.pending["word_checked"] = registers[2] == encoded_word(self.pending["word"])
        elif pc == COMPARE_PC:
            scalar = int.from_bytes(read(sp + 0xDD8, 4), "little")
            expected = expected_word(scalar)
            row = {
                "decoded_word_observed": self.pending is not None,
                "word_transform_matches_native": bool(self.pending and self.pending["word_checked"]),
                "expected_header_is_itun": (expected ^ 0xCE757BC5) == int.from_bytes(b"itun", "big"),
                "native_comparison_equal": registers[2] == expected,
                "decoded_header_is_itun": bool(self.pending and self.pending["word"] == int.from_bytes(b"itun", "big")),
                "error_destination_observed": False,
            }
            self.rows.append(row)
            self.pending = None
        elif pc == ERROR_PC and self.rows:
            self.rows[-1]["error_destination_observed"] = True

    def report(self):
        return {"schema": "crumb-native-sinf-header-audit-v1", "decisions": self.rows,
                "guest_memory_modified_by_audit": False, "license_validated_by_audit": False}

# Public encoding constant from generated_raw_source_current_sinf.go for the
# pinned binary. This is not an account key. The private key is read only from
# the replay's own native frame at the production candidate boundary.
MATERIAL_PC = 0x100234CDC
KEY_MASK = bytes.fromhex("b80f7230553fa153e71c812f65e931fe")


def sinf_field(sinf, wanted):
    def atoms(start, end):
        result = []
        while start < end:
            if end - start < 8:
                raise ValueError("invalid SINF bounds")
            size = int.from_bytes(sinf[start:start + 4], 'big')
            if size < 8 or size > end - start:
                raise ValueError("invalid SINF bounds")
            result.append((sinf[start + 4:start + 8], start + 8, start + size))
            start += size
        return result
    def unique(entries, name):
        values = [(start, end) for typ, start, end in entries if typ == name]
        if len(values) != 1:
            raise ValueError("missing or duplicate SINF field")
        return values[0]
    roots = atoms(0, len(sinf))
    if len(roots) != 1 or roots[0][0] != b'sinf':
        raise ValueError("invalid SINF root")
    start, end = unique(atoms(roots[0][1], roots[0][2]), b'schi')
    start, end = unique(atoms(start, end), wanted)
    return sinf[start:end]


class MobileHeaderComparison:
    """Compare the native work buffer; never accept a license on a match.

    Python/cryptography cannot guarantee erasure of intermediate key copies.
    This optional host research tool is never installed in the iOS runtime.
    """
    def __init__(self, sinf):
        import hashlib
        name = sinf_field(sinf, b'name').split(b'\0', 1)[0]
        iviv = sinf_field(sinf, b'iviv')
        private = sinf_field(sinf, b'priv')
        if len(iviv) != 16 or len(private) not in (440, 448):
            raise ValueError("unsupported mobile SINF geometry")
        if len(private) == 440 and private[-8:] != bytes(8):
            raise ValueError("invalid mobile SINF trailer")
        self.iv = hashlib.md5(name + iviv).digest()
        self.block = private[:-8] if len(private) == 440 else private
        self.candidates = []
        self.rows = []

    def observe(self, pc, registers, sp, read):
        if pc == MATERIAL_PC:
            from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
            material = read(registers[29] - 0x70, 16)
            # Mobile credential = model AES key XOR 0x30; mobile AES reverses
            # that XOR. Compare directly with the effective AES key, without
            # creating or exporting an alleged account credential.
            key = bytes(a ^ b for a, b in zip(material, KEY_MASK))
            decrypt = Cipher(algorithms.AES(key), modes.CBC(self.iv)).decryptor()
            first = (decrypt.update(self.block) + decrypt.finalize())[:432]
            if first not in self.candidates:
                self.candidates.append(first)
        elif pc == DECODE_PC:
            offsets = [(registers[8] + registers[i]) & MASK for i in (20, 3, 4, 5)]
            contiguous = offsets == list(range(offsets[0], offsets[0] + 4))
            matches = []
            first_matches = []
            if contiguous and registers[26] < 256:
                native = bytes(x ^ registers[26] for x in read(registers[25] + offsets[0], 432))
                matches = [native == candidate for candidate in self.candidates]
                first_matches = [native[:16] == candidate[:16] for candidate in self.candidates]
            self.rows.append({"candidate_count_at_read": len(self.candidates),
                              "native_offsets_contiguous": contiguous,
                              "native_first_block_matches_candidates": first_matches,
                              "compared_work_bytes": 432,
                              "native_work_buffer_matches_candidates": matches})

    def report(self):
        return {"schema": "crumb-native-mobile-header-comparison-v1",
                "native_candidate_count": len(self.candidates), "comparisons": self.rows,
                "license_validated_by_comparison": False, "guest_memory_modified_by_comparison": False}
