# Extract tag name

Extract tag information from release branch

## Outputs

### `tag`

tag name

### `tagnum`

tag number

## Example usage

```yaml
- uses: vesoft-inc/.github/actions/tagname@master
  id: tag

- name: Other step
  run: |
    echo ${{ steps.tag.outputs.tag }}
    echo ${{ steps.tag.outputs.tagnum }}
```
