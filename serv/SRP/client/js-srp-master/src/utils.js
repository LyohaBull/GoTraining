/**
 * Returns the length of the bigint in bits.
 *
 * @param {bigint} n bigint to check
 * @return {number} count of bits needed to store bigint
 */
function bitLength(n) {
    return n.toString(2).length;
}
/**
 * Returns the length of the bigint in bytes. Rounds up to the nearest byte.
 *
 * @param {bigint} n bigint to check
 * @return {number} count of bytes needed to store bigint
 */
function byteLength(n) {
    return ((bitLength(n) + 7) / 8) | 0;
}
/**
 * Deserializes a buffer of bytes into a bigint.
 *
 * @param {Uint8Array} buf buffer containing a serialized bigint
 * @return {bigint} deserialized value parsed from buf
 */
function bigintFromBytes(buf) {
    let ret = 0n;
    for (const i of buf.values()) {
        ret = (ret << 8n) + BigInt(i);
    }
    return ret;
}
/**
 * Serializes a bigint into a buffer of bytes.
 *
 * @param {bigint} v value to serialize
 * @return {Uint8Array} serialized form of v
 */
function bytesFromBigint(v) {
    const bytes = new Uint8Array(byteLength(v));
    for (let i = bytes.length - 1; v > 0; i--, v >>= 8n) {
        bytes[i] = Number(v & 0xffn);
    }
    return bytes;
}
/**
 * Returns cryptographically-safe random bytes into a buffer.
 *
 * @param {number} numBytes number of bytes
 * @return {Uint8Array} buffer containing random bytes
 */
function randomBytes(numBytes) {
    if (numBytes < 1) {
        throw new RangeError("numBytes must be >= 1");
    }
    const bytes = new Uint8Array(numBytes);
    crypto.getRandomValues(bytes);
    return bytes;
}
/**
 * Returns the smallest positive value in the multiplicative group of integers
 * modulo n that is congruent to a.
 *
 * @param {bigint} a value to find congruent value of
 * @param {bigint} n modulo of multiplicative group
 * @return {bigint} smallest positive congruent value of a in integers modulo n
 */
function toZn(a, n) {
    if (n < 1n) {
        throw new RangeError("n must be > 0");
    }
    const aZn = a % n;
    return aZn < 0n ? aZn + n : aZn;
}
/**
 * Solves for values g, x, y, such that g = gcd(a, b) and g = ax + by.
 *
 * @param {bigint} a
 * @param {bigint} b
 * @return {{g: bigint, x: bigint, y: bigint }}
 */
function eGcd(a, b) {
    if (a < 1n || b < 1n) {
        throw new RangeError("a and b must be > 0");
    }
    let x = 0n;
    let y = 1n;
    let u = 1n;
    let v = 0n;
    while (a !== 0n) {
        const q = b / a;
        const r = b % a;
        const m = x - u * q;
        const n = y - v * q;
        b = a;
        a = r;
        x = u;
        y = v;
        u = m;
        v = n;
    }
    return { g: b, x, y };
}
/**
 * Calculates the modular inverse of a in the multiplicative group of integers
 * modulo n.
 *
 * @param {bigint} a
 * @param {bigint} n
 * @return {bigint}
 */
function modInv(a, n) {
    const egcd = eGcd(toZn(a, n), n);
    if (egcd.g !== 1n) {
        throw new RangeError();
    }
    else {
        return toZn(egcd.x, n);
    }
}
/**
 * Calculates the value of x ^ y % m efficiently.
 *
 * @param {bigint} x
 * @param {bigint} y
 * @param {bigint} m
 * @return {bigint}
 */
function modPow(x, y, m) {
    if (m < 1n) {
        throw new RangeError("n must be > 0");
    }
    else if (m === 1n) {
        return 0n;
    }
    x = toZn(x, m);
    if (y < 0n) {
        return modInv(modPow(x, y >= 0 ? y : -y, m), m);
    }
    let r = 1n;
    while (y > 0) {
        if (y % 2n === 1n) {
            r = (r * x) % m;
        }
        y = y / 2n;
        x = x ** 2n % m;
    }
    return r;
}
/**
 * Concatenates multiple buffers into one new buffer.
 *
 * @param {Uint8Array[]} a buffers to concatenate
 * @return {Uint8Array} a new buffer containing the concatenated contents
 */
function concatBytes(...a) {
    let length = 0;
    for (const b of a) {
        length += b.byteLength;
    }
    const buf = new Uint8Array(length);
    let offset = 0;
    for (const b of a) {
        buf.set(b, offset);
        offset += b.byteLength;
    }
    return buf;
}
/**
 * XORs two equal-size byte arrays together.
 *
 * @param {Uint8Array[]} a buffers to concatenate
 * @return {Uint8Array} a new buffer containing the concatenated contents
 */
function xorBytes(a, b) {
    if (a.length !== b.length) {
        throw new Error('xorBytes: buffers must be same length');
    }
    const length = a.length;
    const result = new Uint8Array(length);
    for (let i = 0; i < length; i++) {
        result[i] = a[i] ^ b[i];
    }
    return result;
}
/**
 * Encodes a buffer into a hexadecimal string.
 *
 * @param {Uint8Array} buffer buffer to encode
 * @return {string} hex-encoded form of buffer
 */
function toHex(buffer) {
    return [...buffer].map((x) => x.toString(16).padStart(2, "0")).join("");
}
/**
 * Decodes a hexadecimal string into a new buffer.
 *
 * @param {string} str hexadecimal string to decode
 * @return {Uint8Array} buffer of bytes decoded from str
 */
function fromHex(str) {
    var _a, _b;
    return Uint8Array.from((_b = (_a = str.match(/.{2}/g)) === null || _a === void 0 ? void 0 : _a.map((byte) => parseInt(byte, 16))) !== null && _b !== void 0 ? _b : []);
}
/**
 * Compares two buffers with constant-time execution.
 *
 * @param {Uint8Array} a first buffer to compare
 * @param {Uint8Array} b second buffer to compare
 * @return {boolean} true if a == b, otherwise false
 */
function constantTimeCompare(a, b) {
    if (a.length !== b.length) {
        return false;
    }
    const len = a.length;
    let out = 0;
    for (let i = 0; i < len; i++) {
        out |= a[i] ^ b[i];
    }
    return out === 0;
}
/**
 * Enumeration of hash types. This is a subset of the Go hash enumeration, with
 * only algorithms supported by WebCrypto.
 */
var Hash;
(function (Hash) {
    Hash[Hash["SHA1"] = 3] = "SHA1";
    Hash[Hash["SHA256"] = 5] = "SHA256";
    Hash[Hash["SHA384"] = 6] = "SHA384";
    Hash[Hash["SHA512"] = 7] = "SHA512";
})(Hash || (Hash = {}));
/**
 * Returns the result of applying a hash to the given buffer.
 *
 * @param {Hash} hash algorithm to use
 * @param {BufferSource} data data to hash
 * @return {Promise<ArrayBuffer>} digest
 */
function hash(hash, data) {
    switch (hash) {
        case Hash.SHA1:
            return crypto.subtle.digest("SHA-1", data);
        case Hash.SHA256:
            return crypto.subtle.digest("SHA-256", data);
        case Hash.SHA384:
            return crypto.subtle.digest("SHA-384", data);
        case Hash.SHA512:
            return crypto.subtle.digest("SHA-512", data);
    }
}
async function hashInterleave(h, data) {
    let copy = new Uint8Array(data, 0, data.byteLength);
    for (var i = 0; i < copy.length; i++) {
        if (copy[i] !== 0) {
            if ((data.byteLength - i) % 2 === 1)
                i++;
            copy = new Uint8Array(data, i, data.byteLength - i);
            break;
        }
    }
    const halfl = copy.length / 2;
    const even = new Uint8Array(halfl);
    const odd = new Uint8Array(halfl);
    for (let i = 0; i < copy.length; i++) {
        even[i] = copy[i * 2];
        odd[i] = copy[i * 2 + 1];
    }
    const hash1 = new Uint8Array(await hash(h, even));
    const hash2 = new Uint8Array(await hash(h, odd));
    const result = new Uint8Array(hash1.byteLength * 2);
    for (let i = 0; i < hash1.byteLength; i++) {
        result[i * 2] = hash1[i];
        result[i * 2 + 1] = hash2[i];
    }
    return result;
}

