# Вибір концепції

Розгортання Kubernetes кластерів на базі k3d для застосунку IlonaBot з використанням GitHub Actions.

## Вступ

Завдання хакатону – розробка боту для Slack, що автоматизує процеси розгортання та відстеження стану артефактів в середовищах `dev`, `qa`, `taging` та `production`, що дозволить розробниками самостійно планувати та відстежувати процес розробки з мінімальним залученням інших команд. Цей документ ставить на меті вибір відповідного інструменту для розгортання кластерів Kubernetes. Розглянемо можливі варіанти: minikube, kind, k3d.

## Характеристики

### Minikube

- **Підтримувані ОС та архітектури:** Підтримує різні ОС, включаючи Linux, macOS, та Windows.
  
- **Можливість автоматизації:** Обмежена можливість автоматизації порівняно з іншими інструментами.

- **Додаткові функції:** Має вбудований функціонал, такий як Dashboard для моніторингу та огляду кластера.

### kind (Kubernetes IN Docker)

- **Підтримувані ОС та архітектури:** Широкі можливості, оскільки використовує Docker як базовий шар.

- **Можливість автоматизації:** Легко інтегрується в CI/CD системи для автоматичного тестування.

- **Додаткові функції:** Орієнтований на швидкість та легкість використання. Підтримує багатовузлові (включаючи HA) кластери. Kind підтримує створення збірок релізів з сирців.

### k3d

- **Підтримувані ОС та архітектури:** Підтримує різні ОС, але інтегрований з Docker.

- **Можливість автоматизації:** Легко автоматизується через інтерфейс командного рядка.

- **Додаткові функції:** Забезпечує легке створення та тестування кластерів Kubernetes у Docker-контейнерах. Є розширення VSCode для роботи з кластером k3d у VSCode – [vscode-k3d](https://github.com/inercia/vscode-k3d/).

## Переваги та Недоліки

| Характеристика | [Minikube](https://minikube.sigs.k8s.io/)| [kind](https://kind.sigs.k8s.io) | [k3d](https://k3d.io/) |
|--|--|--|--|
| Легкість використання | Зручний для початківців, але може бути обмежений | Легко встановлюється і використовується | Простий у використанні та конфігурації |
| Швидкість розгортання | Залежить від обраного гіпервізору, може бути повільним | Швидке розгортання у Docker-контейнерах | Швидке створення та тестування у Docker-контейнерах |
| Стабільність роботи | Стабільний, але може виникати обмеження масштабування | Стабільний та надійний | Стабільний і легко масштабується |
| Документація та спільнота | Добра документація, активна спільнота | Прийнятна документація, сильна спільнота | Задовільна документація, активна спільнота |
| Налаштування та використання | Має багато параметрів, що може бути складним для новачків | Простий у використанні та конфігурації | Легко конфігурується та використовується |
| Масштабування | Тільки один вузол | Це єдине локальне рішення з HA-кластером з кількома control-plane вузлами. | Підтримує кілька кластерів і підтримує кілька робочих вузлів на кластер. Легко зупиняє та запускає кластери без втрати їх стану. |
| Використання з Podman | Обмежене (експериментальне використання) | Краща підтримка роботи, робота у rootless режимі | Обмежене (експериментальне використання)|

## Демонстрація

Рекомендований інструмент: k3d Розгортання програми "Hello World" на Kubernetes

[![Застосунок на Kubernetes](https://asciinema.org/a/622883.svg)](https://asciinema.org/a/622883)

## Висновки

Після уважного порівняльного аналізу було вирішено використовувати k3d для розгортання кластерів Kubernetes. Це легкий інструмент, що дозволяє швидко створювати та тестувати кластери у Docker-контейнерах. Його легко інтегрувати в CI/CD системи для автоматичного тестування. Інструмент має добру документацію та активну спільноту, що дозволить швидко розвʼязувати проблеми, що виникають під час розробки.

Для керування розгортаннями артефактів в кластері Kubernetes було вирішено використовувати [Flux](https://fluxcd.io/), як складову Continuous Delivery, що дозволить автоматизувати процес розгортання та відстеження стану артефактів в середовищах `dev`, `qa`, `taging` та `production` та налагодити комунікацію з інфраструктурою за допомогою бота Slack. Процеси Continuous Integration будуть реалізовані за допомогою GitHub Actions.

---
← [Повернутись до змісту](../README.md)  
