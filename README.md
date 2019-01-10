# image-sync
sync  image tags to another registry

config are like below:

```yaml
apiVersion: v1
sync:
  from:
    registry: registry.a.com
    username: admin
    password: xxxx
  to: 
    registry: registry.b.com
    username: root
    password: yyyy
  names:
  - "framework/notify-service"
  - "framework/api-gateway"
  replace:
    - old: framework
      new: vinkdong
  rules:
  - name: release
    value: "^v?(\\d+.)*\\d+$"
```

this program can sync `registry.a.com/framework/notify-service` all tag that 
match regex `"^v?(\\d+.)*\\d+$"` to 
`registry.b.com/vinkdong/notify-service`, we can config rule very simply as config file shown.

just add additional flag `-d` to run as daemon.

```bash
./image-sync sync -c config.yml -d 
``` 