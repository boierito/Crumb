Crumb vendors apple-codesign-quick 0.1.0 from crates.io (Dadoum/apple-crates).
License: LGPL-2.1-or-later; see COPYING.LESSER. Upstream copyright attribution
remains with its authors. Crumb's change in src/bundle.rs excludes
EmbeddedPayload from dylib signing: those files are sealed as resources so
Apple's original signatures and FairPlay material survive helper installation.
The complete source, Cargo.lock and build scripts are provided for rebuilding
and relinking the statically linked library with a modified copy.
