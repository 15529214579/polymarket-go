# Live Trading Guard

Live order submission is fail-closed. Paper trading, reports, and push alerts do
not use this guard.

## Required State

- `db/live-trading.disabled` must not exist. If it appears while the bot is
  running, the live process is cancelled within one second.
- `db/live-trading.enabled` must be a regular file with mode `0600`.
- The arm file must bind to the loaded wallet and expire within 24 hours of its
  `armed_at` timestamp.
- The arm file is gitignored and must never contain a key, mnemonic, API key,
  token, or Bitwarden session.

Arm-file schema:

```json
{
  "wallet": "0x015282e9b720E072A9B87eEeaE738C6Bb039Bd9e",
  "armed_at": "2026-08-02T01:00:00+08:00",
  "expires_at": "2026-08-02T13:00:00+08:00"
}
```

## Hard Limits

- Bot default maximum per live BUY: `20U`; SELL exits are not capped by the
  entry limit.
- Bot default maximum filled BUY notional per process: `100U`.
- Standalone `trade` maximum per BUY: `20U`.
- The guard is checked immediately before every order is signed.
- API startup uses only the existing-key derivation endpoint; it does not try
  to create a new API key.

## Wallet Maintenance Boundary

The automated bot creates a read-only Polygon client. It can query pUSD and
conditional-token balances, but it cannot wrap, approve, redeem, transfer, sign
an on-chain transaction, or broadcast one. Redeemable positions generate a
single manual-maintenance alert.

Wallet mutations are separate, explicit `trade` invocations that exit without
placing a CLOB order:

```text
trade -wrap-approve
trade -redeem-all
```

`-wrap-approve` wraps USDC.e and ensures exchange approvals. It no longer sizes
or submits an order. `-redeem-all` only processes positive-value redeemable
positions. Both paths verify the loaded Bitwarden wallet against the configured
expected address before constructing a transaction-capable Polygon client.

The guard does not bypass Polymarket account, jurisdiction, IP, or regional
availability controls. A CLOB restriction remains a hard failure.
