# Changelog

## [v1.2.2](https://github.com/jakbytes/version_actions/compare/v1.2.1...v1.2.2) (2024-02-20)
### Fixes

- ([`e938309`](https://github.com/jakbytes/version_actions/commit/e938309d5ca33b6241ae835dd31d403e13e832a7)) version-action should be version_action

## [v1.2.1](https://github.com/jakbytes/version_actions/compare/v1.2.0...v1.2.1) (2024-02-20)
### Fixes

- ([`07cbc39`](https://github.com/jakbytes/version_actions/commit/07cbc3960f481b17afc9085cd1c3faf537b5d6ab)) best effort

## [v1.2.0](https://github.com/jakbytes/version_actions/compare/v1.1.0...v1.2.0) (2024-02-16)
### Features

- ([`df1846b`](https://github.com/jakbytes/version_actions/commit/df1846bc69b656c36e0e5c45c9e2dced809d336c)) chores should trigger patch version update if they're deps related

## [v1.1.0](https://github.com/jakbytes/version_actions/compare/v1.0.3...v1.1.0) (2024-02-16)
### Features

- ([`f787f07`](https://github.com/jakbytes/version_actions/commit/f787f0781aea855bfd46b686a860678d23338298)) add chore to changelog

### Fixes

- ([`0abefd7`](https://github.com/jakbytes/version_actions/commit/0abefd7355a6db1bd69971ae9e057654f822d7ca)) if PR is removed because of force reset branch it shouldn't fail but instead create a new branch
- ([`199b999`](https://github.com/jakbytes/version_actions/commit/199b999fb7286eab3296e88fca14205d031090f5)) monitor flow
- ([`9498bcc`](https://github.com/jakbytes/version_actions/commit/9498bcce2657d3c8cfde89878b3aafe8ca7d3a07)) sync should ff if available
- ([`636227a`](https://github.com/jakbytes/version_actions/commit/636227af427e8ed0f0635572730768b1cb1602b9)) sync should ff if available
- ([`054a002`](https://github.com/jakbytes/version_actions/commit/054a0028b5df9de13b216a8913abb450221d4e68)) remove newline
- ([`a3c968c`](https://github.com/jakbytes/version_actions/commit/a3c968c7fd4da7633983d853b0142e81f2ebe4dc)) sync should ff if available, commit message if not available

## [v1.0.2](https://github.com/jakbytes/version_actions/compare/v1.0.1...v1.0.2) (2024-02-09)
### Fixes

- ([`0dc65ac`](https://github.com/jakbytes/version_actions/commit/0dc65ac15feb19a905bcef1cf1099362e2e0c7ad)) shouldn't stop PR if version increment is -1 and this is first version

## [v1.0.1](https://github.com/jakbytes/version_actions/compare/v1.0.0...v1.0.1) (2024-02-08)
### Fixes

- ([`362c4ce`](https://github.com/jakbytes/version_actions/commit/362c4ce03fc643e689acb2bccfb1be8ed46149f7)) output was not set properly for identifier

## [v1.0.0](https://github.com/jakbytes/version_actions/compare/v0.1.1...v1.0.0) (2024-02-08)
### ⚠ BREAKING CHANGES

- ([`c5d049e`](https://github.com/jakbytes/version_actions/commit/c5d049ee3e5fc513aace47b03ead2286671feace)) require passed token, updated various actions to account for that
  > 
  > BREAKING CHANGE: tokens are required in some actions where they were not before

### Features

- ([`1f982dc`](https://github.com/jakbytes/version_actions/commit/1f982dcc79afd547a03617ee6fee43bb22b4b27b)) support passing in token for private repositories

### Fixes

- ([`8f34822`](https://github.com/jakbytes/version_actions/commit/8f348224905504dc70598287f53e5ca207dd88a2)) token download not working properly
- ([`6b99a7c`](https://github.com/jakbytes/version_actions/commit/6b99a7c5b40138c8cb28f9190a0187199cf6501f)) use token not github_token

## [v0.1.1](https://github.com/jakbytes/version_actions/compare/v0.1.0...v0.1.1) (2024-02-07)
### Fixes

- ([`efc4e6e`](https://github.com/jakbytes/version_actions/commit/efc4e6ee9c1115e471a8653b6306b91404cac2b3)) correct the module name for imports

### CI/CD

- ([`176d7c3`](https://github.com/jakbytes/version_actions/commit/176d7c347e2632069d3cd84fa848880de86e78a9)) add testing to pull requests for better code coverage

## [v0.1.0](https://github.com/jakbytes/version_actions/compare/v0.0.1...v0.1.0) (2024-02-07)
### Features

- ([`bd3c539`](https://github.com/jakbytes/version_actions/commit/bd3c539485ed2326b9f53b05ed2bccba9989aae5)) support debug types in conventional commit syntax

### Fixes

- ([`2b1b493`](https://github.com/jakbytes/version_actions/commit/2b1b49317a8c94f1bb411fdef538e524c81986ef)) release action needs to use the correct file for the release asset

### Debugging

- ([`ff81697`](https://github.com/jakbytes/version_actions/commit/ff81697d995cf560e47f030abff969c92b01a50c)) occasionally the list of commits does not stop at the correct hash, logging behavior

## [v0.0.1](https://github.com/jakbytes/version_actions/compare/v0.0.0...v0.0.1) (2024-02-07)
### Fixes

- ([`ffa1d53`](https://github.com/jakbytes/version_actions/commit/ffa1d5370ba7bfd933ace9954b0fd369021a9665)) no increment should not generate a PR

### CI/CD

- ([`2c3acd9`](https://github.com/jakbytes/version_actions/commit/2c3acd9472dc067d35a29990f122a71e4fad0372)) softprops/action-gh-release pinned to sha with node20

## [v0.0.0] Initial Version (2024-02-07)
### Features

- ([`3fb5621`](https://github.com/jakbytes/version_actions/commit/3fb562193137e64068da04b1dbb9d3c69b4fc5b3)) initial actions, version, sync, pr, extract_commit, download_release_asset, prerelease
- ([`c2a4629`](https://github.com/jakbytes/version_actions/commit/c2a4629dd8aadaafd8b577cf3738a8ec4eb34624)) clean up prs markdown
- ([`3154777`](https://github.com/jakbytes/version_actions/commit/3154777b22d84f248a31abf98695727df4d84b8e)) body and footer of the commit will end up in the changelog as a subsection
- ([`38f1bd1`](https://github.com/jakbytes/version_actions/commit/38f1bd1091e162416bbcc653da5865b8f70e2c49)) breaking changes text capitalized to call it out strongly:
- ([`7237226`](https://github.com/jakbytes/version_actions/commit/72372265d197605918b127c92eb75375c3715382)) date on version is simplified
- ([`0ba489f`](https://github.com/jakbytes/version_actions/commit/0ba489f5f33d221061c149fed64166c26c6322ae)) extract prerelease identifier action
- ([`c0d1dcd`](https://github.com/jakbytes/version_actions/commit/c0d1dcd0e3483390d8d7405569bcf3eadcce5710)) initial supported actions, version, sync, pull_request, extract_commit, download_release_asset

### Fixes

- ([`deb7523`](https://github.com/jakbytes/version_actions/commit/deb7523fc729ed0e9a1ef8a0a05710af9a783841)) version should have v prefix
- ([`0355658`](https://github.com/jakbytes/version_actions/commit/03556582d5a46e64452c945454d95e2ddc1a4784)) observing
- ([`93297c1`](https://github.com/jakbytes/version_actions/commit/93297c169c9ce8aabf9f0df7292e2e04a6296070)) observing
- ([`70d030b`](https://github.com/jakbytes/version_actions/commit/70d030b01e8d9672076b8017cd10d6d75b001986)) multiline commit handling fixed
- ([`d7ce0b8`](https://github.com/jakbytes/version_actions/commit/d7ce0b88ef4d3c296f7db91d6ec5c14af2233c2b)) resolve issue where existing changelog versions were dropped
  > 
  > The conditional logic for skiping existing lines did not correctly identify when to stop skipping lines for the next version.
- ([`24b21c0`](https://github.com/jakbytes/version_actions/commit/24b21c024993d337061c3ee53f3d179f11293ecb)) prerelease_identifier action input should be version and optional
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

- ([`09744c1`](https://github.com/jakbytes/version_actions/commit/09744c1d845d1c9c26d9831595a27c26f4bacc38)) tweak pull_request to properly handle dev branches
- ([`d9a2852`](https://github.com/jakbytes/version_actions/commit/d9a28521ed93dcac2c43df8137b89eba5e231be2)) action-gh-release version, node 16 to node 20
- ([`a48f0ae`](https://github.com/jakbytes/version_actions/commit/a48f0aeac3a5c4ce3bed5af4e055bff7174bd99f)) fix reference to type
- ([`3b55e7f`](https://github.com/jakbytes/version_actions/commit/3b55e7fbce860c789836006c2c1e93ab3a1554ce)) actions need to reference the correct path
- ([`ed5f7a3`](https://github.com/jakbytes/version_actions/commit/ed5f7a398dd060d3a9769c344206c2b86dad2959)) remove debugging action