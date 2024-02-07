# Changelog

---

## [v0.0.0] Initial Version _2024-02-07 02:47 UTC_
### Features

- ([`c0d1dcd`](https://github.com/jakbytes/version_actions/commit/c0d1dcd0e3483390d8d7405569bcf3eadcce5710)) initial supported actions, version, sync, pull_request, extract_commit, download_release_asset

### Fixes

- ([`58bf05c`](https://github.com/jakbytes/version_actions/commit/58bf05caf571984ec6b2233ddb6f18a109a624ba)) type value needs to be output for further activity
- ([`e1729a9`](https://github.com/jakbytes/version_actions/commit/e1729a947a61a321155939e72779334c88033b47)) action trigger should be set properly
- ([`ba3d06f`](https://github.com/jakbytes/version_actions/commit/ba3d06fc58c65dc4fae5dd39c0d539207d906118)) hanging % needed to be removed from version action
- ([`19bfb4d`](https://github.com/jakbytes/version_actions/commit/19bfb4db2aa5af63bead5067d2d3582e6b67fba2)) don't use best effort
- ([`1a481d7`](https://github.com/jakbytes/version_actions/commit/1a481d72d0715ae6d7d88a9b434502513529c18c)) should be using v4 actions checkout
- ([`1487ff3`](https://github.com/jakbytes/version_actions/commit/1487ff34f740541c9cb5aa3345aa14e6d1d93abc)) commits should be freeform to allow release and others
- ([`68906c8`](https://github.com/jakbytes/version_actions/commit/68906c816d30d62c6f67c4a35b5e6003ccd74fbf)) download_release_asset shouldnt have quotes around the chmod val, version should not modify yml
- ([`42328c0`](https://github.com/jakbytes/version_actions/commit/42328c0dc7d95b59e58c1373f678834420f8c329)) actions should reference version_action, not action
- ([`8d24825`](https://github.com/jakbytes/version_actions/commit/8d24825ef39953f45c2fae275b420777c635ba5c)) a few more references to the old path were not adjusted
- ([`db31802`](https://github.com/jakbytes/version_actions/commit/db31802dc409e7306ca2a4b17a8a1ba3e8332c05)) use the download_release_asset in pull_request, rename action.go to version_action.go

### CI/CD

- ([`a48f0ae`](https://github.com/jakbytes/version_actions/commit/a48f0aeac3a5c4ce3bed5af4e055bff7174bd99f)) fix reference to type
- ([`3b55e7f`](https://github.com/jakbytes/version_actions/commit/3b55e7fbce860c789836006c2c1e93ab3a1554ce)) actions need to reference the correct path
- ([`ed5f7a3`](https://github.com/jakbytes/version_actions/commit/ed5f7a398dd060d3a9769c344206c2b86dad2959)) remove debugging action