---
description: How to set up an environment for developing on RocketBlend.
---

# Getting Setup

## Prerequisites <a href="#prerequisites" id="prerequisites"></a>

* Linux or Windows x64
* Go 1.19
* Git
* GoReleaser (Release Only)
* Cosign (Release Only)

## Building RocketBlend <a href="#building-helm" id="building-helm"></a>

We use Make to build our programs. The simplest way to get started is:

{% tabs %}
{% tab title="CLI" %}
```bash
make build
```
{% endtab %}
{% endtabs %}

## Running tests <a href="#running-tests" id="running-tests"></a>

```bash
make check
```

## Creating a Release

```
make dry
make release version=v1.5.0
```
