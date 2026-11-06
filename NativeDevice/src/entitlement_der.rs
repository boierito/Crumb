//! Bounded decoder for Apple's DER entitlement envelope and legacy dictionary.
use anyhow::{Context, Result, ensure};
fn tlv<'a>(input: &mut &'a [u8]) -> Result<(u8, &'a [u8])> {
    ensure!(input.len() >= 2, "Truncated DER entitlement");
    let tag = input[0];
    let first = input[1];
    *input = &input[2..];
    let length = if first & 0x80 == 0 {
        first as usize
    } else {
        let count = (first & 0x7f) as usize;
        ensure!(
            (1..=4).contains(&count) && input.len() >= count,
            "Invalid DER length"
        );
        ensure!(input[0] != 0, "Noncanonical DER length");
        let mut n = 0usize;
        for b in &input[..count] {
            n = (n << 8) | *b as usize;
        }
        *input = &input[count..];
        ensure!(n >= 128, "Noncanonical DER length");
        n
    };
    ensure!(length <= input.len(), "Truncated DER value");
    let body = &input[..length];
    *input = &input[length..];
    Ok((tag, body))
}
fn dictionary(mut body: &[u8], depth: usize) -> Result<plist::Dictionary> {
    ensure!(depth <= 32, "Entitlements nested too deeply");
    let mut result = plist::Dictionary::new();
    while !body.is_empty() {
        ensure!(result.len() < 10000, "Too many entitlements");
        let (tag, mut entry) = tlv(&mut body)?;
        ensure!(tag == 0x30, "Invalid DER dictionary entry");
        let (tag, key) = tlv(&mut entry)?;
        ensure!(tag == 0x0c, "Invalid DER key");
        let key = std::str::from_utf8(key)?.to_owned();
        ensure!(
            !key.is_empty() && !result.contains_key(&key),
            "Duplicate DER entitlement"
        );
        let value = value(&mut entry, depth + 1)?;
        ensure!(entry.is_empty(), "Trailing DER entry data");
        result.insert(key, value);
    }
    Ok(result)
}
fn value(input: &mut &[u8], depth: usize) -> Result<plist::Value> {
    ensure!(depth <= 32, "Entitlements nested too deeply");
    let (tag, mut body) = tlv(input)?;
    Ok(match tag {
        0x0c => std::str::from_utf8(body)?.to_owned().into(),
        0x01 => {
            ensure!(
                body.len() == 1 && [0, 255].contains(&body[0]),
                "Invalid DER boolean"
            );
            (body[0] != 0).into()
        }
        0x02 => {
            ensure!(!body.is_empty() && body.len() <= 9, "Invalid DER integer");
            if body[0] & 0x80 != 0 {
                ensure!(body.len() <= 8, "DER integer overflow");
                let mut n = -1i64;
                for b in body {
                    n = (n << 8) | *b as i64;
                }
                n.into()
            } else {
                if body.len() == 9 {
                    ensure!(body[0] == 0, "DER integer overflow");
                    body = &body[1..];
                }
                let mut n = 0u64;
                for b in body {
                    n = (n << 8) | *b as u64;
                }
                n.into()
            }
        }
        0x04 => plist::Value::Data(body.to_vec()),
        0x30 => {
            let mut values = Vec::new();
            while !body.is_empty() {
                ensure!(values.len() < 10000, "Too many DER values");
                values.push(value(&mut body, depth + 1)?);
            }
            values.into()
        }
        0x31 => dictionary(body, depth + 1)?.into(),
        _ => anyhow::bail!("Unsupported DER entitlement tag"),
    })
}
pub fn decode(mut bytes: &[u8]) -> Result<plist::Dictionary> {
    ensure!(bytes.len() <= 1024 * 1024, "DER entitlements too large");
    let (tag, mut body) = tlv(&mut bytes)?;
    ensure!(bytes.is_empty(), "Trailing DER signature data");
    if tag == 0x31 {
        return dictionary(body, 0);
    }
    ensure!(tag == 0x70, "Invalid DER envelope");
    let (version, v) = tlv(&mut body)?;
    ensure!(
        version == 2 && v == [1],
        "Unsupported DER entitlement version"
    );
    let (tag, dict) = tlv(&mut body)?;
    ensure!(
        tag == 0xb0 && body.is_empty(),
        "Invalid DER dictionary envelope"
    );
    dictionary(dict, 0).context("Invalid DER entitlements")
}
#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn apples_versioned_envelope_round_trips() {
        let mut d = plist::Dictionary::new();
        d.insert(
            "keychain-access-groups".into(),
            vec![plist::Value::String("TEAM.app".into())].into(),
        );
        d.insert("get-task-allow".into(), false.into());
        d.insert("number".into(), 42u64.into());
        let encoded = apple_codesign_quick::signature::der_entitlements(&d).unwrap();
        assert_eq!(decode(&encoded).unwrap(), d);
        for i in 0..encoded.len() {
            assert!(decode(&encoded[..i]).is_err());
        }
    }
    #[test]
    fn invalid_lengths_are_rejected() {
        for b in [vec![0x70, 0x80], vec![0x70, 0x81, 0], vec![0x70, 0xff]] {
            assert!(decode(&b).is_err());
        }
    }
}
