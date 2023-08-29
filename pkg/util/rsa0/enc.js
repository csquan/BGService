let crypto;
try {
  crypto = require("crypto");
} catch (err) {
  console.log("crypto support is disabled!");
}
const { publicKey, privateKey } = crypto.generateKeyPairSync("rsa", {
  // The standard secure default length for RSA keys is 2048 bits
  modulusLength: 2048,
});

let pubStr = publicKey.export({
  type: "pkcs1",
  format: "pem",
});
let prvStr = privateKey.export({
  type: "pkcs1",
  format: "pem",
});

//console.log(pubStr, prvStr);
const prvB64 =
  "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBelhvSzBWbm9BY3hCVTJpZXZyTkNVVnFzTmVhRWNOZVBEUWNadGhhc2g0ZEtMU3owCmhWNmFaSFhoUWNuRzdENEdjMVlFaUhEQUxnZUI5SW9kbmNwcHJOb0VtbXpNTi9Mait5MEtZb2kwNnY4RjIveHUKTGFnK3pFTVgwWGdQUVkwb2g3alhHUU5zM05jZnFneWNlUFR0b1BBYjJlaVU2aEZ5S2VpVnMwOWpTcGkrS1BOVQpGVzBYVWpFOFpHNnk3QVJ6b3lLOWxNU0RDSkFlSXUyMlJDYmF3UjNtY1FuYlREeWRvWjJnZHl3OTk5YTFLcnd4CmRuWUtraFR3Y016aGJ5akt4WTh2ZjkycGpRYlJ5Lzc4N0E4cjhaeDNOQ3R4TXprZkltbno5SW93NGpMZDd2d3AKcEREVXdJZXhYR29jcG12UDZyTCtBTENrUnZJTzRMRldyLyt0NVFJREFRQUJBb0lCQVFDRjJTRWczS01RUU1DeApGNjZJNVBhblRoeHVCQnB6ODBjQklvWHlJblVDMS9OTzJFRDBlL3F3eEVVVytOeTIyWFNVUHcwMjM5T2dhRTJVClFVQW5vRW9VU2ZURHIybWNiSUZzQmh3RXNhN0FnWjJNZnJRNjA2VU1wQXZvN0g5cm9Rc3MxaEJ2LzlZelNZTWIKMzVreUJjS2htcWRaM0hMY3dyNk9aQzdZSGplOHgvWVJYVlRQOEpFZEVQVnFnUjlkWnl1SkFhQzFGYmg0eSsyWApiY0lycmgvY0FsMUZQQWR4UU03cWxuYjdIZmtmWWdDeFdaczJwTUZRUXVqRkVvZDZiUSt5WWNYYzhCb2RzcEdDCmdDTzQxTVFOTjA4cExDTDdyOG9sanBybzdPeDl4MVdVZ244ZDFOYjhOUzR1VFl1cXBuRWRZUXNYU3VSb0gzODcKeDVzZTZCbWxBb0dCQU9XSjNTeGhUNnI5ZHFBSHpFbEFXTkpMTUtmbTBaM083Z3pYa01kT2NTMHByQk56eWVIMwpDZXQrMnNGQTdBMWgxM3FmVzBTdnd0Q1p3clRBV3FsazVML1pMTDQyU1pFOGJtN1B2dFJiWUlwTVVYNW9BTWJQCjZjRFllR01nYWwxQlFSek5kTWZ0Vy9RUGxKOE1ReWVTTlFkVVJrVzhLWkNLLzh5VDJSSkgzcDZ2QW9HQkFPVXEKRVI5MmNCWHd3T2VBa1paUDNaYVBoVjVQMFdPaFEzbWxEVURYS1UwSXFpcVc5dVZFUFBxZE0vNTJiU21tMkx4MgoyRE1TbEwxMmdZUVpHa21xRUw5SHBNZitZemozcCtlUnYxbS8wWDhLSGdCU0FqQVBDcHpPbnBPR0hhVzlGK1dCCjFWME1BRkZxSlpvMTc5QWU4UFJzZGEwRnAxS1BiQnQ3anlJdVl3R3JBb0dBU2J1cmJHSWw4VXRTRzczbGhYSkMKRmV0SlNlWC9WNjN0RWZyODZzanIyaElVMEhyVlV0ekVOdjJjejQ5SFJGVTFucElxQXpwaVhoZkdUOEdxWGRlbgpFMmx6MGZZbVU1MFI3RTZYZ2llSUwyU3NtT3BYdFlWOEZSSjBPWU5rSjJpYXZlSFJyWmMxZm9TeXZSUjNUZkxOClRmbG9TV1pVQTdaaXpSaUJGam8zN01rQ2dZRUFoZUFUYkx2Mk12c1kxcVZuWjlaMGJ4YWRKVUdmNDRJOE52NVQKUmNQc250SW5Cd1oyYWUxNWFqY1lQdG5VWC9iV3V4TDZycXQyTlZEYnpONFZXMTk4dFNJWGc3WjdKTGFaWWxEawo1bnVHMlo2QmRGSjBjTHI0eWk1eXVXQXFSYjY0RFIzU0ZhK0RLQXpJdHRRM3F1L0llQ0k4aEwyK3lCNTlXM2pOClgzeVYzazBDZ1lFQXlPWEYrYTVSaWwwUDcvOGxPQ2hFYlFPUEJhdVJkaG5uOGs2NGVkVGE4WXczYWF0V1dkQUkKZ2tOVW90SHNiSlZUZ0hGWjJvS2lqOXhQa0lUaEJTcm1kbFBrZjdVbE9QeHhsZjBFUlYwQUkrbUN4MUhHWWp5bwprSWxFZjhpdlQvUFlJUGpjQ3E0M1Bha2paSFR4dTZDRTNCaEpFVUU5WTBDdVFoa1hBcFVZb2xnPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=";
const pubB64 =
  "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF6WG9LMFZub0FjeEJVMmlldnJOQwpVVnFzTmVhRWNOZVBEUWNadGhhc2g0ZEtMU3owaFY2YVpIWGhRY25HN0Q0R2MxWUVpSERBTGdlQjlJb2RuY3BwCnJOb0VtbXpNTi9Mait5MEtZb2kwNnY4RjIveHVMYWcrekVNWDBYZ1BRWTBvaDdqWEdRTnMzTmNmcWd5Y2VQVHQKb1BBYjJlaVU2aEZ5S2VpVnMwOWpTcGkrS1BOVUZXMFhVakU4Wkc2eTdBUnpveUs5bE1TRENKQWVJdTIyUkNiYQp3UjNtY1FuYlREeWRvWjJnZHl3OTk5YTFLcnd4ZG5ZS2toVHdjTXpoYnlqS3hZOHZmOTJwalFiUnkvNzg3QThyCjhaeDNOQ3R4TXprZkltbno5SW93NGpMZDd2d3BwRERVd0lleFhHb2NwbXZQNnJMK0FMQ2tSdklPNExGV3IvK3QKNVFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==";

const encryptStringWithRsaPublicKey = function (toEncrypt, publicKey) {
  const encrypted = crypto.publicEncrypt(
    {
      key: publicKey,
      padding: crypto.constants.RSA_PKCS1_OAEP_PADDING,
      oaepHash: "sha512",
    },
    Buffer.from(toEncrypt, "utf8")
  );
  return encrypted.toString("base64");
};
console.log(
  encryptStringWithRsaPublicKey("MyText", Buffer.from(pubB64, "base64"))
);

const Sign1 =
  "3145d338222b516c9b0a6d774dc269f6a140af64709036f678524763a6fe56f5968f7774d819c2009d7cf92ac97da568747615e765ad130205c622a8396123f0cc23c158c4b42f334136e2804baee1070ff8fa369240e31f7c71afba7fd993abfbe4f77d40da8ad2bb37db3985d090c289a0267a57e5d9787b8969d689e7b8a1cd861f88064e0c490fc416d7cbd321b3c3a8841c826cbadfaf97d0b8c617089cf23115a9f3cb7d26cdca4ab27e1a052f5160ea8c9ac0708e43737964f90e513903e7beed56f91953404f14d313a9b3af44e9da38f81ab37b33c9c39b77016ebbbc3ccdac6af7fbb43bd5adec664988a36fdb0d3535ac4ee45aa62fefaebf0b4f";
const sign2 =
  "1d791be9f7a687ab7af1483f3a50d2541e46abbf27b3a3ef5e2f6191c323fce0545247aa2b0ce3ec0159817f6b2855885144001f7d548cfbf876991fe26daa1a465e8ead085f6a0ba1cd7cb7eab6da7037452b6b024bd10547914d32c415f9f7cb30b55ea17acd2f82570b8965ba6174dd9236e4cb317423477a4ea70d0a7b75ffc5585c55a792bb32137cbeb5957c5ddebe75b49afc93cefcd447abf8cc9781862b79ce1e6b072fbfdc912a24c816994b9c33696558999b90e88b960fd26ca5902ab7ea045d0e4f427946c67274ff7f786c8d5502bb944bea53dd5d918b202b52c1117da03dab8f09276d208cdc00aa2ac68447d25081c3ba8657d6443a5743";
const sign3 =
  "9be813f5acffc1c0b1aa7dc4e1baabe54cdfdb49e07797c3e024028af42a2cc460bf8b5b69eee4cf83fc0107b201dd62db565fbf5b8edd8eb484b967b386fc75152bb81656b02ae67301c36bd8ca9777b509aaadb405bee28d92fc2dc541bfe384de9a8708006718ee34e2e21ec205056bb9203918bbb6a8f143f443b297459a303199d4932d7595c1113002b535a74e19d7d9b0bef7f37060f20d6cbd66b262de532f4409ff6505c9ddaac4319f1ce7d8e48bdfdbc963236a1ce2d047110d18e7155196577d0a1abb53a2fe1a0e9eff5f9a94efc3bec2aff88fd015c31f26def6052bccc11bceda4df5c0efd6def4aa6e6b435455c0615d8a64522792597d74";
const sign4 =
  "aec49e481f48e2169d9e0bf62ae7caed5bbbeaa8ac9412516700917139aefc1f7334fff6af1920a771b9410e58f1e8dbee64f18df95f5073ce4641689fe027cedda5763a95b510e63c22adda67df9f4d4acc16b51062300f65085018e9a74b213394085019262b24cde9322ad1581e09933422892dd9f4731319e664984d8fac9d979f1ed005676cf61788ed487c0705542917d51f16a4f0a2c5f389604628acb0bf48e47c3ac77559440ccbe9892e760d0a8637117beeab6a2969df06c0ba561521e56be8b5b665be580b9f23d648ed1f27515b47d3b1b4a560ef66b0b0dc9c6ad575b9cb09c628370dff599187ecdd8a38918d359501082cc618deece3aa10";

const m0 = { hello: "你好", world: "世界" };

// Convert a hex string to a byte array
function hexToBytes(hex) {
  let bytes = [];
  let c = 0;
  for (; c < hex.length; c += 2) bytes.push(parseInt(hex.substr(c, 2), 16));
  return bytes;
}

// Convert a byte array to a hex string
function bytesToHex(bytes) {
  let hex = [];
  let i = 0;
  for (; i < bytes.length; i++) {
    const current = bytes[i] < 0 ? bytes[i] + 256 : bytes[i];
    hex.push((current >>> 4).toString(16));
    hex.push((current & 0xf).toString(16));
  }
  return hex.join("");
}
// To verify the data, we provide the same hashing algorithm and
// padding scheme we provided to generate the signature, along
// with the signature itself, the data that we want to
// verify against the signature, and the public key
const isVerified = crypto.verify(
  "SHA512",
  Buffer.from(JSON.stringify(m0)),
  {
    key: publicKey,
    padding: crypto.constants.RSA_PKCS1_PSS_PADDING,
    saltLength: 190
  },
  Buffer.from(hexToBytes(Sign1))
);
console.log(isVerified);
