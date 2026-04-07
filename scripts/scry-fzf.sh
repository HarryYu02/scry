#!/usr/bin/env bash

if [ $# -eq 0 ]; then
    echo "ERROR: source not provided"
    exit 1
fi

src=$1

fzf --disabled \
    --bind "start,change:reload(scry search --n=0 --docs --meta $1 {q} | jq -r '.id + \"\t\" + .title')" \
    --delimiter '\t' \
    --with-nth 2 \
    --border="rounded" \
    --border-label="$1" \
    --preview "scry open --stdout $1 {1}" \
    --preview-window="right:50%:wrap" \
    | cut -f1 | xargs -r scry open $1
