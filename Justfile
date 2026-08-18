# The contents API only inlines `content` for blobs under 1 MiB; above that it
# answers `encoding: "none"` with an empty string, and `base64 -d` of nothing is
# a silent 0-byte spec that regenerates an empty client. The raw media type has
# no size limit, and the JSON check below makes any future failure loud.
sync-openapi branch="main":
  echo "Syncing OpenAPI..."
  gh api "/repos/nudibranches-tech/hyperfluid/contents/apis/generated/console-external.openapi.json?ref={{branch}}" \
    -H "Accept: application/vnd.github.raw" > sdk/controlplaneapiclient/control_plane_api.openapi.json
  python3 -c "import json,sys; d=json.load(open('sdk/controlplaneapiclient/control_plane_api.openapi.json')); sys.exit('fetched spec has no schemas') if not d.get('components',{}).get('schemas') else print('spec ok:', len(d['components']['schemas']), 'schemas')"
  go generate ./...
