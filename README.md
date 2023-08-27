# ccc - Cloud Cost Checker

[![license](https://img.shields.io/github/license/kunitsucom/ccc)](LICENSE)
[![pkg](https://pkg.go.dev/badge/github.com/kunitsucom/ccc)](https://pkg.go.dev/github.com/kunitsucom/ccc)
[![goreportcard](https://goreportcard.com/badge/github.com/kunitsucom/ccc)](https://goreportcard.com/report/github.com/kunitsucom/ccc)
[![workflow](https://github.com/kunitsucom/ccc/workflows/go-lint/badge.svg)](https://github.com/kunitsucom/ccc/tree/main)
[![workflow](https://github.com/kunitsucom/ccc/workflows/go-test/badge.svg)](https://github.com/kunitsucom/ccc/tree/main)
[![codecov](https://codecov.io/gh/kunitsucom/ccc/branch/main/graph/badge.svg?token=4UML9FB7BX)](https://codecov.io/gh/kunitsucom/ccc)
[![sourcegraph](https://sourcegraph.com/github.com/kunitsucom/ccc/-/badge.svg)](https://sourcegraph.com/github.com/kunitsucom/ccc)

[ccc - Cloud Cost Checker](https://github.com/kunitsucom/ccc) collects, calculates, graphs and notifies IaaS costs.  

## Project Goal

- Inform you of IaaS costs so that you are aware of sudden cost increases
- Don't bother opening the console to see costs

## Supported

### IaaS

- Google Cloud Platform

### Method of saving Cost Graph Image

- Post to Slack
- Save to local directory

## How to use

### 1. Slack

#### 1-1. Create Slack Bot and Invite channel

- Please create a Bot and issue access tokens by referring to this document.
  - [Create a bot for your workspace | Slack](https://slack.com/help/articles/115005265703)
- **Don't forget to invite Slack Bot User to your channel!**

#### 1-2. Download ccc for your execution environment

- [Releases · kunitsucom/ccc](https://github.com/kunitsucom/ccc/releases)

Download ccc for your execution environment as follows:

```bash
VERSION=v0.0.6

# download
curl -fLROSs https://github.com/kunitsucom/ccc/releases/download/${VERSION}/ccc_${VERSION}_darwin_arm64.zip

# unzip
unzip -j ccc_${VERSION}_darwin_arm64.zip '*/ccc'
```

#### 1-3. Run ccc

<!--
// secretlint-disable @secretlint/secretlint-rule-slack
-->

```bash
# Authenticate with Google User or Service Account that has permissions
# equivalent to `roles/bigquery.dataViewer` and `roles/bigquery.user` in some way.
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json

# Run ccc
./ccc \
  -tz Asia/Tokyo \
  -project your-gcp-project \
  -billing-table your-gcp-project.billing_dataset.gcp_billing_export_v1_FFFFFF_FFFFFF_FFFFFF \
  -billing-project your-gcp-project \
  -message '```your-gcp-project Cost last 30 days (last 30 days)```' \
  -slack-token xoxb-999999999999-9999999999999-ZZZZZZZZZZZZZZZZZZZZZZZZ \
  -slack-channel '#your-bot-invited-channel' \
  -days 30 \
  -debug
```

<!--
// secretlint-enable
-->

It will be posted as follows:  

[![cost](/docs/images/example.png)](/docs/images/example.png)

## If you want to post cost graphs to Slack on a regular basis

I highly recommend this GitHub Actions: [ccc-actions - GitHub Actions for Cloud Cost Checker
](https://github.com/kunitsucom/ccc-actions)  

## TODO

- IaaS
  - [x] Google Cloud Platform
  - [ ] Amazon Web Service
- Method of saving Cost Graph Image
  - [x] Post to Slack
  - [x] Save to local directory
- [x] Add tests
