# Setup MinIO Action

This GitHub Action sets up the MinIO client (`mc`) and configures an alias for your MinIO server, making it easier to interact with your MinIO buckets in subsequent workflow steps.

## Inputs

- `minio_url`: The URL of your MinIO server. Required.
- `access_key`: Your MinIO access key. Required.
- `secret_key`: Your MinIO secret key. Required.

## Usage

To use this action in your workflow, add the following step:

```yaml
- name: Setup MinIO
  uses: vesoft-inc/nebula-utils/setup-minio@master
  with:
    minio_url: ${{ secrets.MINIO_ENDPOINT }}
    access_key: ${{ secrets.MINIO_KEY }}
    secret_key: ${{ secrets.MINIO_SECRET }}

- name: Copy files from MinIO
  run: mc cp minio/mybucket/myfile ./myfile

- name: Copy dir to MinIO
  run: mc cp ./mydir minio/mybucket/somedir/ 
```
