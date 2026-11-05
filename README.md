# KubeSay

KubeSay explains why a Kubernetes workload is failing — read-only — using a local or free LLM.

## What it does

KubeSay is a Go-native, read-only Kubernetes troubleshooting CLI. It gathers the
evidence a human operator would collect — pod status, events, current and previous
container logs, the owner chain, and resource limits — and uses a local or cloud-free
LLM to explain the likely cause. The explanation separates observed facts from
hypotheses and recommends next steps. KubeSay never modifies the cluster.

By default it runs against a local [Ollama](https://ollama.com/) backend, so it works
offline with no API keys. The free [Groq](https://groq.com/) tier is supported as an
optional cloud-free backend.

## Usage

```sh
# List pods in a namespace
kubesay list pods -n <namespace>

# Explain why a target (pod or higher-level workload) is failing
kubesay investigate <target>
```

## Status

v1 is under active development (Phase 0 scaffolding). The build is phased:

- **Phase 0 — Setup:** Cobra CLI that connects to a cluster and lists pods.
- **Phase 1 — Evidence Collector:** read-only, concurrently assembled evidence bundle.
- **Phase 2 — Grounded Reasoner:** `kubesay investigate <target>` explains a real failure.
- **Phase 3 — Eval harness + polish:** measurable per-backend diagnosis accuracy.

Module path: `github.com/rimasid/kubesay`.

> This is an early stub. A full README with a demo, install instructions, architecture
> diagram, and evaluation results will follow as v1 matures.
