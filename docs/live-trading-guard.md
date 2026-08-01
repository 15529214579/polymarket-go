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

- Bot default maximum per live order: `20U`.
- Bot default maximum filled BUY notional per process: `100U`.
- Standalone `trade` maximum per order: `20U`.
- The guard is checked immediately before every order is signed.
- API startup uses only the existing-key derivation endpoint; it does not try
  to create a new API key.

The guard does not bypass Polymarket account, jurisdiction, IP, or regional
availability controls. A CLOB restriction remains a hard failure.
