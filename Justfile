sync-openapi branch="main":
  echo "Syncing OpenAPI..."
  gh api "/repos/nudibranches-tech/hyperfluid/contents/apis/generated/console-external.openapi.json?ref={{branch}}" --template '{{{{ .content }}' | base64 -d > sdk/controlplaneapiclient/control_plane_api.openapi.json
  go generate ./...
