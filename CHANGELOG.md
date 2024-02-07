# Changelog

---

## [v0.0.0] Initial Version _2024-02-07 01:28 UTC_
### Features

- ([`c0d1dcd`](https://github.com/jakbytes/version_actions/commit/c0d1dcd0e3483390d8d7405569bcf3eadcce5710)) initial supported actions, version, sync, pull_request, extract_commit, download_release_asset

### Fixes

- ([`68906c8`](https://github.com/jakbytes/version_actions/commit/68906c816d30d62c6f67c4a35b5e6003ccd74fbf)) download_release_asset shouldnt have quotes around the chmod val, version should not modify yml
- ([`42328c0`](https://github.com/jakbytes/version_actions/commit/42328c0dc7d95b59e58c1373f678834420f8c329)) actions should reference version_action, not action
- ([`8d24825`](https://github.com/jakbytes/version_actions/commit/8d24825ef39953f45c2fae275b420777c635ba5c)) a few more references to the old path were not adjusted
- ([`db31802`](https://github.com/jakbytes/version_actions/commit/db31802dc409e7306ca2a4b17a8a1ba3e8332c05)) use the download_release_asset in pull_request, rename action.go to version_action.go

### CI/CD

- ([`3b55e7f`](https://github.com/jakbytes/version_actions/commit/3b55e7fbce860c789836006c2c1e93ab3a1554ce)) actions need to reference the correct path
- ([`ed5f7a3`](https://github.com/jakbytes/version_actions/commit/ed5f7a398dd060d3a9769c344206c2b86dad2959)) remove debugging action
