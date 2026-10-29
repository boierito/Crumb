"""Local, conservative explicit-data/address origin audit; never a validator.

Tracks executed instructions, not a guessed white-box calling convention.
Control dependence and same-value writes by unknown external services are not
established. Reports contain no guest values, paths, material or secret hashes.
"""
import struct
import base64
from array import array

from capstone import Cs, CS_ARCH_ARM64, CS_MODE_ARM
from capstone.arm64_const import ARM64_OP_REG
from unicorn import arm64_const as arm, UcError

MATERIAL_PC = 0x100234CDC
BASE = 0x100000000
KINDS = ('initial_context', 'readonly_image', 'mutable_image', 'native_literal',
         'service', 'sidb', 'sidd', 'sinf', 'supf', 'executable_or_other_file',
         'unclassified')
BITS = {name: 1 << i for i, name in enumerate(KINDS)}


class ByteOrigins:
    """Bounded sparse pages, avoiding a Python object graph per guest byte."""
    def __init__(self):
        self.pages = {}

    def get(self, address, default=None):
        page = self.pages.get(address >> 12)
        offset = address & 4095
        if page is None or page[0][offset] == 0:
            return default
        return page[0][offset] - 1, page[1][offset] or None

    def __setitem__(self, address, value):
        number = address >> 12
        page = self.pages.get(number)
        if page is None:
            if len(self.pages) >= 16384:
                raise ValueError('origin audit memory bound exceeded')
            page = self.pages[number] = (array('I', [0]) * 4096, array('Q', [0]) * 4096)
        origin, writer = value
        page[0][address & 4095] = origin + 1
        page[1][address & 4095] = writer or 0

    def __getitem__(self, address):
        value = self.get(address)
        if value is None:
            raise KeyError(address)
        return value

    def pop(self, address, default=None):
        value = self.get(address, default)
        page = self.pages.get(address >> 12)
        if page is not None:
            page[0][address & 4095] = 0
            page[1][address & 4095] = 0
        return value


def image_segments(image):
    if len(image) < 32 or struct.unpack_from('<I', image)[0] != 0xFEEDFACF:
        raise ValueError('expected pinned Mach-O image')
    count, length = struct.unpack_from('<II', image, 16)
    end = 32 + length
    if end > len(image):
        raise ValueError('invalid Mach-O load commands')
    pos, result = 32, []
    for _ in range(count):
        if pos + 8 > end:
            raise ValueError('truncated Mach-O command')
        command, size = struct.unpack_from('<II', image, pos)
        if size < 8 or pos + size > end:
            raise ValueError('invalid Mach-O command')
        if command == 0x19:
            if size < 72:
                raise ValueError('truncated segment')
            address, size_vm = struct.unpack_from('<QQ', image, pos + 24)
            protection = struct.unpack_from('<I', image, pos + 60)[0]
            if protection & 1:
                result.append((address, address + size_vm,
                               'mutable_image' if protection & 2 else 'readonly_image'))
        pos += size
    if pos != end:
        raise ValueError('load command census mismatch')
    return result


def canonical_register(name):
    if name in ('xzr', 'wzr'):
        return 'zero'
    if name in ('fp', 'lr'):
        return {'fp': 'x29', 'lr': 'x30'}[name]
    if name == 'wsp':
        return 'sp'
    if name.startswith('w') and name[1:].isdigit():
        return 'x' + name[1:]
    if name[:1] in ('b', 'h', 's', 'd', 'q', 'v') and name[1:].isdigit():
        return 'v' + name[1:]
    return name


class CandidateProvenance:
    def __init__(self, image):
        self.image = image
        self.segments = image_segments(image)
        self.md = Cs(CS_ARCH_ARM64, CS_MODE_ARM)
        self.md.detail = True
        self.cache = {}
        self.registers = {}
        self.memory = ByteOrigins()
        self.default_pages = {}
        self.pending = None
        self.fd = {}
        self.rows = []
        self.materials = []  # Private in-memory equality only; never serialized.
        self.visits = 0
        self.instructions = 0
        self.unsupported = 0

    def default_origin(self, address):
        page = address >> 12
        if page in self.default_pages:
            return self.default_pages[page]
        for start, end, kind in self.segments:
            if start <= address < end:
                if start % 4096 == 0 and end % 4096 == 0:
                    self.default_pages[page] = BITS[kind]
                return BITS[kind]
        return BITS['initial_context']

    def memory_origin(self, address, size):
        origin = 0
        for a in range(address, address + size):
            origin |= self.memory.get(a, (self.default_origin(a), None))[0]
        return origin

    def reg_origin(self, name):
        key = canonical_register(name)
        return BITS['native_literal'] if key == 'zero' else self.registers.get(key, BITS['initial_context'])

    def finish(self):
        if self.pending is None:
            return
        pc, writes, origin = self.pending
        for name in writes:
            self.registers[name] = origin
        self.pending = None

    def instruction(self, u, address, size, _):
        self.finish()
        if not BASE <= address < BASE + len(self.image):
            return
        self.instructions += 1
        if address == MATERIAL_PC:
            self.candidate(u)
        info = self.cache.get(address)
        if info is None:
            instruction = next(self.md.disasm(self.image[address-BASE:address-BASE+4], address, count=1), None)
            if instruction is None:
                self.unsupported += 1
                return
            reads, writes = instruction.regs_access()
            reads = {canonical_register(instruction.reg_name(r)) for r in reads}
            writes = {canonical_register(instruction.reg_name(r)) for r in writes}
            # White-box lookup outputs depend on the index as well as table
            # bytes. Keep address-register origins, even when this also
            # conservatively retains initial stack/base pointer origins.
            if instruction.mnemonic.startswith('st'):
                # str x1,[x1] still stores x1's value even though x1 is also
                # the address base. Never discard that explicit value source.
                reads.update(canonical_register(instruction.reg_name(op.reg))
                             for op in instruction.operands if op.type == ARM64_OP_REG)
            info = (reads, writes, instruction.mnemonic)
            self.cache[address] = info
        reads, writes, mnemonic = info
        origin = 0
        for name in reads:
            origin |= self.reg_origin(name)
        if not reads:
            origin |= BITS['native_literal']
        self.pending = (address, writes, origin)

    def read(self, u, access, address, size, value, _):
        if self.pending is not None:
            pc, writes, origin = self.pending
            self.pending = (pc, writes, origin | self.memory_origin(address, size))

    def write(self, u, access, address, size, value, _):
        if self.pending is None:
            return
        pc, _, origin = self.pending
        for a in range(address, address + size):
            self.memory[a] = (origin, pc)

    def candidate(self, u):
        self.visits += 1
        start = u.reg_read(arm.UC_ARM64_REG_X29) - 0x70
        value = bytes(u.mem_read(start, 16))
        if value not in self.materials:
            self.materials.append(value)
        candidate = self.materials.index(value) + 1
        entries = [self.memory.get(a, (self.default_origin(a), None)) for a in range(start, start + 16)]
        origin = self.memory_origin(start, 16)
        row = {'candidate_ordinal': candidate,
               'origins': [name for name in KINDS if origin & BITS[name]],
               'all_bytes_written_by_native': all(pc is not None for _, pc in entries),
               'writer_pcs': sorted({'text+0x%x' % (pc-BASE) for _, pc in entries if pc is not None}),
               'visits': 1}
        if self.rows and all(self.rows[-1][key] == row[key] for key in row if key != 'visits'):
            self.rows[-1]['visits'] += 1
        else:
            self.rows.append(row)

    def service(self, u, event):
        self.finish()
        cpu, after = event['CPU'], event['After']
        pc, args = cpu['PC'], cpu['Registers']
        returned = after['Registers'][0]
        # The service trace supplies page deltas, not an exact write log.
        # Attribute changed bytes; unknown same-value writes remain explicit
        # incompleteness in report(). Known read/mmap contracts override below.
        for page, encoded in (event.get('Pages') or {}).items():
            address = int(page)
            content = base64.b64decode(encoded, validate=True)
            if len(content) != 4096:
                raise ValueError('invalid service page')
            try:
                previous = bytes(u.mem_read(address, 4096))
            except UcError:
                previous = None
            for offset, byte in enumerate(content):
                if previous is None or previous[offset] != byte:
                    self.memory[address + offset] = (BITS['service'], None)
        for page in event.get('Removed') or []:
            for address in range(page, page + 4096):
                self.memory.pop(address, None)
        if cpu['Flags'] != after['Flags']:
            self.registers['nzcv'] = BITS['service']
        if pc == 0xF10000000001F000 and returned not in (0xFFFFFFFF, 0xFFFFFFFFFFFFFFFF):
            name = bytearray()
            for i in range(4096):
                c = bytes(u.mem_read(args[0] + i, 1))[0]
                if not c:
                    break
                name.append(c)
            lower = bytes(name).lower()
            kind = next((k for k in ('sidb', 'sidd', 'sinf', 'supf') if lower.endswith(b'.'+k.encode())), 'executable_or_other_file')
            self.fd[returned] = kind
        if pc == 0xF10000000001A000:
            self.fd.pop(args[0], None)
        target = None
        if pc == 0xF100000000008000 and returned <= args[2]:
            target = (args[1], returned, self.fd.get(args[0], 'service'))
        if pc == 0xF100000000023000 and returned not in (0xFFFFFFFF, 0xFFFFFFFFFFFFFFFF):
            target = (returned, args[1], self.fd.get(args[4], 'service'))
        if target:
            address, size, kind = target
            if size > 0x10000000:
                raise ValueError('invalid audited service size')
            for a in range(address, address + size):
                self.memory[a] = (BITS[kind], None)
        for index, (before, value) in enumerate(zip(args, after['Registers'])):
            if before != value or index == 0:
                self.registers['x%d' % index] = BITS['service']

    def report(self):
        return {'schema': 'crumb-candidate-provenance-v1', 'candidate_count': len(self.materials),
                'landmark_visits': self.visits, 'runs': self.rows,
                'instructions_observed': self.instructions, 'undecoded_instructions': self.unsupported,
                'control_dependencies_complete': False,
                'explicit_address_registers_included': True,
                'address_dependencies_complete': False,
                'unknown_service_same_value_writes_verified': False,
                'causal_chain_complete': False, 'license_validated': False,
                'guest_memory_modified': False}
