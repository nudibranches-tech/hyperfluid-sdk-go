sync-openapi:
  echo "Syncing OpenAPI..."
  gh api /repos/nudibranches-tech/hyperfluid/contents/apis/generated/console-external.openapi.json --template '{{{{ .content }}' | base64 -d > sdk/controlplaneapiclient/control_plane_api.openapi.json
