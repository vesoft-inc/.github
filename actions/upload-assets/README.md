# Upload file to release assets

Upload file to release assets

## Inputs

### `tag`

**Required** tag name of release branch. Default `${{ github.ref }}`.

### `asset-path`

**Required** file path to be uploaded. Default `''`.

## Example usage

```yaml
uses: vesoft-inc/.github/actions/upload-assets@master
with:
  asset-path: build/cpack_output
  tag: ${{ steps.tag.outputs.tag  }}
```
