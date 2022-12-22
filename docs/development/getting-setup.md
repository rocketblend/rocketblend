---
description: >-
  This guide explains how to set up your environment for developing on
  RocketBlend.
---

# Getting Setup

## Prerequisites <a href="#prerequisites" id="prerequisites"></a>

* Linux or Windows x64
* Go 1.19
* Git

## Building RocketBlend <a href="#building-helm" id="building-helm"></a>

We use Make to build our programs. The simplest way to get started is:

{% tabs %}
{% tab title="CLI" %}
```bash
make build
```
{% endtab %}

{% tab title="Launcher" %}
```bash
make launcher
```
{% endtab %}
{% endtabs %}

## Running tests <a href="#running-tests" id="running-tests"></a>

```bash
make check
```

