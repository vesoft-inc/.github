# Setup OSS Environment

This action installs ossutil and configures the OSS environment for further operations.

## Inputs

### `oss-id`

**Required** OSS access key ID.

### `oss-secret`

**Required** OSS access key secret.

### `oss-endpoint`

**Required** OSS endpoint.

## Example usage

```yaml
uses: vesoft-inc/.github/actions/setup-oss@master
with:
  oss-id: ${{ secrets.OSS_ID }}
  oss-secret: ${{ secrets.OSS_SECRET }}
  oss-endpoint: ${{ secrets.OSS_ENDPOINT }}
