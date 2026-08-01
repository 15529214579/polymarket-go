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
single pending-maintenance alert while the separate task owns execution.

Wallet mutations are separate `trade` maintenance invocations that exit without
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

## Hourly Redemption

`com.polymarket-go.hourly-live-redeem` checks once per hour and runs only the
standalone `trade -redeem-all` maintenance path. It does not run at load and it
cannot submit or cancel CLOB orders.

The task is fail-closed. All of the following must be true:

- `db/live/redeem.disabled` is absent. Its presence always wins.
- `db/live/redeem.enabled` is a regular, non-symlink file owned by the current
  user with mode `0600`.
- The enable file contains exactly the expected public wallet address.
- A valid `BW_SESSION` is supplied to the launchd process environment. The
  session is never stored in the plist, repository, state file, or logs.

The installer loads the hourly schedule but deliberately does not create the
enable file or provide a Bitwarden session. Runtime status, lock data, alert
deduplication, and redeemed-position state are all under `db/live/` and are
gitignored. Paper positions, journals, PnL, and risk state retain their existing
paths and are never read or written by the redemption task.
