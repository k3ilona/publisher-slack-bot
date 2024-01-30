# Flux bootstrap auto-image на прикладі ibot
Мета – забезпечити автоматичне розгортання образу з репозиторію ghcr.io в кластері.

Ми налаштуємо Flux Automate image updates to Git для автоматичного розгортання образу з репозиторію ghcr.io в нашому кластері. Flux слідкуватиме за репозиторієм Helm, автоматично оновлюючи випуски Helm до новіших версій, відповідно до заданих політик версіювання.

## Передумови
Перед початком роботи, переконайтеся, що у вас є:
1. Кластер Kubernetes версії 1.21 або вище. Для локального тестування можна використовувати [kind](https://kind.sigs.k8s.io) або [k3d](https://k3d.io/).
2. Обліковий запис GitHub і персональний токен доступу із дозволами на створення і видалення репозиторію.
3. Встановлений [Flux CD](https://fluxcd.io/flux/get-started/).

## Налаштування
1. Перевіряємо умови для встановлення flux:
```shell
flux check --pre
```
Отримуємо відповідь, що все нормально:
```
► checking prerequisites
✔ Kubernetes 1.27.4+k3s1 >=1.26.0-0
✔ prerequisites checks passed
```
2. Встановлюємо flux до кластеру з параметрами `--components-extra=image-reflector-controller,image-automation-controller` :
```shell
export GITHUB_TOKEN=[ваш TOKEN]
export GITHUB_OWNER=[імʼя користувача]
export GITHUB_REPO=[назва репозиторію]

flux bootstrap github --components-extra=image-reflector-controller,image-automation-controller \
--token-auth --owner=${GITHUB_OWNER} --repository=${GITHUB_REPO} \
--branch=main --path=clusters/[cluster name] --read-write-key --personal
```
Чекаємо завершення встановлення.
3. Після встановлення, перевіряємо наявність і працездатність необхідних компонентів:
```shell
k get po -A
```
Всі поди запущені і працюють:
```
NAMESPACE            NAME                                           READY   STATUS    RESTARTS   AGE
flux-system          helm-controller-865448769d-xz228               1/1     Running   0          3m21s
flux-system          image-automation-controller-69b75845c5-4ljnj   1/1     Running   0          3m21s
flux-system          image-reflector-controller-84f896568c-ngpff    1/1     Running   0          3m21s
flux-system          kustomize-controller-5c8878fd86-x4g98          1/1     Running   0          3m21s
flux-system          notification-controller-59696fbb58-v8qxl       1/1     Running   0          3m21s
flux-system          source-controller-fc5555fb-2nw8x               1/1     Running   0          3m21s
```
4. Створимо розгортання для бота першої версії в кластері `clusters/dev`:
```yaml
apiVersion: v1
kind: Namespace
metadata:
name: ibot-dev
```
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ibot
  namespace: ibot-dev
spec:
  selector:
    matchLabels:
      app: ibot
  template:
    metadata:
      labels:
        app: ibot
    spec:
      containers:
        - name: ibot
          image: ghcr.io/k3ilona/ibot:0.0.1-dev-533f12b
          imagePullPolicy: IfNotPresent
          env:
            - name: SLACK_APP_TOKEN
              valueFrom:
                secretKeyRef:
                  name: ibot
                  key: appToken
            - name: SLACK_AUTH_TOKEN
              valueFrom:
                secretKeyRef:
                  name: ibot
                  key: token
            - name: SLACK_CHANNEL_ID
              valueFrom:
                secretKeyRef:
                  name: ibot
                  key: channel-id
          resources:
            limits:
              cpu: "0.5"
              memory: "512Mi"
            requests:
              cpu: "0.1"
              memory: "256Mi"
```
5. Перевіримо інформацію про деплоймент:
```shell
k get deployment/ibot -n ibot-dev -o yaml | grep 'image:'
```
або 
```shell
k describe deploy -n ibot-dev ibot | grep 'Image:'
```
Якщо отримали відповідь з назвою нашого образу, то розгортання пройшло успішно:
```
Image:        ghcr.io/k3ilona/ibot:0.0.1-dev-533f12b
```
6. Створюємо `ImageRepository`, щоб вказати Flux де шукати оновлення образів і з яким інтервалом:
```shell
flux create image repository ibot \
--image=ghcr.io/k3ilona/ibot \
--interval=3m \
--export > dev-repo.yaml
```
Буде згенерований наступний маніфест:
```yaml
---
apiVersion: image.toolkit.fluxcd.io/v1beta2
kind: ImageRepository
metadata:
  name: ibot
  namespace: flux-system
spec:
  image: ghcr.io/k3ilona/ibot
  interval: 3m0s
```
7. Для приватних зображень ви можете створити секрет Kubernetes у тому самому просторі імен, що й `ImageRepository`with `kubectl create secret ghcr-registry`. Тоді ви можете налаштувати Flux на використання облікових даних, посилаючись на секрет Kubernetes у `ImageRepository`:
```yaml
kind: ImageRepository
spec:
  secretRef:
    name: regcred
```
8. Створюємо `ImagePolicy`, щоб вказати Flux, яку політику оновлення версій використовувати:
```shell
flux create image policy ibot \
--image-ref=ibot \
--select-semver=1.x \
--export > dev-policy.yaml
```
Наведена вище команда створює такий маніфест:
```yaml
---
apiVersion: image.toolkit.fluxcd.io/v1beta2
kind: ImagePolicy
metadata:
  name: ibot
  namespace: flux-system
spec:
  imageRepositoryRef:
    name: ibot
  policy:
    semver:
    range: '>=1.0.0'
```
При потребі, можна додати фільтри або правила для оновлення версій: https://fluxcd.io/flux/guides/image-update/#imagepolicy-examples
Приклад використання:
```yaml
---
apiVersion: image.toolkit.fluxcd.io/v1beta2
kind: ImagePolicy
metadata:
  name: ibot
  namespace: flux-system
spec:
  imageRepositoryRef:
    name: ibot
  filterTags:
    pattern: '^\d\.\d\.\d-dev-[a-f0-9]+-(?P<ts>\d+)'
    extract: '$ts'
  policy:
    numerical:
      order: asc
```
9. Зафіксуйте та надішліть зміни до головної гілки та зачекайте, поки Flux отримає список тегів зображень із реєстру контейнерів. Потім можна перевірити список доступних тегів:
```shell
flux get image repository ibot
```
Отримаємо коротку відповідь:
```
NAME	LAST SCAN           	SUSPENDED	READY	MESSAGE                        
ibot	2024-01-30T09:28:15Z	False    	True 	successful scan: found 37 tags	
```
або
```shell
k -n flux-system describe imagerepositories ibot
```
Отримаємо розгорнуту відповідь:
```
Name:         ibot
Namespace:    flux-system
Labels:       kustomize.toolkit.fluxcd.io/name=flux-system
              kustomize.toolkit.fluxcd.io/namespace=flux-system
Annotations:  reconcile.fluxcd.io/requestedAt: 2024-01-29T17:22:23.186705496Z
API Version:  image.toolkit.fluxcd.io/v1beta2
Kind:         ImageRepository
Metadata:
  Creation Timestamp:  2024-01-28T12:08:43Z
  Finalizers:
    finalizers.fluxcd.io
  Generation:        1
  Resource Version:  264445
  UID:               4b58eb1a-95cd-42bc-86d9-fdf6241aec17
Spec:
  Exclusion List:
    ^.*\.sig$
  Image:     ghcr.io/k3ilona/ibot
  Interval:  3m0s
  Provider:  generic
Status:
  Canonical Image Name:  ghcr.io/k3ilona/ibot
  Conditions:
    Last Transition Time:     2024-01-29T22:23:14Z
    Message:                  successful scan: found 37 tags
    Observed Generation:      1
    Reason:                   Succeeded
    Status:                   True
    Type:                     Ready
  Last Handled Reconcile At:  2024-01-29T17:22:23.186705496Z
  Last Scan Result:
    Latest Tags:
      0.1.4-dev-f5301ad-1706562879
      0.1.4-dev-f522cfd
      0.1.4-dev-d4e1d8e-1706561827
      0.1.4-dev-c2a06db
      0.1.4-dev-b3f2699
      0.1.4-dev-9e8aa04
      0.1.4-dev-64a1422-1706563823
      0.1.4-dev-49e40b9-1706564586
      0.1.4-dev-4550524
      0.1.3-dev-c97f46b
    Scan Time:  2024-01-30T09:28:15Z
    Tag Count:  37
  Observed Exclusion List:
    ^.*\.sig$
  Observed Generation:  1
Events:
  Type    Reason     Age                 From                        Message
  ----    ------     ----                ----                        -------
  Normal  Succeeded  6m (x860 over 45h)  image-reflector-controller  no new tags found, next scan in 3m0s
```
10.  Налаштуємо оновлення образів. Для цього оновлюємо файл розгортання нашого боту `dev-deploy.yaml`та додаєм маркер для Flux, який додає політику оновлення image в контейнері:
```yaml
spec:
  containers:
  - name: podinfod
    image: ghcr.io/k3ilona/ibot:0.0.1-dev-533f12b # {"$imagepolicy": "flux-system:ibot"}
```
11. Створюємо `ImageUpdateAutomation`, щоб вказати Flux, до якого репозиторію Git (повинен містити маніфест розгортання) і з яким інтервалом записувати оновлення зображень:
```shell
flux create image update ibot \                 
--interval=3m \
--git-repo-ref=flux-system \
--git-repo-path="./clusters/dev/ibot" \
--checkout-branch=main \
--push-branch=main \
--author-name=fluxcdbot \
--author-email=fluxcdbot@users.noreply.github.com \
--commit-template="{{range .Updated.Images}}{{println .}}{{end}}" \
--export > dev-automatiom.yaml
```
Наведена вище команда створює такий маніфест:
```yaml
---
apiVersion: image.toolkit.fluxcd.io/v1beta1
kind: ImageUpdateAutomation
metadata:
  name: ibot
  namespace: flux-system
spec:
  git:
    checkout:
      ref:
        branch: main
    commit:
      author:
        email: fluxcdbot@users.noreply.github.com
        name: fluxcdbot
      messageTemplate: '{{range .Updated.Images}}{{println .}}{{end}}'
    push:
      branch: main
  interval: 3m0s
  sourceRef:
    kind: GitRepository
    name: flux-system
  update:
    path: ./clusters/dev/ibot
    strategy: Setters
```
12. Перевіряємо що все працює правильно:
```shell
flux get image repository kbot
flux get image policy kbot
flux get image update kbot
```
або
```shell
flux get images all --all-namespaces
```
Якщо все вірно налаштовано, Flux підтягне оновлену версію образу з репозиторію і оновить її в кластері:
```shell
NAMESPACE  	NAME                	LAST SCAN           	SUSPENDED	READY	MESSAGE                        
flux-system	imagerepository/ibot	2024-01-30T09:37:17Z	False    	True 	successful scan: found 37 tags	

NAMESPACE  	NAME            	LATEST IMAGE                                     	READY	MESSAGE                                                                                                               
flux-system	imagepolicy/ibot	ghcr.io/k3ilona/ibot:0.1.4-dev-49e40b9-1706564586	True 	Latest image tag for 'ghcr.io/k3ilona/ibot' updated from 0.1.3-dev-51d4a0e-1706565110 to 0.1.4-dev-49e40b9-1706564586	

NAMESPACE  	NAME                      	LAST RUN            	SUSPENDED	READY	MESSAGE                                                      
flux-system	imageupdateautomation/ibot	2024-01-30T09:37:18Z	False    	True 	no updates made; last commit eecc501 at 2024-01-29T22:05:11Z
```
