# minio
$env:MINIO_ROOT_USER = "admin"
$env:MINIO_ROOT_PASSWORD = "12345678"
D:\App\MinIO\minio.exe server D:\App\MinIO\Data --console-address ":9001"

# this is general command for mc.exe, run in cmd
D:\App\MinIO\mc.exe alias set 'myminio' 'http://ip:port' 'USER' 'PASSWORD'

# redis
redis-server
