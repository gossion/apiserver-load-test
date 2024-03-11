1. build
    ```
    cd app-http
    make docker-build
    docker push $img
    ```

1. deploy
    ```
    k apply -f deploy/example-pod.yml
    k apply -f deploy/sa.yml
    k apply -f deploy/ds-alternative.yml
    ```

1. get logs
    ```
    k logs -l app=apiserver-load-tester --tail 2
    ```

1. get established connections
    * can't work when nat port exhausted, we can get established connections by `ss` on the nodes.
    * app has sleep interval, we don't know if the connections are still truly alive, `ss` is more acurate.
    ```
    bash scripts/get-connetions.sh
    ```

