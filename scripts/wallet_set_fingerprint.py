#!/usr/bin/env python3
import argparse
import hashlib
import re


ADDRESS_RE = re.compile(r"(?i)^0x[0-9a-f]{40}$")
LIST_RE = re.compile(r"(?:^|\s)list=([^\s#]+)", re.IGNORECASE)


def read_entries(paths):
    for path in paths:
        try:
            with open(path, encoding="utf-8") as handle:
                for raw in handle:
                    line = raw.strip()
                    if not line or line.startswith("#"):
                        continue
                    body, _, comment = line.partition("#")
                    fields = body.split()
                    if not fields or not ADDRESS_RE.fullmatch(fields[0]):
                        continue
                    address = fields[0].lower()
                    list_match = LIST_RE.search(comment)
                    list_name = list_match.group(1).lower() if list_match else ""
                    yield address, list_name
        except FileNotFoundError:
            continue


def effective_entries(sources, excludes, mode):
    blocked = {address for address, _ in read_entries(excludes)}
    entries = {}
    for address, list_name in read_entries(sources):
        if address in blocked or address in entries:
            continue
        entries[address] = list_name if mode == "address-list" else ""
    return sorted(entries.items())


def fingerprint(sources, excludes, mode):
    body = "".join(f"{address}|{list_name}\n" for address, list_name in effective_entries(sources, excludes, mode))
    return hashlib.sha256(body.encode("ascii")).hexdigest()


def main():
    parser = argparse.ArgumentParser(description="Fingerprint an effective wallet set")
    parser.add_argument("--mode", choices=("address", "address-list"), default="address")
    parser.add_argument("--source", action="append", default=[])
    parser.add_argument("--exclude", action="append", default=[])
    args = parser.parse_args()
    print(fingerprint(args.source, args.exclude, args.mode))


if __name__ == "__main__":
    main()
