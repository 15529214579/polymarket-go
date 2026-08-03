import os
import tempfile
import unittest

import wallet_set_fingerprint as subject


class WalletSetFingerprintTest(unittest.TestCase):
    def write(self, directory, name, body):
        path = os.path.join(directory, name)
        with open(path, "w", encoding="utf-8") as handle:
            handle.write(body)
        return path

    def test_ignores_metadata_and_order_for_address_mode(self):
        with tempfile.TemporaryDirectory() as directory:
            first = self.write(directory, "first", "0x" + "1" * 40 + " # score=1\n0x" + "2" * 40 + "\n")
            second = self.write(directory, "second", "0x" + "2" * 40 + " # score=99\n0x" + "1" * 40 + " # changed\n")
            self.assertEqual(subject.fingerprint([first], [], "address"), subject.fingerprint([second], [], "address"))

    def test_address_list_mode_tracks_effective_first_list_and_excludes(self):
        with tempfile.TemporaryDirectory() as directory:
            address = "0x" + "1" * 40
            source_a = self.write(directory, "source-a", address + " # list=core score=1\n")
            source_b = self.write(directory, "source-b", address + " # list=watch score=1\n")
            blocked = self.write(directory, "blocked", address + "\n")
            core = subject.fingerprint([source_a], [], "address-list")
            watch = subject.fingerprint([source_b], [], "address-list")
            self.assertNotEqual(core, watch)
            self.assertEqual(subject.effective_entries([source_a, source_b], [], "address-list"), [(address, "core")])
            self.assertEqual(subject.effective_entries([source_a], [blocked], "address-list"), [])

    def test_ignores_addresses_in_comments(self):
        with tempfile.TemporaryDirectory() as directory:
            address = "0x" + "1" * 40
            path = self.write(directory, "comments", "# archived " + address + "\n")
            self.assertEqual(subject.effective_entries([path], [], "address"), [])


if __name__ == "__main__":
    unittest.main()
