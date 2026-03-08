import { NextResponse } from 'next/server';
import crypto from 'crypto';

const KEY_HEX = process.env.DON_ENCRYPTION_KEY ?? '';

export async function POST(request: Request) {
  const { recipient, amount } = await request.json();

  if (!recipient || !/^0x[0-9a-fA-F]{40}$/.test(recipient)) {
    return NextResponse.json({ error: 'Invalid recipient address' }, { status: 400 });
  }
  const amountNum = parseFloat(amount);
  if (!amountNum || amountNum <= 0) {
    return NextResponse.json({ error: 'Invalid amount' }, { status: 400 });
  }

  const key = Buffer.from(KEY_HEX, 'hex');
  if (key.length !== 32) {
    return NextResponse.json({ error: 'Server encryption key not configured' }, { status: 500 });
  }

  // AES-256-CTR: [12-byte nonce][20-byte ciphertext] → 32 bytes → bytes32
  // IV = nonce(12 bytes) zero-padded to 16 bytes (matches Go workflow)
  const nonce = crypto.randomBytes(12);
  const iv = Buffer.concat([nonce, Buffer.alloc(4)]);
  const cipher = crypto.createCipheriv('aes-256-ctr', key, iv);
  const addrBytes = Buffer.from(recipient.slice(2), 'hex'); // 20 bytes
  const ciphertext = Buffer.concat([cipher.update(addrBytes), cipher.final()]);

  const encryptedRecipient =
    '0x' + Buffer.concat([nonce, ciphertext]).toString('hex').padEnd(64, '0');

  // USDC has 6 decimals
  const amountRaw = BigInt(Math.round(amountNum * 1_000_000)).toString();

  return NextResponse.json({ encryptedRecipient, amountRaw });
}
