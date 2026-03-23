# Scope & hypotheses (M2 MMU prototype)

## Goal
Concevoir et implémenter un **prototype de MMU logicielle** pour DPUs UPMEM dans **uPIMulator (golang_vm)**, en s’inspirant de Faho (ou “prototype équivalent” comme autorisé par le sujet).

## Execution model (SWAP)
Modèle **SWAP** (Faho-like):
- À un instant donné, **un seul process** est résident en WRAM/IRAM (zone user).
- Le context switch sauvegarde/restaure l’image WRAM/IRAM user via MRAM.

## Constraints (UPMEM)
- Pas de MMU matérielle, pas de trap sur LOAD/STORE.
- WRAM: accès direct par instructions LW/SW.
- IRAM: fetch par PC.
- MRAM: accès via DMA (donc médiable par syscalls/runtime).

## What is in scope
- IRAM: invariant trampoline à IRAM[0x0000] + guard de chargement .text.
- MRAM: virtualisation via paging 64KB, table par process, faults.
- WRAM: séparation user/kernel + sbrk/malloc/free sûrs + batch validation.
- Faults software: violation -> kill process + métriques.
- Mesures cycle-level dans uPIMulator.

## Out of scope (MVP)
- Intégration complète du vrai code Faho (K1/K2/K3) si non disponible.
- Validation sur matériel UPMEM réel (non dispo).
