#!/usr/bin/env bash
# shellcheck disable=SC2148,SC2155

# NOTE: If .env file doesn't exist, create a template file.
[ -f .env ] || tee .env >/dev/null <<'EOF'
# NOTE: Define environment variables that are not shared by Git.
EOF

# NOTE: Load .env files
dotenv .versenv.env
dotenv .env

# NOTE: Define environment variables that are shared by Git and not referenced in the Container here.
export REPO_ROOT=$(git rev-parse --show-toplevel || pwd || echo '.')
export PATH="${REPO_ROOT:?}/.bin:${REPO_ROOT:?}/.local/bin:${PATH:?}"
