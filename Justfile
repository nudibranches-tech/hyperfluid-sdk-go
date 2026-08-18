sync-openapi branch="main":
  #!/usr/bin/env bash
  set -euo pipefail
  # The contents API only inlines `content` for blobs under 1 MiB. Above that it
  # answers `encoding: "none"` with an empty string, so base64-decoding it wrote a
  # 0-byte spec and regenerated an empty client without failing. The raw media type
  # has no size limit.
  #
  # Fetch to a temp file and validate before replacing the committed spec: a plain
  # `>` redirect truncates the good spec before gh runs, so a failed fetch used to
  # leave nothing behind.
  dest=sdk/controlplaneapiclient/control_plane_api.openapi.json
  tmp=$(mktemp)
  trap 'rm -f "$tmp"' EXIT
  echo "Syncing OpenAPI from {{branch}}..."
  gh api "/repos/nudibranches-tech/hyperfluid/contents/apis/generated/console-external.openapi.json?ref={{branch}}" \
    -H "Accept: application/vnd.github.raw" > "$tmp"
  jq -e '(.components.schemas | length > 0) and (.paths | length > 0)' "$tmp" > /dev/null \
    || { echo "fetched spec is empty or malformed; keeping the committed one" >&2; exit 1; }
  echo "spec ok: $(jq '.components.schemas | length' "$tmp") schemas, $(jq '.paths | length' "$tmp") paths"
  # The raw media type returns no trailing newline; end-of-file-fixer requires one.
  [ -n "$(tail -c1 "$tmp")" ] && printf '\n' >> "$tmp"
  mv "$tmp" "$dest"
  trap - EXIT
  go generate ./...
