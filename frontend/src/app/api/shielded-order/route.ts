import { NextResponse } from 'next/server';
import crypto from 'crypto';

const KEY_HEX = process.env.DON_ENCRYPTION_KEY ?? '';

export async function POST(request: Request) {
  const { side, asset, price, quantity } = await request.json();

  if (!side || !asset || !price || !quantity) {
    return NextResponse.json({ error: 'Missing order fields' }, { status: 400 });
  }

  const key = Buffer.from(KEY_HEX, 'hex');
  if (key.length !== 32) {
    return NextResponse.json({ error: 'Server encryption key not configured' }, { status: 500 });
  }

  // Generate orderId: 16 UUID bytes as first 16 bytes of bytes32, rest zeros
  const uuidBytes = crypto.randomBytes(16);
  const orderIdBytes = Buffer.concat([uuidBytes, Buffer.alloc(16)]);
  const orderId = '0x' + orderIdBytes.toString('hex');

  // Build order JSON payload (id will be injected by workflow from on-chain orderId)
  const orderPayload = JSON.stringify({ side, asset, price, quantity });

  // AES-256-CTR: [12-byte nonce][ciphertext]
  const nonce = crypto.randomBytes(12);
  const iv = Buffer.concat([nonce, Buffer.alloc(4)]);
  const cipher = crypto.createCipheriv('aes-256-ctr', key, iv);
  const plainBytes = Buffer.from(orderPayload, 'utf8');
  const ciphertext = Buffer.concat([cipher.update(plainBytes), cipher.final()]);

  const encryptedOrder = '0x' + Buffer.concat([nonce, ciphertext]).toString('hex');

  return NextResponse.json({ encryptedOrder, orderId });
}
