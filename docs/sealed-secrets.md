# Налаштування Secrets у Flux за допомогою Sealed Secrets

## 1. Встановлюємо kubeseal – Sealed Secrets CLI

```shell
KUBESEAL_VERSION='0.25.0' # Set this to, for example, KUBESEAL_VERSION='0.23.0'
wget "https://github.com/bitnami-labs/sealed-secrets/releases/download/v${KUBESEAL_VERSION:?}/kubeseal-${KUBESEAL_VERSION:?}-linux-amd64.tar.gz"
tar -xvzf kubeseal-${KUBESEAL_VERSION:?}-linux-amd64.tar.gz kubeseal
sudo install -m 755 kubeseal /usr/local/bin/kubeseal
```

## 2. Створюємо HelmRepository для Sealed Secrets

```shell
flux create source helm sealed-secrets \
  --url https://bitnami-labs.github.io/sealed-secrets \
  --export > sealed-secrets-repo.yaml
```

Отримуємо маніфест:

```yaml
---
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: HelmRepository
metadata:
  name: sealed-secrets
  namespace: flux-system
spec:
  interval: 1m0s
  url: https://bitnami-labs.github.io/sealed-secrets
```

## 3. Створюємо HelmRelease для Sealed Secrets

```shell
flux create helmrelease sealed-secrets \
  --chart sealed-secrets \
  --source HelmRepository/sealed-secrets \
  --target-namespace flux-system \
  --release-name sealed-secrets-controller \
  --crds CreateReplace \
  --chart-version ">=1.15.0-0" \
  --export > sealed-secrets-helmrelease.yaml
```

Отримуємо маніфест:

```yaml
---
apiVersion: helm.toolkit.fluxcd.io/v2beta2
kind: HelmRelease
metadata:
  name: sealed-secrets
  namespace: flux-system
spec:
  chart:
    spec:
      chart: sealed-secrets
      reconcileStrategy: ChartVersion
      sourceRef:
        kind: HelmRepository
        name: sealed-secrets
      version: '>=1.15.0-0'
  install:
    crds: Create
  interval: 1m0s
  releaseName: sealed-secrets-controller
  targetNamespace: flux-system
  upgrade:
    crds: CreateReplace
```

## 4. Розгортаємо маніфести у Flux і перевіряємо працездатність розгортання

```shell
kubectl get po -A
```

## 5. Створюємо публічний ключ для шифрування секрету

```shell
kubeseal --fetch-cert \
  --controller-name=sealed-secrets-controller \
  --controller-namespace=flux-system \
  > pub-sealed-secrets.pem
```

## 6. Створюємо Secret маніфест (на прикладі боту і токена для нього)

```shell
read -s SECRET_TOKEN
```

```shell
kubectl -n app_namespace create secret generic app_name \
--dry-run=client \
--from-literal=token=$SECRET_TOKEN \
-o yaml > secret.yaml
```

Отримуємо маніфест з токеном:

```yaml
apiVersion: v1
data:
  token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx==
kind: Secret
metadata:
  creationTimestamp: null
  name: app_name
  namespace: app_namespace
```

## 7. Шифруємо наш токен за допомогою публічного ключа

```shell
kubeseal --format=yaml \
  --cert=pub-sealed-secrets.pem \
  < secret.yaml > secret-sealed.yaml
rm secret.yaml
```

Отримуємо маніфест з зашифрованим токеном:

```yaml
---
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: app_name
  namespace: app_namespace
spec:
  encryptedData:
    token: AgCRaiMVEYlfHmvMGNkqiLGTf7bRBV1iFzpjrA5GvXWT0EKgVxMORcN6XMmsKChdw8q6af6D/SylMzGuBHTCuvrD7ONfUJQlL4xXubjyqHUU+Vkz369BGro7xLCb0yCjH8yQOXibR4K5RWo45igx83eFCgNEUvEh2xjfiE467PZ9caYO/8QFn4Rygotaryqa5OUM+iv1hMSOFCyP2Aa/9RrRijyjFXuJ+uz5QOzMqNAH2VfVQDkfme7WubmiE6Zp1t0t9C3lvdGLvUv1YCkb3GOsB9LgSV485EUURAC+GFAbx/+SkvITP0Iqv5Zp7LzpffVlvj5ZLjR4DqrQ3lTanyFnfy4fjPVf16jbiXkV158hRgS4A+WAVP3DVg69cgVN1pae2Txn2An8ruZoOha8lSFnrBtzzKYnNJjs3qmckDgOQDMvO/MNpE4ZXk8RlJZ1GXDHAkiWGbCJwTdQ+LQ+Vk5HaE9OhGlgPkAdCxSKQpjhS+KCC8hyP3F2mQgi6mb9ZZelp/4gAPtAcM79dI1AazdhFllbDhX8I7UiO3o04j7XHi9CNUl7kqYKOn4wJ1u0zOAqIFs6NBvkgeJCwhzmqdCidZfwRg0TSH/QmoC2Ih4UNmnkCHbVZNiYn2x+o61CxkUQ0NU5Yc/Z948FJW0dMZrHB5t35e4IC+xy8wIg4JIzMxVCk3lavCrkU/G/1ehjrym9ccUV+bUPnbkWjCnSz30XNPAFObdjKWumdrPJGCSbHcead2Ne4Tw9ihh6+s9K
  template:
    metadata:
      creationTimestamp: null
      name: app_name
      namespace: app_namespace
```

## 8. Розміщуємо маніфест у репозиторії Flux та через декілька хвилин перевіряємо наявність потрібного секрету

```shell
kubectl get secrets -A
```

---
 ← [Повернутись до змісту](../README.md)  