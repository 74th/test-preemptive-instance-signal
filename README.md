Test Preemptive Instance

## ノードプールを0にして、新しいノードプールに移動させた時の挙動

- PodのIDは変わる
- SIGTERMが送られる
- SIGTERMが送られてから、すぐに別のノードでPodが起動する

![](./delete_node.png)

- 08:42:59 SIGTERMが送られる
- 08:43:03 次のPodが作られる
- 08:43:27 前のPodの最後のログ

## プリエンプティブインスタンスの挙動

- NodeのIDは変わらない
- PodのIDは変わらない
- SIGTERMは送られない
- いきなりKILLされてるっぽい
- 瞬断は防げない
- ノードプールに複数ノードあっても、ずらしてくれたりしない

![](./preemptive.png)

`kubectl get pods`

![](./get_pods.png)

`kubectl get nodes`

![](./get_nodes.png)
