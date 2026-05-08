#!/usr/bin/env python3
"""
Copytrade backtest — pull historical trades for 73 whale wallets,
simulate copy-trading, and rank wallets by profitability.

Usage:
  python3 scripts/copytrade_backtest.py [--days 90] [--min-size 100] [--output reports/copytrade_backtest.md]
"""

import json, os, sys, time, argparse, math
from datetime import datetime, timezone, timedelta
from pathlib import Path
from urllib.request import urlopen, Request
from urllib.error import HTTPError
from collections import defaultdict
from concurrent.futures import ThreadPoolExecutor, as_completed

SGT = timezone(timedelta(hours=8))
BASE_URL = "https://data-api.polymarket.com"
LIMIT = 500
MAX_PAGES = 20  # safety cap per wallet
CONCURRENCY = 5
CACHE_DIR = Path(__file__).parent.parent / "db" / "whale_history"

def load_wallets(path):
    wallets = []
    with open(path) as f:
        for line in f:
            line = line.split("#")[0].strip()
            if line.startswith("0x"):
                wallets.append(line.lower())
    return wallets

def fetch_activity(addr, limit=LIMIT, offset=0):
    url = f"{BASE_URL}/activity?user={addr}&limit={limit}&offset={offset}"
    req = Request(url, headers={"User-Agent": "polymarket-go-backtest/1.0"})
    try:
        with urlopen(req, timeout=30) as resp:
            return json.loads(resp.read())
    except HTTPError as e:
        print(f"  HTTP {e.code} for {addr[:10]}… offset={offset}", file=sys.stderr)
        return []

def fetch_all_trades(addr, min_ts=0):
    """Fetch all trades for a wallet, paginating until exhausted or min_ts reached."""
    all_trades = []
    offset = 0
    for _ in range(MAX_PAGES):
        batch = fetch_activity(addr, LIMIT, offset)
        if not batch:
            break
        all_trades.extend(batch)
        oldest = min(t.get("timestamp", 0) for t in batch)
        if oldest < min_ts or len(batch) < LIMIT:
            break
        offset += LIMIT
        time.sleep(0.3)
    return all_trades

def fetch_positions(addr):
    url = f"{BASE_URL}/positions?user={addr}&sizeThreshold=0.1&limit=200"
    req = Request(url, headers={"User-Agent": "polymarket-go-backtest/1.0"})
    try:
        with urlopen(req, timeout=30) as resp:
            return json.loads(resp.read())
    except HTTPError:
        return []

def cache_path(addr):
    return CACHE_DIR / f"{addr}.jsonl"

def load_cache(addr):
    p = cache_path(addr)
    if not p.exists():
        return []
    trades = []
    with open(p) as f:
        for line in f:
            line = line.strip()
            if line:
                trades.append(json.loads(line))
    return trades

def save_cache(addr, trades):
    CACHE_DIR.mkdir(parents=True, exist_ok=True)
    seen = set()
    unique = []
    for t in trades:
        key = t.get("transactionHash", "") + str(t.get("timestamp", "")) + t.get("asset", "")
        if key not in seen:
            seen.add(key)
            unique.append(t)
    unique.sort(key=lambda t: t.get("timestamp", 0))
    with open(cache_path(addr), "w") as f:
        for t in unique:
            f.write(json.dumps(t) + "\n")
    return unique

def pull_wallet(addr, min_ts):
    """Pull trades for one wallet, merging with cache."""
    cached = load_cache(addr)
    fresh = fetch_all_trades(addr, min_ts)
    merged = cached + fresh
    unique = save_cache(addr, merged)
    return addr, unique

def simulate_copytrade(trades, min_size=100, min_price=0.05, max_price=0.95):
    """
    Simulate copy-trading on a wallet's trade history.
    For each qualifying BUY: paper-buy at that price.
    Track positions by asset. On SELL: paper-sell proportionally.
    On REDEEM: settle at 1.0 (win) or 0.0 (loss).
    """
    positions = {}  # asset -> {qty, avg_price, cost}
    closed_trades = []

    buy_trades = [t for t in trades if t.get("side") == "BUY" and t.get("type") == "TRADE"]
    sell_trades = [t for t in trades if t.get("side") == "SELL" and t.get("type") == "TRADE"]
    redeems = [t for t in trades if t.get("type") == "REDEEM"]

    for t in sorted(trades, key=lambda x: x.get("timestamp", 0)):
        side = t.get("side", "")
        typ = t.get("type", "")
        price = t.get("price", 0)
        size = t.get("size", 0)
        asset = t.get("asset", "")
        title = t.get("title", "")
        outcome = t.get("outcome", "")
        ts = t.get("timestamp", 0)
        usd_value = size * price if price > 0 else size

        if typ == "TRADE" and side == "BUY":
            if usd_value < min_size:
                continue
            if price < min_price or price > max_price:
                continue
            if asset not in positions:
                positions[asset] = {"qty": 0, "cost": 0, "avg_price": 0, "title": title, "outcome": outcome, "entry_ts": ts}
            pos = positions[asset]
            old_cost = pos["cost"]
            old_qty = pos["qty"]
            pos["qty"] += size
            pos["cost"] += size * price
            pos["avg_price"] = pos["cost"] / pos["qty"] if pos["qty"] > 0 else 0

        elif typ == "TRADE" and side == "SELL":
            if asset in positions and positions[asset]["qty"] > 0:
                pos = positions[asset]
                sell_qty = min(size, pos["qty"])
                entry_price = pos["avg_price"]
                pnl = sell_qty * (price - entry_price)
                closed_trades.append({
                    "asset": asset,
                    "title": pos["title"],
                    "outcome": pos["outcome"],
                    "entry_price": entry_price,
                    "exit_price": price,
                    "qty": sell_qty,
                    "pnl": pnl,
                    "exit_reason": "sell",
                    "entry_ts": pos["entry_ts"],
                    "exit_ts": ts,
                })
                pos["qty"] -= sell_qty
                pos["cost"] = pos["qty"] * pos["avg_price"]
                if pos["qty"] <= 0:
                    del positions[asset]

        elif typ == "REDEEM":
            if asset in positions and positions[asset]["qty"] > 0:
                pos = positions[asset]
                pnl = pos["qty"] * (1.0 - pos["avg_price"])
                closed_trades.append({
                    "asset": asset,
                    "title": pos["title"],
                    "outcome": pos["outcome"],
                    "entry_price": pos["avg_price"],
                    "exit_price": 1.0,
                    "qty": pos["qty"],
                    "pnl": pnl,
                    "exit_reason": "redeem_win",
                    "entry_ts": pos["entry_ts"],
                    "exit_ts": ts,
                })
                del positions[asset]

    # remaining open positions — use current price from last trade or mark as open
    open_positions = []
    for asset, pos in positions.items():
        open_positions.append({
            "asset": asset,
            "title": pos["title"],
            "outcome": pos["outcome"],
            "avg_price": pos["avg_price"],
            "qty": pos["qty"],
            "cost": pos["cost"],
            "entry_ts": pos["entry_ts"],
        })

    return closed_trades, open_positions

def wallet_stats(closed, open_pos):
    if not closed:
        return None
    total_pnl = sum(t["pnl"] for t in closed)
    wins = [t for t in closed if t["pnl"] > 0]
    losses = [t for t in closed if t["pnl"] <= 0]
    win_rate = len(wins) / len(closed) * 100 if closed else 0
    avg_win = sum(t["pnl"] for t in wins) / len(wins) if wins else 0
    avg_loss = sum(t["pnl"] for t in losses) / len(losses) if losses else 0
    total_cost = sum(t["qty"] * t["entry_price"] for t in closed)
    roi = (total_pnl / total_cost * 100) if total_cost > 0 else 0

    pnls = [t["pnl"] for t in closed]
    if len(pnls) > 1:
        mean_pnl = sum(pnls) / len(pnls)
        var = sum((p - mean_pnl) ** 2 for p in pnls) / (len(pnls) - 1)
        std = math.sqrt(var) if var > 0 else 0
        sharpe = mean_pnl / std if std > 0 else 0
    else:
        sharpe = 0

    open_cost = sum(p["cost"] for p in open_pos)

    return {
        "trades": len(closed),
        "wins": len(wins),
        "losses": len(losses),
        "win_rate": win_rate,
        "total_pnl": total_pnl,
        "avg_win": avg_win,
        "avg_loss": avg_loss,
        "roi": roi,
        "sharpe": sharpe,
        "open_positions": len(open_pos),
        "open_cost": open_cost,
    }

def tier_label(stats):
    if stats is None or stats["trades"] < 5:
        return "D", "insufficient data"
    if stats["roi"] > 20 and stats["win_rate"] > 55:
        return "A", "high ROI + high WR"
    if stats["roi"] > 10 or (stats["win_rate"] > 50 and stats["total_pnl"] > 0):
        return "B", "profitable"
    if stats["total_pnl"] > -50 and stats["win_rate"] > 40:
        return "C", "marginal"
    return "D", "unprofitable"

def main():
    parser = argparse.ArgumentParser(description="Copytrade historical backtest")
    parser.add_argument("--days", type=int, default=90, help="How many days back to pull (default: 90)")
    parser.add_argument("--min-size", type=float, default=100, help="Min trade size USD to copy (default: 100)")
    parser.add_argument("--min-price", type=float, default=0.05, help="Min price filter (default: 0.05)")
    parser.add_argument("--max-price", type=float, default=0.95, help="Max price filter (default: 0.95)")
    parser.add_argument("--output", default="reports/copytrade_backtest.md", help="Output report path")
    parser.add_argument("--wallets", default="wallets.txt", help="Wallets file path")
    parser.add_argument("--pull-only", action="store_true", help="Only pull data, skip analysis")
    args = parser.parse_args()

    project_root = Path(__file__).parent.parent
    wallets_path = project_root / args.wallets
    wallets = load_wallets(wallets_path)
    print(f"Loaded {len(wallets)} wallets")

    cutoff = int(time.time()) - args.days * 86400

    # Phase 1: Pull data
    print(f"\n=== Phase 1: Pulling historical trades (last {args.days} days) ===")
    wallet_trades = {}

    with ThreadPoolExecutor(max_workers=CONCURRENCY) as pool:
        futures = {pool.submit(pull_wallet, w, cutoff): w for w in wallets}
        done = 0
        for fut in as_completed(futures):
            done += 1
            addr, trades = fut.result()
            wallet_trades[addr] = trades
            trade_count = len(trades)
            buys = sum(1 for t in trades if t.get("side") == "BUY" and t.get("type") == "TRADE")
            print(f"  [{done}/{len(wallets)}] {addr[:10]}… {trade_count} trades ({buys} buys)")

    total_trades = sum(len(v) for v in wallet_trades.values())
    print(f"\nTotal: {total_trades} trades across {len(wallets)} wallets")

    if args.pull_only:
        print("Pull-only mode, skipping analysis.")
        return

    # Phase 2: Simulate copytrade per wallet
    print(f"\n=== Phase 2: Simulating copytrade (min_size=${args.min_size}, price {args.min_price}-{args.max_price}) ===")

    results = {}
    for addr in wallets:
        trades = wallet_trades.get(addr, [])
        # filter to cutoff period
        period_trades = [t for t in trades if t.get("timestamp", 0) >= cutoff]
        closed, open_pos = simulate_copytrade(period_trades, args.min_size, args.min_price, args.max_price)
        stats = wallet_stats(closed, open_pos)
        tier, reason = tier_label(stats)
        results[addr] = {
            "stats": stats,
            "tier": tier,
            "reason": reason,
            "closed": closed,
            "open": open_pos,
            "total_trades_raw": len(period_trades),
        }

    # Phase 3: Generate report
    print(f"\n=== Phase 3: Generating report ===")

    ranked = sorted(
        [(addr, r) for addr, r in results.items() if r["stats"] is not None],
        key=lambda x: x[1]["stats"]["total_pnl"],
        reverse=True
    )
    no_data = [(addr, r) for addr, r in results.items() if r["stats"] is None]

    tier_counts = defaultdict(int)
    for _, r in results.items():
        tier_counts[r["tier"]] += 1

    total_sim_pnl = sum(r["stats"]["total_pnl"] for _, r in ranked)
    total_sim_trades = sum(r["stats"]["trades"] for _, r in ranked)
    total_sim_wins = sum(r["stats"]["wins"] for _, r in ranked)
    overall_wr = total_sim_wins / total_sim_trades * 100 if total_sim_trades > 0 else 0

    report_lines = []
    report_lines.append(f"# Copytrade Backtest Report")
    report_lines.append(f"")
    report_lines.append(f"**Generated:** {datetime.now(SGT).strftime('%Y-%m-%d %H:%M SGT')}")
    report_lines.append(f"**Period:** last {args.days} days")
    report_lines.append(f"**Filters:** min_size=${args.min_size}, price {args.min_price}-{args.max_price}")
    report_lines.append(f"**Wallets:** {len(wallets)} tracked")
    report_lines.append(f"")
    report_lines.append(f"## Summary")
    report_lines.append(f"")
    report_lines.append(f"- **Total simulated trades:** {total_sim_trades}")
    report_lines.append(f"- **Win rate:** {overall_wr:.1f}%")
    report_lines.append(f"- **Total PnL:** ${total_sim_pnl:+.2f}")
    report_lines.append(f"- **Tier distribution:** A={tier_counts['A']}, B={tier_counts['B']}, C={tier_counts['C']}, D={tier_counts['D']}")
    report_lines.append(f"")
    report_lines.append(f"## Wallet Rankings")
    report_lines.append(f"")
    report_lines.append(f"| Rank | Wallet | Tier | Trades | WR% | PnL | ROI% | Sharpe | Open |")
    report_lines.append(f"|------|--------|------|--------|-----|-----|------|--------|------|")

    for i, (addr, r) in enumerate(ranked, 1):
        s = r["stats"]
        label = f"`{addr[:6]}…{addr[-4:]}`"
        report_lines.append(
            f"| {i} | {label} | {r['tier']} | {s['trades']} | {s['win_rate']:.0f}% | ${s['total_pnl']:+.2f} | {s['roi']:+.1f}% | {s['sharpe']:.2f} | {s['open_positions']} |"
        )

    if no_data:
        report_lines.append(f"")
        report_lines.append(f"**No qualifying trades ({len(no_data)} wallets):** ")
        report_lines.append(", ".join(f"`{a[:8]}…`" for a, _ in no_data))

    # Tier breakdown
    report_lines.append(f"")
    report_lines.append(f"## Tier Definitions")
    report_lines.append(f"")
    report_lines.append(f"- **A** — ROI > 20% AND WR > 55% → follow with $20-50")
    report_lines.append(f"- **B** — ROI > 10% OR (WR > 50% AND profitable) → follow with $5-10")
    report_lines.append(f"- **C** — PnL > -$50 AND WR > 40% → follow with $2 (observe)")
    report_lines.append(f"- **D** — Unprofitable or insufficient data → keep monitoring, don't follow")
    report_lines.append(f"")
    report_lines.append(f"*Note: Tiers are conservative. D wallets are NOT eliminated — they stay in monitoring.*")

    # Top performers detail
    top_n = min(10, len(ranked))
    if top_n > 0:
        report_lines.append(f"")
        report_lines.append(f"## Top {top_n} Detailed")
        report_lines.append(f"")
        for i, (addr, r) in enumerate(ranked[:top_n], 1):
            s = r["stats"]
            report_lines.append(f"### #{i} `{addr}`")
            report_lines.append(f"- Tier: **{r['tier']}** ({r['reason']})")
            report_lines.append(f"- Trades: {s['trades']} ({s['wins']}W / {s['losses']}L)")
            report_lines.append(f"- PnL: **${s['total_pnl']:+.2f}** | ROI: {s['roi']:+.1f}%")
            report_lines.append(f"- Avg win: ${s['avg_win']:+.2f} | Avg loss: ${s['avg_loss']:+.2f}")
            report_lines.append(f"- Sharpe: {s['sharpe']:.2f}")
            report_lines.append(f"- Open positions: {s['open_positions']} (${s['open_cost']:.2f})")
            if r["closed"]:
                biggest_win = max(r["closed"], key=lambda t: t["pnl"])
                biggest_loss = min(r["closed"], key=lambda t: t["pnl"])
                report_lines.append(f"- Biggest win: ${biggest_win['pnl']:+.2f} ({biggest_win['title'][:50]})")
                report_lines.append(f"- Biggest loss: ${biggest_loss['pnl']:+.2f} ({biggest_loss['title'][:50]})")
            report_lines.append(f"- PM: https://polymarket.com/profile/{addr}")
            report_lines.append(f"")

    # Save report
    output_path = project_root / args.output
    output_path.parent.mkdir(parents=True, exist_ok=True)
    with open(output_path, "w") as f:
        f.write("\n".join(report_lines))
    print(f"\nReport saved to {output_path}")

    # Save raw results as JSON for further analysis
    json_path = project_root / "db" / "copytrade_backtest_results.json"
    summary = {}
    for addr, r in results.items():
        summary[addr] = {
            "tier": r["tier"],
            "reason": r["reason"],
            "stats": r["stats"],
            "total_trades_raw": r["total_trades_raw"],
        }
    with open(json_path, "w") as f:
        json.dump(summary, f, indent=2)
    print(f"Results JSON saved to {json_path}")

    # Print summary to stdout
    print(f"\n{'='*60}")
    print(f"COPYTRADE BACKTEST RESULTS ({args.days} days)")
    print(f"{'='*60}")
    print(f"Wallets: {len(wallets)} | With trades: {len(ranked)} | No data: {len(no_data)}")
    print(f"Total PnL: ${total_sim_pnl:+.2f} | Trades: {total_sim_trades} | WR: {overall_wr:.1f}%")
    print(f"Tiers: A={tier_counts['A']} B={tier_counts['B']} C={tier_counts['C']} D={tier_counts['D']}")
    print(f"\nTop 5:")
    for i, (addr, r) in enumerate(ranked[:5], 1):
        s = r["stats"]
        print(f"  {i}. {addr[:10]}… [{r['tier']}] {s['trades']}T {s['win_rate']:.0f}%WR ${s['total_pnl']:+.2f}")
    print(f"\nBottom 5:")
    for i, (addr, r) in enumerate(ranked[-5:], len(ranked)-4):
        s = r["stats"]
        print(f"  {i}. {addr[:10]}… [{r['tier']}] {s['trades']}T {s['win_rate']:.0f}%WR ${s['total_pnl']:+.2f}")

if __name__ == "__main__":
    main()
