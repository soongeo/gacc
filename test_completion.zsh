source <(./gacc completion zsh)
compdef _gacc ./gacc
autoload -Uz compinit && compinit
