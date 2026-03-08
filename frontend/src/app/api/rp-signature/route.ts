import { NextResponse } from "next/server";
import { signRequest } from "@worldcoin/idkit/signing";

export async function POST(request: Request): Promise<Response> {
    const { action } = await request.json();
    // Strip any 0x prefix and whitespace — signRequest expects raw hex
    const signingKey = (process.env.RP_SIGNING_KEY ?? '')
      .trim()
      .replace(/^0x/i, '');

    const { sig, nonce, createdAt, expiresAt } = signRequest(action, signingKey);

    return NextResponse.json({
        sig,
        nonce,
        created_at: createdAt,
        expires_at: expiresAt,
    });
}
