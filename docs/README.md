---
description: A brief introduction to RocketBlend.
---

# Introduction to RocketBlend

{% hint style="warning" %}
Please note that RocketBlend is still in the early stages of development and may undergo significant changes as it continues to evolve.
{% endhint %}

## What is RocketBlend?

[RocketBlend ](https://github.com/rocketblend/rocketblend)is a command line (CLI) dependency manager and toolkit for the 3D graphics software, [Blender](https://www.blender.org/). It helps manage builds, add-ons, and offers extended blender commands for your projects.

## Features

* Define dependencies for your projects, including builds and add-ons, and have them set up correctly when you open your project.
* Render animations with improved output path templating, versioning, resumability, render engine selection, and more.
* Use arguments in any order without needing to remember specific commands or the order of flags, unlike Blender's command line arguments.
* Lightweight and requires no GUI or background processes.
* Customisable, with all data stored in .json files instead of internal .blend file data blocks or file headers.
* Write your own packages for any Blender build or add-on, with support for alternative forks of Blender, such as Bforartists and K-Cycles.
* Host your own library of packages, providing flexibility beyond the official ones. Fix broken or missing packages and continue working. Libraries are just git repositories containing `rocketpack.json` files that define an add-on or build. You can fork existing ones or create your own for use in your projects.
* Cross-platform compatibility with macOS, Windows, and Linux.
* Open source and receptive to community feedback.

## Want to jump right in?

Feeling like an eager beaver? Jump into the quick start guide to start using RocketBlend

{% content-ref url="getting-started/quick-start.md" %}
[quick-start.md](getting-started/quick-start.md)
{% endcontent-ref %}

## Want a deeper dive?

Dive a little deeper and start exploring our references documentation to get an idea of everything that's possible with RocketBlend:

{% content-ref url="reference/cli/" %}
[cli](reference/cli/)
{% endcontent-ref %}

## Not a Fan of Command Lines?

Give the desktop application a try. It offers most of the same functionality as the CLI, plus some extra features, all within an easy-to-use interface.

{% embed url="https://github.com/rocketblend/rocketblend-desktop" %}
