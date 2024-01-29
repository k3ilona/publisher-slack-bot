# Розгортання контролера Flux2 для відстеження бажаного стану кластерів

Розглянемо налаштування розгортання чотирьох кластерів: `staging`, `production`, `dev`, `qa`, із метою оптимізації управління через Flux та Kustomize.

Мета – забезпечити ефективне управління без непотрібного повторення налаштувань між кластерами, використовуючи ці інструменти для автоматизації та синхронізації конфігурацій.

Ми налаштуємо Flux для автоматичного керування демонстраційним додатком, використовуючи HelmRepository та HelmRelease. Це забезпечить встановлення, тестування та оновлення додатка. Flux слідкуватиме за репозиторієм Helm, автоматично оновлюючи випуски Helm до новіших версій, відповідно до заданих діапазонів семантичного версіювання (semver).

На кожному з кластерів буде встановлено відповідний артефакт застосунку, відкритий інтерфейс для Flux, для візуалізації та моніторингу робочих навантажень, керованих Flux. Цей приклад демонструє застосування принципів GitOps в керуванні бажаним станом кластерів та застосунків в них.

## Передумови

Перед початком роботи, переконайтеся, що у вас є:

1. Кластер Kubernetes версії 1.21 або новішої. Для локального тестування можна використовувати [k3d](https://k3d.io/stable/).
2. Обліковий запис GitHub і [персональний токен доступу](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) із відповідними дозволами.

## Встановлення Flux та k3d

1. **Встановлення Flux CLI**

   - MacOS або Linux через Homebrew:

     ```sh
     brew install fluxcd/tap/flux
     ```

   - Завантаження скомпільованих бінарників через Bash-скрипт:

     ```sh
     curl -s https://fluxcd.io/install.sh | sudo bash
     ```

2. **Встановлення k3d**:

   - Встановлення останньої версії k3d:

     ```sh
     curl -sfL https://get.k3s.io | sh -
     ```

   - Перевірка встановлення:

     ```sh
     k3d --version
     ```

## Структура репозиторію та Bootstrap кластерів

Структура репозиторію виглядає наступним чином:

```sh
./clusters/
├── production/
│   ├── apps.yaml
│   └── infrastructure.yaml
├── staging/
│   ├── apps.yaml
│   └── infrastructure.yaml
├── dev/
│   ├── apps.yaml
│   └── infrastructure.yaml
└── qa/
    ├── apps.yaml
    └── infrastructure.yaml
```

Кожен кластер має свої конфігурації Flux `Kustomization`. Наприклад, у `clusters/staging/` ми маємо таку конфігурацію:

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: apps
  namespace: flux-system
spec:
  interval: 10m0s
  dependsOn:
    - name: infra-configs
  sourceRef:
    kind: GitRepository
    name: flux-system
  path: ./apps/staging
  prune: true
  wait: true
```

## Ініціалізація та налаштування кластерів

Виконайте наступні кроки для кожного кластера (`staging`, `production`, `dev`, `qa`):

1. **Створення кластера через k3d**:

   Виконайте наступну команду для створення кластера:

   ```sh
   k3d cluster create <назва кластера>
   ```

2. **Встановлення контексту kubectl**:

   Експортуйте вміст файлу `kubeconfig` у змінну середовища `KUBECONFIG`:
  
   ```sh
   export KUBECONFIG=$(k3d kubeconfig write <назва кластера>) 
   ```

3. **Ініціалізація Flux**:

   - Форкніть [репозиторій](https://github.com/fluxcd/flux2-kustomize-helm-example) до вашого облікового запису GitHub.
   - Експортуйте токен доступу до GitHub, ім'я користувача та ім'я репозиторію:

     ```sh
     export GITHUB_TOKEN=<ваш токен>
     export GITHUB_USER=<ваше імʼя користувача>
     export GITHUB_REPO=<назва репозиторію>
     ```

   - Перевірте необхідні умови:

     ```sh
     flux check --pre
     ```

   - Ініціалізуйте Flux:

     ```sh
     flux bootstrap github \
         --context=k3d-<назва кластера> \
         --owner=${GITHUB_USER} \
         --repository=${GITHUB_REPO} \
         --branch=main \
         --path=clusters/<назва кластера>
     ```

4. **Перевірка роботи**:

   - Використовуйте `port-forward` для доступу до додатків через локальний порт:

     ```sh
     kubectl -n ingress-nginx port-forward svc/ingress-nginx-controller 8080:80 &
     ```

   - Тестування через cURL:

     ```sh
     curl -H "Host: podinfo.<назва кластера>" http://localhost:8080
     ```

5. **Перевірка випусків Helm**:

   ```sh
   flux get helmreleases --all-namespaces 
   ```

Повторіть ці кроки для кожного кластера, змінюючи відповідні параметри (`<назва кластера>`, `<назва репозиторію>` тощо). Це забезпечить автоматичне управління додатками та інфраструктурою через Flux, Kustomize та Helm, оптимізуючи процес управління конфігураціями.

З результатом можете ознайомитися тут – <https://github.com/k3ilona/multicluster>

---
 ← [Повернутись до змісту](../README.md)  
