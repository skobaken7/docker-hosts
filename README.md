# docker-hosts

Docker containers use port mapping to export container port to host. When some containers use the same port, We select each port carefully.
But, by using container ip addresses, we need not use port mapping because it enables to access container port directly.
docker-hosts modifys /etc/hosts so that we can use container ip addresses easily.

## How to use

``` sh
docker run -d \
  --name docker-hosts \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /etc/hosts:/etc/hosts \
  ghcr.io/skobaken7/docker-hosts
```

Now, you can connect to containers by `${container_alias}.${network_name}.docker.internal`

### container_alias?

The value of `docker inspect ${container_id} | jq -r ".[0].NetworkSettings.Networks.${network_name}.Aliases[0]" `.
This is service name if the container is created by docker-compose.
