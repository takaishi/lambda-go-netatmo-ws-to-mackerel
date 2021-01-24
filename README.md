# lambda-go-netatmo-ws-to-mackerel

```
make build_container
docker tag netatmo-ws-to-mackerel:latest ${ECR_HOST}/netatmo-ws-to-mackerel:${TAG}
docker push ${ECR_HOST}/netatmo-ws-to-mackerel:${TAG}
```

# Licence

MIT Licence

original code by aereal and itchyny:

- https://github.com/aereal/gae-go-netatmo-ws-mackerel